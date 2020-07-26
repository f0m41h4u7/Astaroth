package linux

import (
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
		require.Nil(t, err)
		require.Equal(t, "0.94", la[0])
		require.Equal(t, "1.49", la[1])
		require.Equal(t, "1.63", la[2])
		require.Equal(t, "1", la[3])
		require.Equal(t, "924", la[4])
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
