package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/f0m41h4u7/Astaroth/pkg/collector"
)

// Available metrics.
const (
	On  isIncluded = "on"
	Off isIncluded = "off"

	LoadAvg      metricsType = "load_avg"
	CPU          metricsType = "cpu_usage"
	DiskData     metricsType = "disk_data"
	TopTalkers   metricsType = "top_talkers"
	NetworkStats metricsType = "network_stats"
)

var (
	errCannotReadConfig  = errors.New("cannot read config file")
	errCannotParseConfig = errors.New("cannot parse config file")

	metricsIfaces = map[metricsType]collector.Metric{
		LoadAvg:      &collector.LoadAvgMetric{LoadAvg: new(api.LoadAvg)},
		CPU:          &collector.CPUMetric{CPU: new(api.CPU)},
		DiskData:     &collector.DiskDataMetric{DiskData: new(api.DiskData)},
		NetworkStats: &collector.NetworkStatsMetric{NetworkStats: new(api.NetworkStats)},
		TopTalkers:   &collector.TopTalkersMetric{TopTalkers: new(api.TopTalkers)},
	}
)

type (
	isIncluded  string
	metricsType string
)

// Config states which metrics are required.
type Config struct {
	Metrics map[metricsType]isIncluded `json:"metrics"`
}

// ReadConfig returns array of required metrics.
func ReadConfig(cfgFile string) ([]collector.Metric, error) {
	required := Config{}
	if cfgFile == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		cfgFile = cwd + "/configs/config.json"
	}
	conf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errCannotReadConfig
	}
	err = json.Unmarshal(conf, &required)
	if err != nil {
		return nil, errCannotParseConfig
	}

	metrics := []collector.Metric{}
	for mt, status := range required.Metrics {
		if status == On {
			metrics = append(metrics, metricsIfaces[mt])
		}
	}

	return metrics, nil
}
