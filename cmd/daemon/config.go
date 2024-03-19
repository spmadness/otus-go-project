package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Metrics MetricsConf
}

type MetricsConf struct {
	LoadAverageSystem bool `mapstructure:"system"`
	LoadAverageCPU    bool `mapstructure:"cpu"`
	LoadDisks         bool `mapstructure:"disk"`
	NetTopTalkers     bool `mapstructure:"netTopTalkers"`
	NetStatistics     bool `mapstructure:"netStats"`
}

func NewConfig(configFile string) Config {
	var config Config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config read error: %s", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("config unmarshal error: %s", err)
		os.Exit(1)
	}

	return config
}
