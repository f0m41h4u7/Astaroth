package linux

import (
	"bufio"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var ErrWrongTopTalkersData = errors.New("cannot parse top talkers data")

func (c *Collector) getTopTalkers(wg *sync.WaitGroup, ss *Snapshot) error {
	defer wg.Done()

	tcpdump := "tcpdump -tnn -c 40 -i any -Q inout -l | uniq | head"
	cmd := exec.Command("bash", "-c", tcpdump)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)

	out := ""
	for scanner.Scan() {
		out = out + scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	ss.TopTalkers, err = parseTopTalkers(out)

	return err
}

func parseTopTalkers(str string) (*api.TopTalkers, error) {
	tt := api.TopTalkers{
		ByProtocol: []*api.ByProtocol{},
		ByTraffic:  []*api.ByTraffic{},
	}
	if str == "" {
		return nil, ErrWrongTopTalkersData
	}

	lines := strings.Split(str, "\n")
	protocols := map[string]int64{}
	var sumBytes int64
	sumBytes = 0

	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		if fields[len(fields)-2] != "length" {
			continue
		}
		bytes, err := strconv.ParseInt(fields[len(fields)-1], 10, 64)
		if err != nil {
			return nil, ErrWrongTopTalkersData
		}

		if _, ok := protocols[fields[0]]; !ok {
			protocols[fields[0]] = bytes
		} else {
			protocols[fields[0]] += bytes
		}
		sumBytes += bytes

		tt.ByTraffic = append(tt.ByTraffic, &api.ByTraffic{
			SourceAddr: fields[1],
			DestAddr:   fields[3],
			Protocol:   fields[0],
			Bps:        bytes,
		})
	}

	for protocol, bytes := range protocols {
		tt.ByProtocol = append(tt.ByProtocol, &api.ByProtocol{
			Protocol:   protocol,
			Bytes:      bytes,
			Percentage: 100 * bytes / sumBytes,
		})
	}

	return &tt, nil
}
