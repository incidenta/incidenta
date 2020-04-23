package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/incidenta/incidenta/pkg/config"
	"github.com/incidenta/incidenta/pkg/models"
	"github.com/incidenta/incidenta/pkg/server"
)

func newWebCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start Incidenta web server",
		RunE:  webRun,
	}

	flags := cmd.PersistentFlags()
	flags.String("env-file", ".env", "Read in a file of environment variables")

	return cmd
}

func webRun(cmd *cobra.Command, _ []string) error {
	envFile, err := cmd.Flags().GetString("env-file")
	if err != nil {
		return err
	}

	_ = godotenv.Load(envFile)

	cfg, err := config.Environ()
	if err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	initLogging(cfg)

	logrus.Infof("Starting incidenta on %s...", cfg.Server.Addr)

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.Infoln("Configuration:", cfg.String())
	}

	if err := models.NewEngine(context.Background(), cfg.Database); err != nil {
		return err
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

	return models.Close()
}
