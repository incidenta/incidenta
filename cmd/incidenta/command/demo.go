package command

import (
	"strings"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/generate"
)

func newDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Create demo content",
		RunE:  demoRun,
	}

	flags := cmd.PersistentFlags()
	flags.String("address", "http://127.0.0.1:8080", "Address of API server")

	return cmd
}

func demoRun(cmd *cobra.Command, _ []string) error {
	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		return err
	}

	projectOpts := []*apiv1.ProjectCreateOptions{
		{
			Name:          "Service Production",
			Description:   "Alerts from production environment",
			SlackURL:      "http://slack.com",
			SlackChannel:  "#service-prod",
			AckButton:     true,
			ResolveButton: true,
			SnoozeButton:  true,
		},
		{
			Name:         "Service Staging",
			Description:  "Alerts from staging environment",
			SlackURL:     "http://slack.com",
			SlackChannel: "#service-stage",
		},
		{
			Name:         "Service Testing",
			Description:  "Alerts from testing environment",
			SlackURL:     "http://slack.com",
			SlackChannel: "#service-test",
		},
	}

	cli := apiv1.NewClient(nil, addr)

	alertNames := []string{
		"ServiceDown",
		"NodeDown",
		"FreeDiskSpace",
	}

	for _, opts := range projectOpts {
		project, _, err := cli.Projects.Create(opts)
		if err != nil {
			return err
		}
		logrus.Infof("Project %q [uid=%s] created", project.Name, project.UID)
		env := strings.TrimPrefix(strings.ToLower(project.Name), "service ")
		for _, alertName := range alertNames {
			fingerprint, err := generate.GetRandomString(16)
			if err != nil {
				return err
			}
			alert := template.Alert{
				Status: "firing",
				Labels: map[string]string{
					"alertname": alertName,
					"instance":  "service.local",
					"env":       env,
				},
				Annotations: map[string]string{
					"summary": alertName,
				},
				StartsAt:     time.Now(),
				EndsAt:       time.Now().Add(10 * time.Minute),
				GeneratorURL: "http://prometheus.local",
				Fingerprint:  fingerprint,
			}
			_, _, err = cli.Integrations.AlertmanagerEvent(project.UID, &webhook.Message{
				Data: &template.Data{
					Alerts: []template.Alert{alert},
				},
			})
			if err != nil {
				return err
			}
			logrus.Info("Alertmanager event triggered")
			alert.Status = "resolved"
			_, _, err = cli.Integrations.AlertmanagerEvent(project.UID, &webhook.Message{
				Data: &template.Data{
					Alerts: []template.Alert{alert},
				},
			})
			if err != nil {
				return err
			}
			logrus.Info("Alertmanager event triggered")
		}
	}

	return nil
}
