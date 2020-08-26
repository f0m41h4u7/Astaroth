package collector

import (
	"math"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type LoadAvgMetric struct {
	LoadAvg *api.LoadAvg
}

func (l *LoadAvgMetric) Average(wg *sync.WaitGroup, st *api.Stats, snapshots *[]Snapshot, idx int) {
	defer wg.Done()
	size := int64(len(*snapshots))

	st.LoadAvg = &api.LoadAvg{
		OneMin:       0.0,
		FiveMin:      0.0,
		FifteenMin:   0.0,
		ProcsRunning: 0,
		TotalProcs:   0,
	}

	for _, snap := range *snapshots {
		st.LoadAvg.OneMin += snap.Metrics[idx].(*LoadAvgMetric).LoadAvg.OneMin
		st.LoadAvg.FiveMin += snap.Metrics[idx].(*LoadAvgMetric).LoadAvg.FiveMin
		st.LoadAvg.FifteenMin += snap.Metrics[idx].(*LoadAvgMetric).LoadAvg.FifteenMin
		st.LoadAvg.ProcsRunning += snap.Metrics[idx].(*LoadAvgMetric).LoadAvg.ProcsRunning
		st.LoadAvg.TotalProcs += snap.Metrics[idx].(*LoadAvgMetric).LoadAvg.TotalProcs
	}
	st.LoadAvg.OneMin = math.Round(st.LoadAvg.OneMin/float64(size)*10) / 10
	st.LoadAvg.FiveMin = math.Round(st.LoadAvg.FiveMin/float64(size)*10) / 10
	st.LoadAvg.FifteenMin = math.Round(st.LoadAvg.FifteenMin/float64(size)*10) / 10
	st.LoadAvg.ProcsRunning /= size
	st.LoadAvg.TotalProcs /= size
}
