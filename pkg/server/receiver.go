package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/models"
	"github.com/incidenta/incidenta/pkg/timeutil"
	"github.com/incidenta/incidenta/pkg/validator"
)

// swagger:route GET /v1/receiver Receiver ListReceiver
//
// List operation
//
// 	Responses:
// 		200: []Receiver
// 		500: GenericError
func (h *HTTPServer) ReceiverListRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &models.SearchReceiversOptions{}
	receivers, _, err := models.SearchReceivers(opts)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}
	var apiReceivers []*apiv1.Receiver
	for _, receiver := range receivers {
		apiReceivers = append(apiReceivers, receiver.APIFormat())
	}
	return JSON(200, apiReceivers)
}

// swagger:route POST /v1/receiver Receiver CreateReceiver
//
// Create operation
//
// 	Responses:
// 		201: Receiver
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) ReceiverCreateRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &apiv1.ReceiverCreateOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}

	if err := validator.Validate(opts); err != nil {
		return Error(400, "Validation failed", err)
	}

	rec := &models.Receiver{
		Name:          opts.Name,
		Description:   opts.Description,
		SlackURL:      opts.SlackURL,
		TemplateID:    opts.TemplateID,
		AckButton:     opts.AckButton,
		ResolveButton: opts.ResolveButton,
		SnoozeButton:  opts.SnoozeButton,
		CreatedUnix:   timeutil.TimeStampNow(),
	}

	err := models.CreateReceiver(rec)
	if err != nil {
		if models.IsErrReceiverAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(201, rec.APIFormat())
}

// swagger:route GET /v1/receiver/{receiver_id} Receiver GetReceiver
//
// Get operation
//
// 	Responses:
// 		200: Receiver
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ReceiverGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["receiver_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	ar, err := models.GetReceiverByID(id)
	if err != nil {
		if models.IsErrReceiverNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(200, ar.APIFormat())
}

// swagger:route POST /v1/receiver/{receiver_id} Receiver EditReceiver
//
// Edit operation
//
// 	Responses:
// 		200: Receiver
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ReceiverEditRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["receiver_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	opts := &apiv1.ReceiverEditOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}
	if err := validator.Validate(opts); err != nil {
		return Error(400, "Validation failed", err)
	}

	ar, err := models.GetReceiverByID(id)
	if err != nil {
		if models.IsErrReceiverNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	if opts.Name != nil {
		ar.Name = *opts.Name
	}

	if opts.Description != nil {
		ar.Description = *opts.Description
	}

	if opts.SlackURL != nil {
		ar.SlackURL = *opts.SlackURL
	}

	if opts.TemplateID != nil {
		ar.TemplateID = *opts.TemplateID
	}

	if opts.AckButton != nil {
		ar.AckButton = *opts.AckButton
	}

	if opts.ResolveButton != nil {
		ar.ResolveButton = *opts.ResolveButton
	}

	if opts.SnoozeButton != nil {
		ar.SnoozeButton = *opts.SnoozeButton
	}

	err = models.EditReceiver(ar)
	if err != nil {
		if models.IsErrReceiverAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(200, ar)
}

// swagger:route DELETE /v1/receiver/{receiver_id} Receiver DeleteReceiver
//
// Delete operation
//
// 	Responses:
// 		204:
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ReceiverDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["receiver_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	ar, err := models.GetReceiverByID(id)
	if err != nil {
		if models.IsErrReceiverNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	err = models.DeleteReceiver(ar)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}

	return Empty(204)
}

// swagger:route GET /v1/receiver/{{ receiver_id }}/alerts Receiver ListReceiverAlerts
//
// List alerts operation
//
// 	Responses:
// 		200: []Alert
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) ReceiverAlertsRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["receiver_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	opts := &models.SearchAlertsOptions{
		ReceiverID: id,
	}
	alerts, _, err := models.SearchAlerts(opts)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}

	var apiAlerts []*apiv1.Alert
	for _, alert := range alerts {
		apiAlerts = append(apiAlerts, alert.APIFormat())
	}

	return JSON(200, apiAlerts)
}
