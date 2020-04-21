package server

import (
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"

	"github.com/incidenta/incidenta/pkg/alertmanager"
)

// swagger:parameters Alertmanager
type swaggerMessage struct {
	// in: path
	ProjectUID string `json:"project_uid"`
	// in: body
	Body *webhook.Message
}

// swagger:route POST /v1/integrations/alertmanager/{project_uid} Integration Alertmanager
//
// Alertmanager integration
//
// 	Responses:
// 		201:
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertmanagerRequest(_ http.ResponseWriter, r *http.Request) Response {
	message, err := alertmanager.ParseMessage(r)
	if err != nil {
		return Error(400, "Bad Request", err)
	}
	return JSON(201, map[string]int{
		"count": len(message.Alerts()),
	})
}
