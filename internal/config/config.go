package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// Available metrics.
const (
	On  isIncluded = "on"
	Off isIncluded = "off"

	LoadAvg      metricsType = "load_avg"
	CPU          metricsType = "cpu_usage"
	DiskUsage    metricsType = "disk_usage"
	DiskData     metricsType = "disk_data"
	TopTalkers   metricsType = "top_talkers"
	NetworkStats metricsType = "network_stats"
)

var (
	// RequiredMetrics is a global Config.
	RequiredMetrics Config

	errCannotReadConfig  = errors.New("cannot read config file")
	errCannotParseConfig = errors.New("cannot parse config file")
)

type (
	isIncluded  string
	metricsType string
)

// Config states which metrics are required.
type Config struct {
	Metrics map[metricsType]isIncluded `json:"metrics"`
}

// InitConfig initializes RequiredMetrics config.
func InitConfig(cfgFile string) error {
	if cfgFile == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cfgFile = cwd + "/configs/config.json"
	}
	conf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return errCannotReadConfig
	}
	err = json.Unmarshal(conf, &RequiredMetrics)
	if err != nil {
		return errCannotParseConfig
	}
	return nil
}
