// +build !windows

package collector

import (
	"bufio"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var (
	user   = 0.0
	system = 0.0
)

func (c *Collector) getCPU(wg *sync.WaitGroup, ss *Snapshot) error {
	defer wg.Done()

	cpu := new(api.CPU)
	var err error
	cpu.User, err = calculateCPU(user, "/sys/fs/cgroup/cpu/cpuacct.usage_user")
	if err != nil {
		return err
	}
	user = cpu.User
	cpu.System, err = calculateCPU(system, "/sys/fs/cgroup/cpu/cpuacct.usage_sys")
	if err != nil {
		return err
	}
	system = cpu.System
	ss.CPU = cpu

	return nil
}

func calculateCPU(prev float64, fname string) (float64, error) {
	tstart := time.Now().UnixNano()
	cstart, err := readCPUFile(fname)
	if err != nil {
		return prev, err
	}
	time.Sleep(10 * time.Millisecond)
	cstop, err := readCPUFile(fname)
	if err != nil {
		return prev, err
	}
	tstop := time.Now().UnixNano()
	if cstop > cstart {
		duration := tstop - tstart
		prev = float64(cstop-cstart) / float64(duration) * 100.0
	}

	return prev, nil
}

func readCPUFile(fname string) (int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return strconv.Atoi(line)
}
