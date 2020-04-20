package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/incidenta/incidenta/pkg/config"
	"github.com/incidenta/incidenta/pkg/models"
	"github.com/incidenta/incidenta/pkg/server"
)

var (
	envFile = flag.String("env-file", ".env", "Read in a file of environment variables")
)

func main() {
	flag.Parse()

	_ = godotenv.Load(*envFile)

	cfg, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	initLogging(cfg)

	logrus.Infof("Starting incidenta on %s...", cfg.Server.Addr)

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		fmt.Println(cfg.String())
	}

	if err := models.NewEngine(context.Background(), cfg.Database); err != nil {
		logrus.Fatal(err)
	}

	s := server.New(cfg.Server)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := s.Run()
		if err != nil {
			logrus.Fatal(err)
		}
		os.Exit(0)
	}()

	<-stop

	logrus.Info("Shutting down")

	models.Close()
}

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
