package command

import (
	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	for _, opts := range projectOpts {
		project, _, err := cli.Projects.Create(opts)
		if err != nil {
			return err
		}
		logrus.Infof("Project %q [uid=%s] created", project.Name, project.UID)
	}

	return nil
}
