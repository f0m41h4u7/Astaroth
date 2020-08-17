package collector

import (
	"strconv"
	"testing"
	"fmt"

	"github.com/stretchr/testify/require"
)

func TestCPU(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cpu, err := readCPUFile("../../tests/testdata/cpu.txt")
		require.Nil(t, err)
		require.Equal(t, 42055566250653, cpu)
	})
	t.Run("wrong file", func(t *testing.T) {
		_, err := readCPUFile("../../tests/testdata/bad_file.txt")
		require.NotNil(t, err)
	})
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := readCPUFile("cpu.txt")
		require.NotNil(t, err)
	})
}

func TestLoadAvg(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		la, err := readLoadAvgFile("../../tests/testdata/loadavg.txt")
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
		_, err := readLoadAvgFile("../../tests/testdata/bad_file.txt")
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
/*
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
*/
func TestTopTalkers(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		data := `     ARP, Request who-has 192.168.0.110 tell 192.168.0.1, length 46
      IP 192.168.0.1.51891 > 255.255.255.255.7437: UDP, length 173
      IP 192.168.0.112.50858 > 192.168.0.255.137: UDP, length 50
      IP 192.168.88.1.36963 > 239.255.255.250.1900: UDP, length 167
      IP 192.168.0.107.49891 > 239.255.255.250.1900: UDP, length 167
      ARP, Request who-has 192.168.0.100 tell 192.168.0.1, length 46
      IP 192.168.0.107.34300 > 68.232.34.200.443: Flags [.], ack 1, win 501, length 0
      IP 68.232.34.200.443 > 192.168.0.107.34300: Flags [.], ack 1, win 136, length 0
      IP 68.232.34.200.443 > 192.168.0.107.34300: Flags [.], ack 1296915859, win 136, length 0`

		_, err := parseTopTalkers(data)
		require.Nil(t, err)
	})
}

func TestWindows(t *testing.T) {
	t.Run("simple", func(t *testing.T){
		data := `
Активные подключения

  Имя    Локальный адрес        Внешний адрес          Состояние       PID
  TCP    0.0.0.0:7658           0.0.0.0:0              LISTENING       4
  TCP    0.0.0.0:1234           0.0.0.0:0              LISTENING       1004
  TCP    0.0.0.0:8768           0.0.0.0:0              LISTENING       856
  TCP    127.0.0.1:4321         127.0.0.1:23245        ESTABLISHED     4080
  TCP    127.0.0.1:5432         127.0.0.1:9876         TIME_WAIT       0
  TCP    [::]:3456              [::]:0                 LISTENING       3236
  TCP    127.0.0.1:11343        127.0.0.1:12309        ESTABLISHED     11332
  TCP    127.0.0.1:6543         0.0.0.0:0              LISTENING       12200
  TCP    [::]:23456             [::]:0                 LISTENING       2484
  TCP    [::]:34567             [::]:0                 LISTENING       972
  TCP    127.0.0.1:45678        12.45.23.0:443         CLOSE_WAIT      12460`
	
	procs := `
Имя образа                     PID Имя сессии          № сеанса       Память
========================= ======== ================ =========== ============
System Idle Process              0 Services                   0         8 КБ
System                           4 Services                   0     1 308 КБ
wininit.exe                    856 Services                   0     6 532 КБ
services.exe                   972 Services                   0    10 532 КБ
lsass.exe                     1004 Services                   0    24 032 КБ
CxAudMsg64.exe                4072 Services                   0     8 960 КБ
mDNSResponder.exe             4080 Services                   0     6 736 КБ
OneApp.IGCC.WinService.ex     3236 Services                   0    37 184 КБ
svchost.exe                   6676 Services                   0     6 372 КБ
ctfmon.exe                    7872 Console                    1    15 160 КБ
steam.exe                    12200 Console                    1   216 984 КБ
SynTPEnh.exe                  7976 Console                    1    34 292 КБ
mediaget.exe                 11332 Console                    1    89 032 КБ
svchost.exe                   2484 Services                   0    17 788 КБ
EpicGamesLauncher.exe        12460 Console                    1   147 596 КБ
UnrealCEFSubProcess.exe      13096 Console                    1    33 764 КБ
ShellExperienceHost.exe       3468 Console                    1    61 532 КБ
tasklist.exe                 13768 Console                    1     8 688 КБ
`
		res, err := parseStates(data)
		require.Nil(t, err)
		require.Equal(t, 6, len(res))
		require.Equal(t, int64(1), res["TIME-WAIT"])
		require.Equal(t, int64(2), res["ESTAB"])
		require.Equal(t, int64(7), res["LISTEN"])
		require.Equal(t, int64(1), res["CLOSE-WAIT"])
		
		lis, err := parseListenSockets(data, procs)
		require.Nil(t, err)
		fmt.Println(lis)
	})
}
