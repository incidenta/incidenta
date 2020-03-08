package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Logging Logging
		Server  Server
	}

	Logging struct {
		LogLevel string `envconfig:"LOGS_LEVEL" default:"info"`
		Color    bool   `envconfig:"LOGS_COLOR" default:"true"`
		Text     bool   `envconfig:"LOGS_TEXT" default:"true"`
		Pretty   bool   `envconfig:"LOGS_PRETTY"`
	}

	Server struct {
		Addr string `envconfig:"SERVER_ADDR" default:":8080"`
	}
)

func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("incidenta", &cfg)
	return cfg, err
}

func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}
