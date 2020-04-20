package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/incidenta/incidenta/pkg/models"
)

// swagger:route GET /v1/log/{log_id} Log GetLog
//
// Get operation
//
// 	Responses:
// 		200: Alert
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) LogGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	alertLogIDRaw := vars["log_id"]
	alertLogID, err := strconv.ParseInt(alertLogIDRaw, 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	al, err := models.GetLogByID(alertLogID)
	if err != nil {
		if models.IsErrLogNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(200, al.APIFormat())
}

// swagger:route DELETE /v1/log/{log_id} Log DeleteLog
//
// Delete operation
//
// 	Responses:
// 		204:
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) LogDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	alertLogIDRaw := vars["log_id"]
	alertLogID, err := strconv.ParseInt(alertLogIDRaw, 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	al, err := models.GetLogByID(alertLogID)
	if err != nil {
		if models.IsErrLogNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	err = models.DeleteLog(al)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}

	return Empty(204)
}
