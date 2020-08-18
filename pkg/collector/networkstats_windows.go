// +build windows
package collector

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var ErrWrongData = errors.New("cannot parse tcp connections data")

func (c *Collector) getNetworkStats(wg *sync.WaitGroup, snap *Snapshot) error {
	defer wg.Done()

	netstat, err := exec.Command("cmd", "/C", "netstat", "-aon").Output()
	if err != nil {
		return err
	}
	procs, err := exec.Command("cmd", "/C", "tasklist").Output()
	if err != nil {
		return err
	}

	states, err := parseStates(string(netstat))
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

	ns.ListenSockets, err = parseListenSockets(string(netstat), string(procs))
	if err != nil {
		return err
	}
	snap.NetworkStats = ns

	return nil
}

func parseStates(data string) (map[string]int64, error) {
	numStates := newMap()
	if data == "" {
		return numStates, ErrWrongData
	}

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "LISTENING"):
			numStates["LISTEN"]++
		case strings.Contains(line, "ESTABLISHED"):
			numStates["ESTAB"]++
		case strings.Contains(line, "FIN_WAIT"):
			numStates["FIN-WAIT"]++
		case strings.Contains(line, "SYN_RECV"):
			numStates["SYN-RCV"]++
		case strings.Contains(line, "TIME_WAIT"):
			numStates["TIME-WAIT"]++
		case strings.Contains(line, "CLOSE_WAIT"):
			numStates["CLOSE-WAIT"]++
		default:
			continue
		}
	}

	return numStates, nil
}

func parseProcs(procs string) map[int64]string {
	procsMap := map[int64]string{}
	lines := strings.Split(procs, "\n")[3:]
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		pid, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		procsMap[pid] = fields[0]
	}
	return procsMap
}

func parseListenSockets(ns string, procs string) ([]*api.Sockets, error) {
	if (ns == "") || (procs == "") {
		return nil, ErrWrongData
	}
	lines := strings.Split(ns, "\n")[3:]
	sockets := []*api.Sockets{}
	ps := parseProcs(procs)

	for i := 0; i < len(lines); i++ {
		fields := strings.Fields(strings.TrimSpace(lines[i]))
		if len(fields) == 0 {
			continue
		}
		if fields[3] == "LISTENING" {
			s := new(api.Sockets)
			s.Protocol = fields[0]

			tmp := strings.Split(fields[1], ":")
			port, err := strconv.ParseInt(tmp[len(tmp)-1], 10, 64)
			if err != nil {
				return nil, err
			}
			s.Port = port
			s.User = ""
			s.PID, err = strconv.ParseInt(fields[4], 10, 64)
			if err != nil {
				return nil, err
			}
			var ok bool
			s.Program, ok = ps[s.PID]
			if !ok {
				return nil, ErrWrongData
			}
			sockets = append(sockets, s)
		}
	}

	return sockets, nil
}
