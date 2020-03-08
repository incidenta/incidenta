package alertmanager

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"
)

type Message webhook.Message

func ParseMessage(r *http.Request) (*Message, error) {
	m := &Message{}
	return m, json.NewDecoder(r.Body).Decode(m)
}
