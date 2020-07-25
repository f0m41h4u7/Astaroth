package linux

import (
	"bufio"
	"os"
	"strconv"
	"time"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var (
	user   = float32(0.0)
	system = float32(0.0)
)

func GetCPU() (cpu *api.CPU, err error) {
	cpu.User, err = calculateCPU(user, "/sys/fs/cgroup/cpu/cpuacct.usage_user")
	if err != nil {
		return nil, err
	}
	cpu.System, err = calculateCPU(system, "/sys/fs/cgroup/cpu/cpuacct.usage_sys")
	if err != nil {
		return nil, err
	}
	return
}

func calculateCPU(prev float32, fname string) (float32, error) {
	tstart := time.Now()
	tmp, err := readCPUFile(fname)
	if err != nil {
		return prev, err
	}
	cstart := float32(tmp)
	time.Sleep(100 * time.Millisecond)
	tmp, err = readCPUFile(fname)
	if err != nil {
		return prev, err
	}
	tstop := time.Now()
	cstop := float32(tmp)
	if cstop > cstart {
		duration := float32(tstop.Sub(tstart).Nanoseconds())
		if prev == 0 {
			prev = (cstop - cstart) / duration * 100.0
		} else {
			prev = 0.8*prev + 0.2*((cstop-cstart)/duration*100.0)
		}
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
