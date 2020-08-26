package collector

import (
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type DiskDataMetric struct {
	DiskData *api.DiskData
}

func (d *DiskDataMetric) Average(wg *sync.WaitGroup, st *api.Stats, snapshots *[]Snapshot, idx int) {
	defer wg.Done()
	size := int64(len(*snapshots))

	st.DiskData = &api.DiskData{
		Data: []*api.FilesystemData{},
	}

	for i := 0; i < len((*snapshots)[0].Metrics[idx].(*DiskDataMetric).DiskData.Data); i++ {
		st.DiskData.Data = append(st.DiskData.Data, (*snapshots)[0].Metrics[idx].(*DiskDataMetric).DiskData.Data[i])
	}
	for i := 1; i < len(*snapshots); i++ {
		for i, d := range (*snapshots)[i].Metrics[idx].(*DiskDataMetric).DiskData.Data {
			st.DiskData.Data[i].Used += d.Used
			st.DiskData.Data[i].Inode += d.Inode
		}
	}
	for _, d := range st.DiskData.Data {
		d.Used /= size
		d.Inode /= size
	}
}
