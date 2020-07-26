package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	On  IsIncluded = "on"
	Off IsIncluded = "off"

	LoadAvg      MetricsType = "load_avg"
	CPU          MetricsType = "cpu_usage"
	DiskUsage    MetricsType = "disk_usage"
	DiskData     MetricsType = "disk_data"
	TopTalkers   MetricsType = "top_talkers"
	NetworkStats MetricsType = "network_stats"
)

var (
	RequiredMetrics Config

	ErrCannotReadConfig  = errors.New("cannot read config file")
	ErrCannotParseConfig = errors.New("cannot parse config file")
)

type (
	IsIncluded  string
	MetricsType string
)

type Config struct {
	Metrics map[MetricsType]IsIncluded `json:"metrics"`
}

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
		return ErrCannotReadConfig
	}
	err = json.Unmarshal(conf, &RequiredMetrics)
	if err != nil {
		return ErrCannotParseConfig
	}
	return nil
}
