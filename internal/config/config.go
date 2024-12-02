package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GoatCounter struct {
		Enabled  bool
		SiteCode string
		APIToken string
	}
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)

	v.SetEnvPrefix("avalanche")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config

	err = v.Unmarshal(&config)
	return &config, err
}
