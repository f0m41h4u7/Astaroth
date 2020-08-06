package linux

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

func (c *Collector) getNetworkStats(wg *sync.WaitGroup, snap *Snapshot) error {
	defer wg.Done()

	ss, err := exec.Command("ss", "-ta").Output()
	if err != nil {
		return err
	}
	states := parseTCPConnections(string(ss))
	ns := new(api.NetworkStats)
	ns.TCPConnStates = &api.States{
		LISTEN:     states["LISTEN"],
		ESTAB:      states["ESTAB"],
		FIN_WAIT:   states["FIN-WAIT"],
		SYN_RCV:    states["SYN-RCV"],
		TIME_WAIT:  states["TIME-WAIT"],
		CLOSE_WAIT: states["CLOSE-WAIT"],
	}

	netstat, err := exec.Command("netstat", "-lntupe").Output()
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

func parseTCPConnections(data string) map[string]int64 {
	lines := strings.Split(data, "\n")[1:]
	numStates := newMap()
	for _, line := range lines {
		state := strings.Fields(strings.TrimSpace(line))[0]
		if _, ok := numStates[state]; !ok {
			numStates[state] = int64(1)
		} else {
			numStates[state]++
		}
	}

	return numStates
}

func newMap() map[string]int64 {
	states := map[string]int64{}
	states["LISTEN"] = 0
	states["ESTAB"] = 0
	states["FIN-WAIT"] = 0
	states["SYN-RCV"] = 0
	states["TIME-WAIT"] = 0
	states["CLOSE-WAIT"] = 0

	return states
}

func parseSockets(data string) ([]*api.Sockets, error) {
	lines := strings.Split(data, "\n")[2:]
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

		for _, f := range fields {
			if strings.Contains(f, "/") {
				pid, err := strconv.ParseInt(strings.Split(f, "/")[0], 10, 64)
				if err != nil {
					return nil, err
				}
				sockets[i].PID = pid
				sockets[i].Program = strings.Split(f, "/")[1]

				break
			}
		}
	}

	return sockets, nil
}
