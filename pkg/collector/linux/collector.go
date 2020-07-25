package linux

import (
	"errors"
)

var ErrWrongIntervals = errors.New("time intervals are less or equal zero")

type MetricsStorage struct {
	// cpu []*api.CPU
	// loadAvg  []*api.LoadAvg
	// diskData []*api.DiskData
	// to be continued ...
}

//nolint:maligned
type Collector struct {
	sendInterval int64
	avgInterval  int64
	storage      MetricsStorage
}

func NewCollector(s int64, a int64) (Collector, error) {
	var c Collector
	if (s <= 0) || (a <= 0) {
		return c, ErrWrongIntervals
	}
	c.sendInterval = s
	c.avgInterval = a
	c.storage = MetricsStorage{}
	return c, nil
}
