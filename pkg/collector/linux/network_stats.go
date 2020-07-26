package linux

import (
	"os/exec"
	"strings"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

func GetNetworkStats() (*api.NetworkStats, error) {
	ss, err := exec.Command("ss", "-ta").Output()
	if err != nil {
		return nil, err
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
	return ns, nil
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
