package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/incidenta/incidenta/handler/api"
	"github.com/incidenta/incidenta/handler/health"

	"github.com/incidenta/incidenta/cmd/incidenta/config"
	"github.com/incidenta/incidenta/server"
	"github.com/incidenta/incidenta/signal"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	envFile = flag.String("env-file", ".env", "Read in a file of environment variables")
)

func main() {
	flag.Parse()

	godotenv.Load(*envFile)

	cfg, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	initLogging(cfg)
	ctx := signal.WithContext(context.Background())

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		fmt.Println(cfg.String())
	}

	s := server.Server{
		Addr:    cfg.Server.Addr,
		Handler: handler(),
	}

	g := errgroup.Group{}
	g.Go(func() error {
		logrus.WithFields(
			logrus.Fields{
				"url": cfg.Server.Addr,
			},
		).Infoln("starting the http server")
		return s.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatalln("program terminated")
	}
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api", api.Handler())
	mux.Handle("/health", health.Handler())
	return mux
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
