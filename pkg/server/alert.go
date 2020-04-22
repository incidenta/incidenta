package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/models"
)

// swagger:route GET /v1/alert Alert ListAlert
//
// List operation
//
// 	Responses:
// 		200: []Alert
// 		500: GenericError
func (h *HTTPServer) AlertListRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &models.SearchAlertsOptions{}
	alerts, _, err := models.SearchAlerts(opts)
	if err != nil {
		return InternalError(err)
	}

	var apiAlerts []*apiv1.Alert
	for _, alert := range alerts {
		apiAlerts = append(apiAlerts, alert.APIFormat())
	}

	return JSON(200, apiAlerts)
}

// swagger:route GET /v1/alert/{alert_id} Alert GetAlert
//
// Get operation
//
// 	Responses:
// 		200: Alert
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["alert_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	ar, err := models.GetAlertByID(id)
	if err != nil {
		if models.IsErrAlertNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	return JSON(200, ar.APIFormat())
}

// swagger:route DELETE /v1/alert/{alert_id} Alert DeleteAlert
//
// Delete operation
//
// 	Responses:
// 		204:
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["alert_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	a, err := models.GetAlertByID(id)
	if err != nil {
		if models.IsErrAlertNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	err = models.DeleteAlert(a)
	if err != nil {
		return InternalError(err)
	}

	return Empty(204)
}

// swagger:route GET /v1/alert/{alert_id}/events Alert ListAlertEvents
//
// Get events operation
//
// 	Responses:
// 		200: []Event
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) AlertEventsRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["alert_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	events, _, err := models.SearchEvents(&models.SearchEventsOptions{
		AlertID: id,
	})
	if err != nil {
		return InternalError(err)
	}

	var apiEvents []*apiv1.Event
	for _, event := range events {
		apiEvents = append(apiEvents, event.APIFormat())
	}

	return JSON(200, apiEvents)
}
