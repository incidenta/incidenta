package models

import (
	"errors"
	"fmt"
	"time"
)

type Config struct {
	Type            string        `envconfig:"DB_TYPE" default:"sqlite3"`
	Schema          string        `envconfig:"DB_SCHEMA" default:"public"`
	Path            string        `envconfig:"DB_PATH" default:"database.db"`
	MaxIdleConn     int           `envconfig:"DB_MAX_IDLE_CONN" default:"2"`
	MaxOpenConn     int           `envconfig:"DB_MAX_OPEN_CONN" default:"0"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"0"`
	LogQueries      bool          `envconfig:"DB_LOG_QUERIES" default:"false"`
}

func (c *Config) GetConnString() (string, error) {
	switch c.Type {
	case "sqlite3":
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc&_busy_timeout=500&_txlock=immediate", c.Path), nil
	}
	return "", errors.New("unknown database type")
}
