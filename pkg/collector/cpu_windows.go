// +build windows

package collector

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

func (c *CPUMetric) Get(wg *sync.WaitGroup) error {
	defer wg.Done()

	data, err := exec.Command("cmd", "/C", "wmic", "cpu", "get", "loadpercentage", "/value").Output()
	if err != nil {
		return err
	}
	num := parseCPU(string(data))

	cpu := new(api.CPU)
	cpu.User, err = strconv.ParseFloat(num, 64)
	if err != nil {
		return err
	}
	c.CPU = cpu

	return nil
}

func parseCPU(data string) string {
	return strings.Split(strings.TrimSpace(data), "=")[1]
}
