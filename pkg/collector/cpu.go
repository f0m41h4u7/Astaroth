package collector

import (
	"math"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type CPUMetric struct {
	CPU *api.CPU
}

func (c *CPUMetric) Average(wg *sync.WaitGroup, st *api.Stats, snapshots *[]Snapshot, idx int) {
	defer wg.Done()
	size := float64(len(*snapshots))

	st.CPU = &api.CPU{
		User:   0,
		System: 0,
	}
	for _, snap := range *snapshots {
		st.CPU.User += snap.Metrics[idx].(*CPUMetric).CPU.User
		st.CPU.System += snap.Metrics[idx].(*CPUMetric).CPU.System
	}
	st.CPU.User = math.Round(st.CPU.User/size*10) / 10
	st.CPU.System = math.Round(st.CPU.System/size*10) / 10
}
