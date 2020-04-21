package server

import (
	"net/http"

	"github.com/incidenta/incidenta/pkg/alertmanager"
)

// swagger:route POST /v1/alertmanager/webhook Alertmanager AlertmanagerWebhook
//
// Alertmanager Webhook operation
//
// 	Responses:
// 		201:
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertmanagerWebhookRequest(_ http.ResponseWriter, r *http.Request) Response {
	message, err := alertmanager.ParseMessage(r)
	if err != nil {
		return Error(400, "Bad Request", err)
	}
	return JSON(201, map[string]int{
		"count": len(message.Alerts()),
	})
}
