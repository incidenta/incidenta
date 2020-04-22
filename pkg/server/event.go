package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/incidenta/incidenta/pkg/models"
)

// swagger:route GET /v1/event/{event_id} Event GetEvent
//
// Get operation
//
// 	Responses:
// 		200: Event
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) EventGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["event_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	al, err := models.GetEventByID(id)
	if err != nil {
		if models.IsErrEventNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	return JSON(200, al.APIFormat())
}

// swagger:route DELETE /v1/event/{event_id} Event DeleteEvent
//
// Delete operation
//
// 	Responses:
// 		204:
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) EventDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["event_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	al, err := models.GetEventByID(id)
	if err != nil {
		if models.IsErrEventNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	err = models.DeleteEvent(al)
	if err != nil {
		return InternalError(err)
	}

	return Empty(204)
}
