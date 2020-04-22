package command

import (
	"github.com/sirupsen/logrus"

	"github.com/incidenta/incidenta/pkg/config"
)

func initLogging(c config.Config) {
	if level, err := logrus.ParseLevel(c.Logging.LogLevel); err == nil {
		logrus.SetLevel(level)
	}
	if c.Logging.Text {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   c.Logging.Color,
			DisableColors: !c.Logging.Color,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: c.Logging.Pretty,
		})
	}
}
