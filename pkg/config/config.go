package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	"github.com/incidenta/incidenta/pkg/models"
	"github.com/incidenta/incidenta/pkg/server"
)

type Config struct {
	Logging  Logging
	Server   server.Config
	Database models.Config
}

func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("incidenta", &cfg)
	return cfg, err
}

func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}
