package linux

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var errCannotParseLoadAvg = errors.New("cannot parse loadavg file")

func (c *Collector) getLoadAvg(wg *sync.WaitGroup, ss *Snapshot) error {
	loadAvg := new(api.LoadAvg)
	defer wg.Done()
	l, err := readLoadAvgFile("/proc/loadavg")
	if err != nil {
		return err
	}
	loadAvg.OneMin, err = strconv.ParseFloat(l[0], 32)
	if err != nil {
		return err
	}
	loadAvg.FiveMin, err = strconv.ParseFloat(l[1], 32)
	if err != nil {
		return err
	}
	loadAvg.FifteenMin, err = strconv.ParseFloat(l[2], 32)
	if err != nil {
		return err
	}
	loadAvg.ProcsRunning, err = strconv.ParseInt(l[3], 10, 64)
	if err != nil {
		return err
	}
	loadAvg.TotalProcs, err = strconv.ParseInt(l[4], 10, 64)
	if err != nil {
		return err
	}
	ss.LoadAvg = loadAvg

	return nil
}

func readLoadAvgFile(fname string) (res [5]string, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return res, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		return res, err
	}

	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 5 {
		return res, errCannotParseLoadAvg
	}
	procs := strings.Split(fields[3], "/")
	if len(procs) != 2 {
		return res, errCannotParseLoadAvg
	}
	res = [5]string{fields[0], fields[1], fields[2], procs[0], procs[1]}

	return res, nil
}
