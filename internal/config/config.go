package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type (
	IsIncluded  string
	MetricsType string
	Metrics     map[MetricsType]IsIncluded
)

const (
	On  IsIncluded = "on"
	Off IsIncluded = "off"

	LoadAvg      MetricsType = "load_avg"
	CPU          MetricsType = "cpu_usage"
	DiskUsage    MetricsType = "disk_usage"
	TopTalkers   MetricsType = "top_talkers"
	NetworkStats MetricsType = "network_stats"
)

var (
	RequiredMetrics Config

	ErrCannotReadConfig  = errors.New("cannot read config file")
	ErrCannotParseConfig = errors.New("cannot parse config file")
)

type Config struct {
	Metrics Metrics
}

func InitConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.SetConfigName("configs/config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return ErrCannotReadConfig
	}

	if err := viper.Unmarshal(&RequiredMetrics); err != nil {
		return ErrCannotParseConfig
	}
	return nil
}
