package linux

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCPU(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cpu, err := readCPUFile("../../../tests/testdata/cpu.txt")
		require.Nil(t, err)
		require.Equal(t, 42055566250653, cpu)
	})
	t.Run("wrong file", func(t *testing.T) {
		_, err := readCPUFile("../../../tests/testdata/bad_file.txt")
		require.NotNil(t, err)
	})
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := readCPUFile("cpu.txt")
		require.NotNil(t, err)
	})
}

func TestLoadAvg(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		la, err := readLoadAvgFile("../../../tests/testdata/loadavg.txt")
		_, err = strconv.ParseFloat(la[0], 32)
		require.Nil(t, err)
		_, err = strconv.ParseFloat(la[1], 32)
		require.Nil(t, err)
		_, err = strconv.ParseFloat(la[2], 32)
		require.Nil(t, err)
		procsRun, err := strconv.Atoi(la[3])
		require.Nil(t, err)
		procsTotal, err := strconv.Atoi(la[4])
		require.Nil(t, err)
		require.Equal(t, "0.94", la[0])
		require.Equal(t, "1.49", la[1])
		require.Equal(t, "1.63", la[2])
		require.Equal(t, 1, procsRun)
		require.Equal(t, 924, procsTotal)
	})
	t.Run("wrong file", func(t *testing.T) {
		_, err := readLoadAvgFile("../../../tests/testdata/bad_file.txt")
		require.NotNil(t, err)
	})
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := readLoadAvgFile("loadavg.txt")
		require.NotNil(t, err)
	})
}

func TestDiskData(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		mb := `Filesystem            1K-blocks    Used Available Use% Mounted on
main/containers/proxy  12043904 1202560  10841344  10% /
none                        492       4       488   1% /dev
udev                     488476       0    488476   0% /dev/tty
tmpfs                       100       0       100   0% /dev/lxd
tmpfs                       100       0       100   0% /dev/.lxd-mounts
tmpfs                    507972      12    507960   1% /dev/shm
tmpfs                    507972   51068    456904  11% /run
tmpfs                      5120       0      5120   0% /run/lock
tmpfs                    507972       0    507972   0% /sys/fs/cgroup
tmpfs                    101596       0    101596   0% /run/user/0`

		inode := `Filesystem              Inodes  IUsed    IFree IUse% Mounted on
main/containers/proxy 21809084 126349 21682735    1% /
none                    126993     26   126967    1% /dev
udev                    122119    623   121496    1% /dev/tty
tmpfs                   126993      1   126992    1% /dev/lxd
tmpfs                   126993      4   126989    1% /dev/.lxd-mounts
tmpfs                   126993      4   126989    1% /dev/shm
tmpfs                   126993    179   126814    1% /run
tmpfs                   126993      3   126990    1% /run/lock
tmpfs                   126993     16   126977    1% /sys/fs/cgroup
tmpfs                   126993      5   126988    1% /run/user/0`
		data, err := parseDiskData(mb, inode)
		require.Nil(t, err)
		require.Equal(t, 10, len(data.Data))
		require.Equal(t, "main/containers/proxy", data.Data[0].Filesystem)
		require.Equal(t, int64(11), data.Data[6].Used)
		require.Equal(t, int64(1), data.Data[3].Inode)
	})

	t.Run("empty data", func(t *testing.T) {
		_, err := parseDiskData("", "")
		require.NotNil(t, err)
	})
}

func TestNetworkStats(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := `State                        Recv-Q                    Send-Q                                         Local Address:Port                                                Peer Address:Port                     Process                    
LISTEN                       0                         5                                                  127.0.0.1:dey-sapi                                                 0.0.0.0:*                                                   
LISTEN                       0                         5                                                    0.0.0.0:cvsup                                                    0.0.0.0:*                                                   
LISTEN                       0                         32                                             192.168.122.1:domain                                                   0.0.0.0:*                                                   
LISTEN                       0                         128                                                  0.0.0.0:ssh                                                      0.0.0.0:*                                                   
LISTEN                       0                         5                                                  127.0.0.1:ipp                                                      0.0.0.0:*                                                   
ESTAB                        0                         0                                              192.168.0.1:8882                                             123.123.23.213:https                                               
ESTAB                        0                         0                                              192.168.0.1:8884                                               123.232.23.23:https                                               
ESTAB                        0                         0                                              192.168.0.1:8885                                               111.16.10.14:https                                               
ESTAB                        0                         0                                              192.168.0.1:8886                                                232.75.122.24:https                                               
ESTAB                        0                         0                                              192.168.0.1:8887                                             151.101.65.140:https                                               
ESTAB                        0                         0                                              192.168.0.1:8888                                             173.194.73.113:https                                               
ESTAB                        0                         0                                              192.168.0.1:8889                                            123.101.245.140:https                                               
TIME-WAIT                    0                         0                                              192.168.0.1:9090                                              123.176.176.76:https                                               
ESTAB                        0                         0                                              192.168.0.1:9091                                            123.143.432.97:https`
		res := parseTCPConnections(data)
		require.Equal(t, 6, len(res))
		require.Equal(t, int64(5), res["LISTEN"])
		require.Equal(t, int64(1), res["TIME-WAIT"])
		require.Equal(t, int64(8), res["ESTAB"])
	})
}
