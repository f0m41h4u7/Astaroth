package linux

import (
	"sync"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

// MetricsStorage stores last metrics measurements.
type MetricsStorage struct {
	idx     int64
	cpu     []*api.CPU
	loadAvg []*api.LoadAvg
	// diskData []*api.DiskData
	// netStats []*api.NetworkStats
	// to be continued ...
}

type Collector struct {
	size    int64
	storage *MetricsStorage
}

func NewCollector(size int64) *Collector {
	return &Collector{
		size: size,
		storage: &MetricsStorage{
			idx:     0,
			cpu:     make([]*api.CPU, size),
			loadAvg: make([]*api.LoadAvg, size),
		},
	}
}

// CollectStats makes measurements and updates MetricsStorage.
func (c *Collector) CollectStats() error {
	var wg sync.WaitGroup
	errs := make(chan error)

	if config.RequiredMetrics.Metrics[config.CPU] == config.On {
		wg.Add(1)
		go func() {
			var mutex sync.RWMutex
			cpu, err := GetCPU(&wg)
			if err != nil {
				errs <- err
				return
			}
			mutex.RLock()
			if c.storage.idx < c.size {
				c.storage.cpu[c.storage.idx] = cpu
				c.storage.idx++
				mutex.RUnlock()
				return
			}
			mutex.RUnlock()
			for i := 0; i < int(c.size-1); i++ {
				c.storage.cpu[i] = c.storage.cpu[i+1]
			}
			c.storage.cpu[c.size-1] = cpu
		}()
	}

	if config.RequiredMetrics.Metrics[config.LoadAvg] == config.On {
		wg.Add(1)
		go func() {
			var mutex sync.RWMutex
			la, err := GetLoadAvg(&wg)
			if err != nil {
				errs <- err
			}
			mutex.RLock()
			if c.storage.idx < c.size {
				c.storage.loadAvg[c.storage.idx] = la
				c.storage.idx++
				mutex.RUnlock()
				return
			}
			mutex.RUnlock()
			for i := 0; i < int(c.size-1); i++ {
				c.storage.loadAvg[i] = c.storage.loadAvg[i+1]
			}
			c.storage.loadAvg[c.size-1] = la
		}()
	}

	wg.Wait()
	if (len(errs) != 0) && (<-errs != nil) {
		return <-errs
	}
	return nil
}

// SendStats returns average values of metrics.
func (c *Collector) SendStats() *api.Stats {
	st := new(api.Stats)
	if config.RequiredMetrics.Metrics[config.CPU] == config.On {
		st.CPU = &api.CPU{
			User:   0,
			System: 0,
		}
		for _, cpu := range c.storage.cpu {
			st.CPU.User += cpu.User
			st.CPU.System += cpu.System
		}
		st.CPU.User /= float64(c.size)
		st.CPU.System /= float64(c.size)
	}
	if config.RequiredMetrics.Metrics[config.LoadAvg] == config.On {
		st.LoadAvg = &api.LoadAvg{
			OneMin:       0.0,
			FiveMin:      0.0,
			FifteenMin:   0.0,
			ProcsRunning: 0,
			TotalProcs:   0,
		}
		for _, la := range c.storage.loadAvg {
			st.LoadAvg.OneMin += la.OneMin
			st.LoadAvg.FiveMin += la.FiveMin
			st.LoadAvg.FifteenMin += la.FifteenMin
			st.LoadAvg.ProcsRunning += la.ProcsRunning
			st.LoadAvg.TotalProcs += la.TotalProcs
		}
		st.LoadAvg.OneMin /= float64(c.size)
		st.LoadAvg.FiveMin /= float64(c.size)
		st.LoadAvg.FifteenMin /= float64(c.size)
		st.LoadAvg.ProcsRunning /= c.size
		st.LoadAvg.TotalProcs /= c.size
	}
	return st
}
