package models

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	x      *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables,
		new(Alert),
		new(Log),
		new(Receiver),
		new(Snooze),
		new(Template),
	)
}

func getEngine(config Config) (*xorm.Engine, error) {
	connStr, err := config.GetConnString()
	if err != nil {
		return nil, err
	}
	engine, err := xorm.NewEngine(config.Type, connStr)
	if err != nil {
		return nil, err
	}
	engine.SetSchema(config.Schema)
	return engine, nil
}

func SetEngine(config Config) (err error) {
	x, err = getEngine(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger := NewXormLogger(logrus.WarnLevel, logrus.WithField("source", "sql"))
	x.SetLogger(logger)

	if config.LogQueries {
		x.ShowSQL(true)
	}

	x.SetMapper(core.GonicMapper{})
	x.SetMaxOpenConns(config.MaxOpenConn)
	x.SetMaxIdleConns(config.MaxIdleConn)
	x.SetConnMaxLifetime(config.ConnMaxLifetime)

	return nil
}

func NewEngine(ctx context.Context, config Config) error {
	if err := SetEngine(config); err != nil {
		return err
	}

	x.SetDefaultContext(ctx)

	if err := x.Ping(); err != nil {
		return err
	}

	if err := x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v", err)
	}

	return nil
}

func Ping() error {
	if x != nil {
		return x.Ping()
	}
	return errors.New("database not configured")
}

func Close() error {
	if x != nil {
		return x.Close()
	}
	return nil
}
