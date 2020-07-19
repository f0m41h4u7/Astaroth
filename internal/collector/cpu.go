package collector

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func getCPU() ([]float64, error) {
	cpu, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	content := string(cpu)
	fields := strings.Fields(strings.Split(content, "\n")[0])
	res := make([]float64, 3)

	// %user_mode
	res[0], err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return res, err
	}

	// %system_mode
	res[1], err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return res, err
	}

	// %idle
	res[2], err = strconv.ParseFloat(fields[4], 64)
	return res, err
}
