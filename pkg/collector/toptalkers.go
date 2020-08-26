package collector

import (
	"sort"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var sumBytes int64

type TopTalkersMetric struct {
	TopTalkers *api.TopTalkers
}

func (t *TopTalkersMetric) Average(wg *sync.WaitGroup, st *api.Stats, snapshots *[]Snapshot, idx int) {
	defer wg.Done()
	protocols := map[string]int64{}
	sumBytes = 0
	st.TopTalkers = &api.TopTalkers{ByProtocol: []*api.ByProtocol{}, ByTraffic: []*api.ByTraffic{}}

	for _, snap := range *snapshots {
		if snap.Metrics[idx].(*TopTalkersMetric).TopTalkers == nil {
			continue
		}
		for _, bp := range snap.Metrics[idx].(*TopTalkersMetric).TopTalkers.ByProtocol {
			if _, ok := protocols[bp.Protocol]; !ok {
				st.TopTalkers.ByProtocol = append(st.TopTalkers.ByProtocol, &api.ByProtocol{Protocol: bp.Protocol, Bytes: bp.Bytes})
				protocols[bp.Protocol] = 1
			} else {
				protocols[bp.Protocol]++
			}
			if len(protocols) == 5 {
				break
			}
			sumBytes += bp.Bytes
		}

		for _, bt := range snap.Metrics[idx].(*TopTalkersMetric).TopTalkers.ByTraffic {
			f := false
			for _, t := range st.TopTalkers.ByTraffic {
				if (t.SourceAddr == bt.SourceAddr) && (t.DestAddr == bt.DestAddr) && (t.Protocol == bt.Protocol) {
					t.Bps += bt.Bps
					f = true

					break
				}
			}
			if !f {
				st.TopTalkers.ByTraffic = append(st.TopTalkers.ByTraffic, &api.ByTraffic{SourceAddr: bt.SourceAddr, DestAddr: bt.DestAddr, Protocol: bt.Protocol, Bps: bt.Bps})
			}
		}
	}

	sort.Slice(st.TopTalkers.ByProtocol, func(i, j int) bool {
		return st.TopTalkers.ByProtocol[i].Bytes > st.TopTalkers.ByProtocol[j].Bytes
	})
	sort.Slice(st.TopTalkers.ByTraffic, func(i, j int) bool {
		return st.TopTalkers.ByTraffic[i].Bps > st.TopTalkers.ByTraffic[j].Bps
	})

	for _, p := range st.TopTalkers.ByProtocol {
		p.Percentage = 100 * p.Bytes / sumBytes
	}
}
