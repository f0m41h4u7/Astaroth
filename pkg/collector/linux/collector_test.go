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
		require.Nil(t, err)
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

	t.Run("no inode", func(t *testing.T) {
		mb := `Filesystem     1K-blocks      Used Available Use% Mounted on
devtmpfs         8057084         0   8057084   0% /dev
tmpfs            8076408      1776   8074632   1% /run
tmpfs            8076408         0   8076408   0% /sys/fs/cgroup
/dev/nvme0n1p5 422120088 230446280 170161616  58% /
/dev/loop0         56320     56320         0 100% /var/lib/snapd/snap/core18/1880
/dev/nvme0n1p2    306584     79256    227328  26% /boot/efi
/dev/loop6        165376    165376         0 100% /var/lib/snapd/snap/gnome-3-28-1804/128`
		inode := `Filesystem       Inodes   IUsed    IFree IUse% Mounted on
devtmpfs        2014271     591  2013680    1% /dev
tmpfs           2019102    1111  2017991    1% /run
tmpfs           2019102      17  2019085    1% /sys/fs/cgroup
/dev/nvme0n1p5 26869760 1396968 25472792    6% /
/dev/loop0        10756   10756        0  100% /var/lib/snapd/snap/core18/1880
/dev/nvme0n1p2        0       0        0     - /boot/efi
/dev/loop6        27798   27798        0  100% /var/lib/snapd/snap/gnome-3-28-1804/128`
		_, err := parseDiskData(mb, inode)
		require.Nil(t, err)
	})
}

func TestNetworkStats(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := `ipv4     2 tcp      6 118 TIME_WAIT src=192.168.0.103 dst=140.82.118.6 sport=56346 dport=443 src=140.82.118.6 dst=192.168.0.103 sport=443 dport=56346 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 431962 ESTABLISHED src=192.168.0.103 dst=151.101.84.133 sport=39822 dport=443 src=151.101.84.133 dst=192.168.0.103 sport=443 dport=39822 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 0 TIME_WAIT src=192.168.0.103 dst=213.36.253.2 sport=59174 dport=80 src=213.36.253.2 dst=192.168.0.103 sport=80 dport=59174 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 unknown  2 255 src=192.168.0.1 dst=224.0.0.251 [UNREPLIED] src=224.0.0.251 dst=192.168.0.1 mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 431975 ESTABLISHED src=192.168.0.103 dst=151.101.84.133 sport=39720 dport=443 src=151.101.84.133 dst=192.168.0.103 sport=443 dport=39720 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 431985 ESTABLISHED src=192.168.0.103 dst=151.101.84.133 sport=39836 dport=443 src=151.101.84.133 dst=192.168.0.103 sport=443 dport=39836 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 4 TIME_WAIT src=192.168.0.103 dst=74.125.205.147 sport=41794 dport=80 src=74.125.205.147 dst=192.168.0.103 sport=80 dport=41794 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2
ipv4     2 tcp      6 4 TIME_WAIT src=192.168.0.103 dst=213.36.253.2 sport=59190 dport=80 src=213.36.253.2 dst=192.168.0.103 sport=80 dport=59190 [ASSURED] mark=0 secctx=system_u:object_r:unlabeled_t:s0 zone=0 use=2`
		res, err := parseTCPConnections(data)
		require.Nil(t, err)
		require.Equal(t, 6, len(res))
		require.Equal(t, int64(4), res["TIME-WAIT"])
		require.Equal(t, int64(3), res["ESTAB"])
	})

	t.Run("empty data", func(t *testing.T) {
		data := ""
		_, err := parseTCPConnections(data)
		require.NotNil(t, err)
	})

	t.Run("netstat", func(t *testing.T) {
		netstat := `tcp        0      0 127.0.0.1:1234          0.0.0.0:*               LISTEN      0          35375      1896/pmcd           
tcp        0      0 0.0.0.0:1337            0.0.0.0:*               LISTEN      0          37433      999/vncserver-virtu 
tcp        0      0 127.0.0.1:666           0.0.0.0:*               LISTEN      0          43739      9092/dnsmasq        
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      0          39249      9783/sshd            
tcp        0      0 127.0.0.1:9090          0.0.0.0:*               LISTEN      0          36389      1345/cupsd`
		ns, err := parseSockets(netstat)
		require.Nil(t, err)
		require.Equal(t, int64(1896), ns[0].PID)
		require.Equal(t, int64(666), ns[2].Port)
		require.Equal(t, "tcp", ns[3].Protocol)
		require.Equal(t, "vncserver-virtu", ns[1].Program)
	})

	t.Run("empty data", func(t *testing.T) {
		data := ""
		_, err := parseSockets(data)
		require.NotNil(t, err)
	})
}
