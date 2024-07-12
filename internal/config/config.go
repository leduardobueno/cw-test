package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LogFile      *LogFileConfig
	MeansOfDeath []string
}

func FromConfigFile() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("/etc/cw-test/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("error parsing config data: %w", err))
	}

	return conf
}
