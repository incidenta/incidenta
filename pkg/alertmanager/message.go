package alertmanager

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
)

// swagger:parameters AlertmanagerWebhook
type swaggerMessage struct {
	// in: body
	Body *webhook.Message
}

type Message webhook.Message

func ParseMessage(r *http.Request) (*Message, error) {
	m := &Message{}
	return m, json.NewDecoder(r.Body).Decode(m)
}

func (m *Message) Alerts() []template.Alert {
	return m.Data.Alerts
}
