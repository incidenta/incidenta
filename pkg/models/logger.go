package models

import (
	"strings"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm/log"
)

type XormLogger struct {
	showSQL bool
	level   logrus.Level
	log     *logrus.Entry
}

func trimFormat(format string) string {
	return strings.TrimLeft(format, "[SQL] ")
}

func NewXormLogger(level logrus.Level, logger *logrus.Entry) log.Logger {
	return &XormLogger{
		level: level,
		log:   logger,
	}
}

func (s *XormLogger) SetLevel(l log.LogLevel) {}

func (s *XormLogger) Error(v ...interface{}) {
	s.log.Error(v...)
}

func (s *XormLogger) Errorf(format string, v ...interface{}) {
	s.log.Errorf(trimFormat(format), v...)
}

func (s *XormLogger) Debug(v ...interface{}) {
	s.log.Debug(v...)
}

func (s *XormLogger) Debugf(format string, v ...interface{}) {
	s.log.Debugf(trimFormat(format), v...)
}

func (s *XormLogger) Info(v ...interface{}) {
	s.log.Info(v...)
}

func (s *XormLogger) Infof(format string, v ...interface{}) {
	s.log.Infof(trimFormat(format), v...)
}

func (s *XormLogger) Warn(v ...interface{}) {
	s.log.Warn(v...)
}

func (s *XormLogger) Warnf(format string, v ...interface{}) {
	s.log.Warnf(trimFormat(format), v...)
}

func (s *XormLogger) Level() log.LogLevel {
	switch s.level {
	case logrus.ErrorLevel:
		return log.LOG_ERR
	case logrus.WarnLevel:
		return log.LOG_WARNING
	case logrus.InfoLevel:
		return log.LOG_INFO
	case logrus.DebugLevel:
		return log.LOG_DEBUG
	default:
		return log.LOG_ERR
	}
}

func (s *XormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

func (s *XormLogger) IsShowSQL() bool {
	return s.showSQL
}
