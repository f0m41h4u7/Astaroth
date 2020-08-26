package collector

import (
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type NetworkStatsMetric struct {
	NetworkStats *api.NetworkStats
}

func (n *NetworkStatsMetric) Average(wg *sync.WaitGroup, st *api.Stats, snapshots *[]Snapshot, idx int) {
	defer wg.Done()
	size := int64(len(*snapshots))
	if size == int64(0) {
		return
	}

	st.NetworkStats = &api.NetworkStats{
		ListenSockets: (*snapshots)[len((*snapshots))-1].Metrics[idx].(*NetworkStatsMetric).NetworkStats.ListenSockets,
		TCPConnStates: &api.States{
			LISTEN:     0,
			ESTAB:      0,
			FIN_WAIT:   0,
			SYN_RCV:    0,
			TIME_WAIT:  0,
			CLOSE_WAIT: 0,
		},
	}

	for _, snap := range *snapshots {
		if snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.String() == "" {
			continue
		}
		st.NetworkStats.TCPConnStates.LISTEN += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.LISTEN
		st.NetworkStats.TCPConnStates.ESTAB += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.ESTAB
		st.NetworkStats.TCPConnStates.FIN_WAIT += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.FIN_WAIT
		st.NetworkStats.TCPConnStates.SYN_RCV += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.SYN_RCV
		st.NetworkStats.TCPConnStates.TIME_WAIT += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.TIME_WAIT
		st.NetworkStats.TCPConnStates.CLOSE_WAIT += snap.Metrics[idx].(*NetworkStatsMetric).NetworkStats.TCPConnStates.CLOSE_WAIT
	}

	st.NetworkStats.TCPConnStates.LISTEN /= size
	st.NetworkStats.TCPConnStates.ESTAB /= size
	st.NetworkStats.TCPConnStates.FIN_WAIT /= size
	st.NetworkStats.TCPConnStates.SYN_RCV /= size
	st.NetworkStats.TCPConnStates.TIME_WAIT /= size
	st.NetworkStats.TCPConnStates.CLOSE_WAIT /= size
}

func newMap() map[string]int64 {
	states := map[string]int64{}
	states["LISTEN"] = 0
	states["ESTAB"] = 0
	states["FIN-WAIT"] = 0
	states["SYN-RCV"] = 0
	states["TIME-WAIT"] = 0
	states["CLOSE-WAIT"] = 0

	return states
}
