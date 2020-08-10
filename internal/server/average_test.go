package server

import (
	"sync"
	"testing"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector/linux"
	"github.com/stretchr/testify/require"
)

var snap = []linux.Snapshot{
	{
		CPU: &api.CPU{
			User:   18.0,
			System: 1.0,
		},
		LoadAvg: &api.LoadAvg{
			OneMin:       2.5,
			FiveMin:      1.31,
			FifteenMin:   1.28,
			ProcsRunning: 1,
			TotalProcs:   882,
		},
		DiskData: &api.DiskData{
			Data: []*api.FilesystemData{
				{
					Filesystem: "tmpfs",
					Used:       10,
					Inode:      1,
				},
				{
					Filesystem: "tmpfs",
					Used:       0,
					Inode:      11,
				},
			},
		},
		NetworkStats: &api.NetworkStats{
			TCPConnStates: &api.States{
				LISTEN:     10,
				ESTAB:      5,
				FIN_WAIT:   3,
				SYN_RCV:    0,
				TIME_WAIT:  0,
				CLOSE_WAIT: 23,
			},
			ListenSockets: []*api.Sockets{
				{
					Program:  "sshd",
					PID:      7844,
					User:     "0",
					Protocol: "tcp",
					Port:     22,
				},
				{
					Program:  "vncserver",
					PID:      123,
					User:     "0",
					Protocol: "tcp",
					Port:     8888,
				},
			},
		},
	},
	{
		CPU: &api.CPU{
			User:   34.4,
			System: 10.6,
		},
		LoadAvg: &api.LoadAvg{
			OneMin:       3.5,
			FiveMin:      2.42,
			FifteenMin:   2.43,
			ProcsRunning: 2,
			TotalProcs:   982,
		},
		DiskData: &api.DiskData{
			Data: []*api.FilesystemData{
				{
					Filesystem: "tmpfs",
					Used:       11,
					Inode:      0,
				},
				{
					Filesystem: "tmpfs",
					Used:       20,
					Inode:      2,
				},
			},
		},
		NetworkStats: &api.NetworkStats{
			TCPConnStates: &api.States{
				LISTEN:     30,
				ESTAB:      1,
				FIN_WAIT:   0,
				SYN_RCV:    2,
				TIME_WAIT:  24,
				CLOSE_WAIT: 2,
			},
			ListenSockets: []*api.Sockets{
				{
					Program:  "sshd",
					PID:      7844,
					User:     "0",
					Protocol: "tcp",
					Port:     22,
				},
			},
		},
	},
}

func TestAverageStats(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var wg sync.WaitGroup
		st := new(api.Stats)
		wg.Add(4)
		go averageCPU(&wg, st, &snap)
		go averageLoadAvg(&wg, st, &snap)
		go averageDiskData(&wg, st, &snap)
		go averageNetworkStats(&wg, st, &snap)
		wg.Wait()

		require.Equal(t, &api.NetworkStats{
			TCPConnStates: &api.States{
				LISTEN:     20,
				ESTAB:      3,
				FIN_WAIT:   1,
				SYN_RCV:    1,
				TIME_WAIT:  12,
				CLOSE_WAIT: 12,
			},
			ListenSockets: []*api.Sockets{
				{
					Program:  "sshd",
					PID:      7844,
					User:     "0",
					Protocol: "tcp",
					Port:     22,
				},
			},
		}, st.NetworkStats)

		require.Equal(t, &api.DiskData{
			Data: []*api.FilesystemData{
				{
					Filesystem: "tmpfs",
					Used:       10,
					Inode:      0,
				},
				{
					Filesystem: "tmpfs",
					Used:       10,
					Inode:      6,
				},
			},
		}, st.DiskData)
	})
}
