package server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector/linux"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

//go:generate protoc --proto_path=../../pkg/api/ --go_out=plugins=grpc:../../pkg/api ../../pkg/api/stats.proto

var errWrongInervals = errors.New("send interval should be less than average interval and positive")

// Server is a GRPC streaming server
type Server struct {
	grpc            *grpc.Server
	addr            string
	sendInterval    int64
	averageInterval int64
	lock            sync.RWMutex
	collector       *linux.Collector
}

// InitServer initializes Server
func InitServer(addr string) *Server {
	s := &Server{
		addr:            addr,
		sendInterval:    3,
		averageInterval: 9,
	}
	grpcServer := grpc.NewServer()
	api.RegisterAstarothServer(grpcServer, s)
	s.grpc = grpcServer
	return s
}

// Start Server
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("server starts")
	err = s.grpc.Serve(lis)
	return err
}

// Connect to Server
func (s *Server) Connect(ctx context.Context, req *api.ConnectRequest) (*empty.Empty, error) {
	resp := &empty.Empty{}
	if (req.SendInterval > req.AverageInterval) || (req.SendInterval <= 0) {
		return resp, errWrongInervals
	}

	s.sendInterval = req.SendInterval
	s.averageInterval = req.AverageInterval
	return resp, nil
}

// GetStats returns collected statistics
func (s *Server) GetStats(_ *empty.Empty, srv api.Astaroth_GetStatsServer) error {
	log.Printf("new stats listener")
	s.collector = linux.NewCollector(s.averageInterval / s.sendInterval)
	ticker := time.NewTicker(time.Duration(s.sendInterval) * time.Second)
	stop := false
	cnt := int64(0)
	for !stop {
		select {
		case <-ticker.C:
			s.lock.Lock()
			err := s.collector.CollectStats()
			s.lock.Unlock()
			if err != nil {
				log.Printf("stats collecting error: %+v", err)
				stop = true
			}

			if cnt++; cnt > (s.averageInterval / s.sendInterval) {
				s.lock.RLock()
				msg := s.collector.SendStats()
				s.lock.RUnlock()
				log.Printf("sending data: %s", msg.String())
				if err := srv.Send(msg); err != nil {
					log.Printf("unable to send message to stats listener: %v", err)
					stop = true
				}
			}
		case <-srv.Context().Done():
			log.Printf("stats listener disconnected")
			stop = true
		}
	}
	return nil
}

// Stop Server
func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
