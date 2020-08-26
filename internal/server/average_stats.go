package server

import (
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector"
)

func (s *Server) averageStats(snapshots []collector.Snapshot) *api.Stats {
	st := new(api.Stats)
	var wg sync.WaitGroup

	for i, mt := range snapshots[0].Metrics {
		wg.Add(1)
		go mt.Average(&wg, st, &snapshots, i)
	}
	wg.Wait()

	return st
}
