package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	ServerPort          string `mapstructure:"SERVER_PORT"`
	ServerEnableTracing bool   //`mapstructure:"SERVER_ENABLE_TRACING"`
	ServerEnv           string `mapstructure:"SERVER_ENV"`
	ServerVersion       string `mapstructure:"SERVER_VERSION"`
	ServerName          string `mapstructure:"SERVER_NAME"`
	Build               string `mapstructure:"BUILD"`
	CollectorUrl        string `mapstructure:"COLLECTOR_URL"`
	ViaCepUrl           string `mapstructure:"VIACEP_API_URL"`
	WeatherApiUrl       string `mapstructure:"WEATHER_API_URL"`
	WeatherApiKey       string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfigs() (*config, error) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	keys := []string{
		"SERVER_PORT",
		"SERVER_ENABLE_TRACING",
		"SERVER_ENV",
		"SERVER_VERSION",
		"SERVER_NAME",
		"BUILD",
		"COLLECTOR_URL",
		"VIACEP_API_URL",
		"WEATHER_API_URL",
		"WEATHER_API_KEY",
	}

	// Setar defaults e bind
	for _, k := range keys {
		viper.SetDefault(k, "")
		if err := viper.BindEnv(k); err != nil {
			return nil, fmt.Errorf("bind error for %s: %w", k, err)
		}
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.ServerPort == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	return &cfg, nil
}
