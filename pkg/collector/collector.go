package collector

import (
	"sync"
	"time"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type Snapshot struct {
	CPU          *api.CPU
	LoadAvg      *api.LoadAvg
	DiskData     *api.DiskData
	NetworkStats *api.NetworkStats
	TopTalkers   *api.TopTalkers
}

type Collector struct {
	mutex       sync.RWMutex
	subscribers []chan Snapshot
}

func NewCollector() *Collector {
	return &Collector{
		subscribers: []chan Snapshot{},
	}
}

func (c *Collector) Subscribe() chan Snapshot {
	ch := make(chan Snapshot)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.subscribers = append(c.subscribers, ch)

	return c.subscribers[len(c.subscribers)-1]
}

func (c *Collector) CollectStats() error {
	var wg sync.WaitGroup
	var ss Snapshot
	errs := make(chan error)

	if config.RequiredMetrics.Metrics[config.CPU] == config.On {
		wg.Add(1)
		go func() {
			errs <- c.getCPU(&wg, &ss)
		}()
	}

	if config.RequiredMetrics.Metrics[config.LoadAvg] == config.On {
		wg.Add(1)
		go func() {
			errs <- c.getLoadAvg(&wg, &ss)
		}()
	}

	if config.RequiredMetrics.Metrics[config.DiskData] == config.On {
		wg.Add(1)
		go func() {
			errs <- c.getDiskData(&wg, &ss)
		}()
	}

	if config.RequiredMetrics.Metrics[config.NetworkStats] == config.On {
		wg.Add(1)
		go func() {
			errs <- c.getNetworkStats(&wg, &ss)
		}()
	}

	if config.RequiredMetrics.Metrics[config.TopTalkers] == config.On {
		wg.Add(1)
		go func() {
			errs <- c.getTopTalkers(&wg, &ss)
		}()
	}

	wg.Wait()
	if (len(errs) != 0) && (<-errs != nil) {
		return <-errs
	}

	c.mutex.RLock()
	for _, ch := range c.subscribers {
		ch <- ss
	}
	c.mutex.RUnlock()

	return nil
}

func (c *Collector) Run(interval int64, done <-chan int) error {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			err := c.CollectStats()
			if err != nil {
				return err
			}
		}
	}
}
