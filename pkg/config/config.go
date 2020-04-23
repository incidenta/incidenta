package config

import (
	"encoding/json"

	"github.com/kelseyhightower/envconfig"

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
	out, _ := json.Marshal(c)
	return string(out)
}
