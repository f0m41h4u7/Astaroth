package collector

import (
	"sync"
	"time"
)

type Snapshot struct {
	Metrics []Metric
}

type Collector struct {
	mutex       sync.RWMutex
	subscribers []chan Snapshot
	snapCurrent *Snapshot
}

func NewCollector(mt []Metric) *Collector {
	//	log.Printf("%+v", mt)
	return &Collector{
		subscribers: []chan Snapshot{},
		snapCurrent: &Snapshot{
			Metrics: mt,
		},
	}
}

func (c *Collector) Subscribe() <-chan Snapshot {
	ch := make(chan Snapshot)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.subscribers = append(c.subscribers, ch)

	return c.subscribers[len(c.subscribers)-1]
}

func (c *Collector) CollectStats() error {
	var wg sync.WaitGroup
	errs := make(chan error)

	for _, mt := range c.snapCurrent.Metrics {
		wg.Add(1)
		go func(mt Metric) {
			errs <- mt.Get(&wg)
		}(mt)
		//	log.Printf("%+v", mt)
	}

	wg.Wait()
	if len(errs) != 0 {
		err := <-errs
		if err != nil {
			return <-errs
		}
	}

	c.mutex.RLock()
	for _, ch := range c.subscribers {
		ch <- *c.snapCurrent
	}
	c.mutex.RUnlock()

	return nil
}

func (c *Collector) Run(interval int64, done <-chan struct{}) error {
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
