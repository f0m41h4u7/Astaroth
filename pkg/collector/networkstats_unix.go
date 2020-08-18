// +build !windows

package collector

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var ErrWrongData = errors.New("cannot parse tcp connections data")

func (c *Collector) getNetworkStats(wg *sync.WaitGroup, snap *Snapshot) error {
	defer wg.Done()

	ss, err := ioutil.ReadFile("/proc/net/nf_conntrack")
	if err != nil {
		return err
	}
	states, err := parseTCPConnections(string(ss))
	if err != nil {
		return err
	}

	ns := new(api.NetworkStats)
	ns.TCPConnStates = &api.States{
		LISTEN:     states["LISTEN"],
		ESTAB:      states["ESTAB"],
		FIN_WAIT:   states["FIN-WAIT"],
		SYN_RCV:    states["SYN-RCV"],
		TIME_WAIT:  states["TIME-WAIT"],
		CLOSE_WAIT: states["CLOSE-WAIT"],
	}

	cmd := "netstat -lntupe | grep LISTEN"
	netstat, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	ns.ListenSockets, err = parseSockets(string(netstat))
	if err != nil {
		return err
	}
	snap.NetworkStats = ns

	return nil
}

func parseTCPConnections(data string) (map[string]int64, error) {
	numStates := newMap()
	if data == "" {
		return numStates, ErrWrongData
	}

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "UNREPLIED"):
			continue
		case strings.Contains(line, "LISTEN"):
			numStates["LISTEN"]++
		case strings.Contains(line, "ESTABLISHED"):
			numStates["ESTAB"]++
		case strings.Contains(line, "FIN_WAIT"):
			numStates["FIN-WAIT"]++
		case strings.Contains(line, "SYN_RECV"):
			numStates["SYN-RCV"]++
		case strings.Contains(line, "TIME_WAIT"):
			numStates["TIME-WAIT"]++
		case strings.Contains(line, "CLOSE"):
			numStates["CLOSE-WAIT"]++
		}
	}

	return numStates, nil
}

func parseSockets(data string) ([]*api.Sockets, error) {
	if data == "" {
		return nil, ErrWrongData
	}
	lines := strings.Split(data, "\n")
	sockets := make([]*api.Sockets, len(lines))

	for i := 0; i < len(lines); i++ {
		fields := strings.Fields(strings.TrimSpace(lines[i]))
		if len(fields) == 0 {
			continue
		}
		sockets[i] = new(api.Sockets)

		sockets[i].Protocol = fields[0]
		tmp := strings.Split(fields[3], ":")
		port, err := strconv.ParseInt(tmp[len(tmp)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		sockets[i].Port = port

		sockets[i].User = fields[6]

		pid, err := strconv.ParseInt(strings.Split(fields[8], "/")[0], 10, 64)
		if err != nil {
			return nil, err
		}
		sockets[i].PID = pid
		sockets[i].Program = strings.Split(fields[8], "/")[1]
	}

	return sockets, nil
}
