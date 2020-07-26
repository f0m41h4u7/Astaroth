package linux

import (
	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

// MetricsStorage stores last metrics measurements.
type MetricsStorage struct {
	idx int64
	cpu []*api.CPU
	// loadAvg  []*api.LoadAvg
	// diskData []*api.DiskData
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
			idx: 0,
			cpu: make([]*api.CPU, size),
		},
	}
}

// CollectStats makes measurements and updates MetricsStorage.
func (c *Collector) CollectStats() error {
	cpu, err := GetCPU()
	if err != nil {
		return err
	}
	if c.storage.idx < c.size {
		c.storage.cpu[c.storage.idx] = cpu
		c.storage.idx++
		return nil
	}
	for i := 0; i < int(c.size-1); i++ {
		c.storage.cpu[i] = c.storage.cpu[i+1]
	}
	c.storage.cpu[c.size-1] = cpu
	return nil
}

// SendStats returns average values of metrics.
func (c *Collector) SendStats() *api.Stats {
	st := new(api.Stats)
	st.CPU = new(api.CPU)
	st.CPU.User = 0
	st.CPU.System = 0
	for _, cpu := range c.storage.cpu {
		st.CPU.User += cpu.User
		st.CPU.System += cpu.System
	}
	st.CPU.User /= float32(c.size)
	st.CPU.System /= float32(c.size)
	return st
}
