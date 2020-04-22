package server

import (
	"testing"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
)

func TestHTTPServer_AlertmanagerRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	project, _, err := te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "main",
		SlackURL:     "https://slack.com",
		SlackChannel: "#general",
	})
	assert.NoError(t, err)

	stats, _, err := te.client.Integrations.AlertmanagerEvent(project.UID, &webhook.Message{
		Data: &template.Data{
			Receiver: "main",
			Alerts: []template.Alert{
				{
					Status: "firing",
					Labels: map[string]string{
						"aletname": "ServiceDown",
						"instance": "service.local",
						"env":      "prod",
					},
					Annotations: map[string]string{
						"summary": "AlertName @ service.local",
					},
					StartsAt:     time.Now(),
					EndsAt:       time.Now().Add(10 * time.Minute),
					GeneratorURL: "http://prometheus.local",
					Fingerprint:  "AABBCCDD",
				},
				{
					Status: "firing",
					Labels: map[string]string{
						"aletname": "NodeDown",
						"instance": "service.local",
						"env":      "prod",
					},
					Annotations: map[string]string{
						"summary": "NodeDown @ service.local",
					},
					StartsAt:     time.Now(),
					EndsAt:       time.Now().Add(10 * time.Minute),
					GeneratorURL: "http://prometheus.local",
					Fingerprint:  "AABBCCFF",
				},
			},
		},
	})
	assert.Equal(t, len(stats.Errors), 0)
	assert.NoError(t, err)

	alerts, _, err := te.client.Projects.Alerts(project.ID)
	assert.NoError(t, err)
	assert.Equal(t, len(alerts), 2)

	for _, alert := range alerts {
		logs, _, err := te.client.Alerts.Logs(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(logs), 1)
	}
}
