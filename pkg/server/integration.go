package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/sirupsen/logrus"

	"github.com/incidenta/incidenta/pkg/alertmanager"
	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/models"
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
// 		201: AlertmanagerResponse
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertmanagerRequest(_ http.ResponseWriter, r *http.Request) Response {
	message, err := alertmanager.ParseMessage(r)
	if err != nil {
		return Error(400, "Bad Request", err)
	}

	vars := mux.Vars(r)

	project, err := models.GetProjectByUID(vars["project_uid"])
	if err != nil {
		if models.IsErrProjectNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	var errs []string
	for _, alert := range message.Alerts() {
		alertInt, err := models.GetAlertByFingerprint(alert.Fingerprint)
		if err != nil {
			if models.IsErrAlertNotExist(err) {
				a, err := models.CreateAlertFromAlertmanagerAlert(project, alert)
				if err != nil {
					errs = append(errs, err.Error())
					logrus.WithError(err).Error("Failed to create alert")
					continue
				}
				alertInt = a
			} else {
				errs = append(errs, err.Error())
				continue
			}
		} else {
			_ = alertInt.SysEvent(alert.Status, "")
		}
		if alertInt.IsSnoozed() {
			continue
		}
	}
	return JSON(201, &apiv1.AlertmanagerResponse{
		Errors: errs,
	})
}
