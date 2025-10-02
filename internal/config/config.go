package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	Version  string `env:"SVC_VERSION" envDefault:"unknown"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
