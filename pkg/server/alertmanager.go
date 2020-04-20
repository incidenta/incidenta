package server

import (
	"net/http"

	"github.com/incidenta/incidenta/pkg/alertmanager"
	"github.com/incidenta/incidenta/pkg/models"
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
	_, err = models.GetReceiverByName(message.Receiver)
	if err != nil {
		if models.IsErrReceiverNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}
	return Empty(201)
}
