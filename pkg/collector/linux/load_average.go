package linux

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var errCannotParseLoadAvg = errors.New("cannot parse loadavg file")

// GetLoadavg collects load average
func GetLoadavg() (loadAvg *api.LoadAvg, err error) {
	l, err := readLoadAvgFile("/proc/loadavg")
	if err != nil {
		return
	}
	var tmp float64

	tmp, err = strconv.ParseFloat(l[0], 32)
	if err != nil {
		return
	}
	loadAvg.OneMin = float32(tmp)

	tmp, err = strconv.ParseFloat(l[1], 32)
	if err != nil {
		return
	}
	loadAvg.FiveMin = float32(tmp)

	tmp, err = strconv.ParseFloat(l[2], 32)
	if err != nil {
		return
	}
	loadAvg.FifteenMin = float32(tmp)

	loadAvg.ProcsRunning, err = strconv.ParseInt(l[3], 10, 64)
	if err != nil {
		return
	}

	loadAvg.TotalProcs, err = strconv.ParseInt(l[4], 10, 64)
	return
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
