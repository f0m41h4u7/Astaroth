package server

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

//go:generate protoc --proto_path=../../pkg/api/ --go_out=plugins=grpc:../../pkg/api ../../pkg/api/stats.proto

var errWrongInervals = errors.New("send interval should be less than average interval and positive")

// Server is a GRPC streaming server.
type Server struct {
	grpc            *grpc.Server
	addr            string
	sendInterval    int64
	averageInterval int64
	collector       *collector.Collector
}

func InitServer(addr string, col *collector.Collector) *Server {
	s := &Server{
		addr:            addr,
		sendInterval:    3,
		averageInterval: 9,
		collector:       col,
	}
	grpcServer := grpc.NewServer()
	api.RegisterAstarothServer(grpcServer, s)
	s.grpc = grpcServer

	return s
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("server starts")
	err = s.grpc.Serve(lis)

	return err
}

func (s *Server) Connect(ctx context.Context, req *api.ConnectRequest) (*empty.Empty, error) {
	resp := &empty.Empty{}
	if (req.SendInterval > req.AverageInterval) || (req.SendInterval <= 0) {
		return resp, errWrongInervals
	}

	s.sendInterval = req.SendInterval
	s.averageInterval = req.AverageInterval

	return resp, nil
}

func (s *Server) GetStats(_ *empty.Empty, srv api.Astaroth_GetStatsServer) error {
	log.Printf("new stats listener")
	statsChan := s.collector.Subscribe()

	size := int(s.averageInterval / s.sendInterval)
	stats := []collector.Snapshot{}
	cnt := 0

	ticker := time.NewTicker(time.Duration(s.sendInterval) * time.Second)
	stop := false
	for !stop {
		select {
		case <-statsChan:
			log.Printf("received: %q", <-statsChan)
			if len(stats) < size {
				stats = append(stats, <-statsChan)
				cnt++
			} else {
				for i := 0; i < size-1; i++ {
					stats[i] = stats[i+1]
				}
				stats[size-1] = <-statsChan
			}
		case <-ticker.C:
			if len(stats) == 0 {
				continue
			}
			msg := s.averageStats(stats)
			stats = stats[:0]
			cnt = 0
			log.Printf("after average")
			if err := srv.Send(msg); err != nil {
				log.Printf("unable to send message to stats listener: %v", err)
				stop = true
			}
		case <-srv.Context().Done():
			log.Printf("stats listener disconnected")
			stop = true
		}
	}

	return nil
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
