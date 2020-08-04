package server

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector/linux"
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
	collector       *linux.Collector
}

// InitServer initializes Server.
func InitServer(addr string, col *linux.Collector) *Server {
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

// Start Server.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("server starts")
	err = s.grpc.Serve(lis)

	return err
}

// Connect to Server.
func (s *Server) Connect(ctx context.Context, req *api.ConnectRequest) (*empty.Empty, error) {
	resp := &empty.Empty{}
	if (req.SendInterval > req.AverageInterval) || (req.SendInterval <= 0) {
		return resp, errWrongInervals
	}

	s.sendInterval = req.SendInterval
	s.averageInterval = req.AverageInterval

	return resp, nil
}

// GetStats returns collected statistics.
func (s *Server) GetStats(_ *empty.Empty, srv api.Astaroth_GetStatsServer) error {
	log.Printf("new stats listener")
	statsChan := s.collector.Subscribe()

	size := int(s.averageInterval / s.sendInterval)
	stats := []linux.Snapshot{}
	cnt := 0

	ticker := time.NewTicker(time.Duration(s.sendInterval) * time.Second)
	stop := false
	for !stop {
		select {
		case <-statsChan:
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
			log.Printf("sending data: %s", msg.String())
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

func (s *Server) averageStats(snapshots []linux.Snapshot) *api.Stats {
	st := new(api.Stats)
	size := len(snapshots)

	if config.RequiredMetrics.Metrics[config.CPU] == config.On {
		st.CPU = &api.CPU{
			User:   0,
			System: 0,
		}
		for _, snap := range snapshots {
			st.CPU.User += snap.CPU.User
			st.CPU.System += snap.CPU.System
		}
		st.CPU.User /= float64(size)
		st.CPU.System /= float64(size)
	}
	if config.RequiredMetrics.Metrics[config.LoadAvg] == config.On {
		st.LoadAvg = &api.LoadAvg{
			OneMin:       0.0,
			FiveMin:      0.0,
			FifteenMin:   0.0,
			ProcsRunning: 0,
			TotalProcs:   0,
		}
		for _, snap := range snapshots {
			st.LoadAvg.OneMin += snap.LoadAvg.OneMin
			st.LoadAvg.FiveMin += snap.LoadAvg.FiveMin
			st.LoadAvg.FifteenMin += snap.LoadAvg.FifteenMin
			st.LoadAvg.ProcsRunning += snap.LoadAvg.ProcsRunning
			st.LoadAvg.TotalProcs += snap.LoadAvg.TotalProcs
		}
		st.LoadAvg.OneMin /= float64(size)
		st.LoadAvg.FiveMin /= float64(size)
		st.LoadAvg.FifteenMin /= float64(size)
		st.LoadAvg.ProcsRunning /= int64(size)
		st.LoadAvg.TotalProcs /= int64(size)
	}

	return st
}

// Stop Server.
func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
