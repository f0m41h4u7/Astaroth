package server

import (
	"math"
	"sync"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector/linux"
)

func averageCPU(wg *sync.WaitGroup, st *api.Stats, snapshots *[]linux.Snapshot) {
	defer wg.Done()
	size := float64(len(*snapshots))

	st.CPU = &api.CPU{
		User:   0,
		System: 0,
	}
	for _, snap := range *snapshots {
		st.CPU.User += snap.CPU.User
		st.CPU.System += snap.CPU.System
	}
	st.CPU.User = math.Round(st.CPU.User/size*10) / 10
	st.CPU.System = math.Round(st.CPU.System/size*10) / 10
}

func averageLoadAvg(wg *sync.WaitGroup, st *api.Stats, snapshots *[]linux.Snapshot) {
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
		st.LoadAvg.OneMin += snap.LoadAvg.OneMin
		st.LoadAvg.FiveMin += snap.LoadAvg.FiveMin
		st.LoadAvg.FifteenMin += snap.LoadAvg.FifteenMin
		st.LoadAvg.ProcsRunning += snap.LoadAvg.ProcsRunning
		st.LoadAvg.TotalProcs += snap.LoadAvg.TotalProcs
	}
	st.LoadAvg.OneMin = math.Round(st.LoadAvg.OneMin/float64(size)*10) / 10
	st.LoadAvg.FiveMin = math.Round(st.LoadAvg.FiveMin/float64(size)*10) / 10
	st.LoadAvg.FifteenMin = math.Round(st.LoadAvg.FifteenMin/float64(size)*10) / 10
	st.LoadAvg.ProcsRunning /= size
	st.LoadAvg.TotalProcs /= size
}

func averageDiskData(wg *sync.WaitGroup, st *api.Stats, snapshots *[]linux.Snapshot) {
	defer wg.Done()
	size := int64(len(*snapshots))

	st.DiskData = &api.DiskData{
		Data: []*api.FilesystemData{},
	}
	for i := 0; i < len((*snapshots)[0].DiskData.Data); i++ {
		st.DiskData.Data = append(st.DiskData.Data, (*snapshots)[0].DiskData.Data[i])
	}
	for i := 1; i < len(*snapshots); i++ {
		for i, d := range (*snapshots)[i].DiskData.Data {
			st.DiskData.Data[i].Used += d.Used
			st.DiskData.Data[i].Inode += d.Inode
		}
	}
	for _, d := range st.DiskData.Data {
		d.Used /= size
		d.Inode /= size
	}
}

func averageNetworkStats(wg *sync.WaitGroup, st *api.Stats, snapshots *[]linux.Snapshot) {
	defer wg.Done()

	size := int64(len(*snapshots))
	if size == int64(0) {
		return
	}
	if (*snapshots)[len((*snapshots))-1].NetworkStats == nil {
		return
	}

	st.NetworkStats = &api.NetworkStats{
		ListenSockets: (*snapshots)[len((*snapshots))-1].NetworkStats.ListenSockets,
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
		st.NetworkStats.TCPConnStates.LISTEN += snap.NetworkStats.TCPConnStates.LISTEN
		st.NetworkStats.TCPConnStates.ESTAB += snap.NetworkStats.TCPConnStates.ESTAB
		st.NetworkStats.TCPConnStates.FIN_WAIT += snap.NetworkStats.TCPConnStates.FIN_WAIT
		st.NetworkStats.TCPConnStates.SYN_RCV += snap.NetworkStats.TCPConnStates.SYN_RCV
		st.NetworkStats.TCPConnStates.TIME_WAIT += snap.NetworkStats.TCPConnStates.TIME_WAIT
		st.NetworkStats.TCPConnStates.CLOSE_WAIT += snap.NetworkStats.TCPConnStates.CLOSE_WAIT
	}

	st.NetworkStats.TCPConnStates.LISTEN /= size
	st.NetworkStats.TCPConnStates.ESTAB /= size
	st.NetworkStats.TCPConnStates.FIN_WAIT /= size
	st.NetworkStats.TCPConnStates.SYN_RCV /= size
	st.NetworkStats.TCPConnStates.TIME_WAIT /= size
	st.NetworkStats.TCPConnStates.CLOSE_WAIT /= size
}

func (s *Server) averageStats(snapshots []linux.Snapshot) *api.Stats {
	st := new(api.Stats)
	var wg sync.WaitGroup

	if config.RequiredMetrics.Metrics[config.CPU] == config.On {
		wg.Add(1)
		go averageCPU(&wg, st, &snapshots)
	}

	if config.RequiredMetrics.Metrics[config.LoadAvg] == config.On {
		wg.Add(1)
		go averageLoadAvg(&wg, st, &snapshots)
	}

	if config.RequiredMetrics.Metrics[config.DiskData] == config.On {
		wg.Add(1)
		go averageDiskData(&wg, st, &snapshots)
	}

	if config.RequiredMetrics.Metrics[config.NetworkStats] == config.On {
		wg.Add(1)
		go averageNetworkStats(&wg, st, &snapshots)
	}
	wg.Wait()

	return st
}
