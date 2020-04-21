package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/models"
	"github.com/incidenta/incidenta/pkg/validator"
)

func (h *HTTPServer) TemplateListRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &models.SearchTemplatesOptions{}
	templates, _, err := models.SearchTemplates(opts)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}
	var apiTemplates []*apiv1.Template
	for _, template := range templates {
		apiTemplates = append(apiTemplates, template.APIFormat(false))
	}
	return JSON(200, apiTemplates)
}

func (h *HTTPServer) TemplateGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["template_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	t, err := models.GetTemplateByID(id)
	if err != nil {
		if models.IsErrTemplateNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(200, t.APIFormat(true))
}

func (h *HTTPServer) TemplateDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["template_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	t, err := models.GetTemplateByID(id)
	if err != nil {
		if models.IsErrTemplateNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	err = models.DeleteTemplate(t)
	if err != nil {
		return Error(500, "Internal Server Error", err)
	}

	return Empty(204)
}

func (h *HTTPServer) TemplateEditRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["template_id"], 10, 64)
	if err != nil {
		return Error(400, "Validation error", err)
	}

	opts := &apiv1.TemplateEditOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}
	if err := validator.Validate(opts); err != nil {
		return Error(400, "Validation failed", err)
	}

	t, err := models.GetTemplateByID(id)
	if err != nil {
		if models.IsErrReceiverNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return Error(500, "Internal Server Error", err)
	}

	if opts.Content != nil {
		t.Content = *opts.Content
	}

	if opts.Name != nil {
		t.Name = *opts.Name
	}

	err = models.EditTemplate(t)
	if err != nil {
		if models.IsErrTemplateAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(200, t)
}

func (h *HTTPServer) TemplateCreateRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &apiv1.TemplateCreateOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}

	if err := validator.Validate(opts); err != nil {
		return Error(400, "Validation failed", err)
	}

	t := &models.Template{
		Name:    opts.Name,
		Content: opts.Content,
	}

	err := models.CreateTemplate(t)
	if err != nil {
		if models.IsErrTemplateAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return Error(500, "Internal Server Error", err)
	}

	return JSON(201, t.APIFormat(true))
}
