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

// swagger:route GET /v1/project Project ListProject
//
// List operation
//
// 	Responses:
// 		200: []Project
// 		500: GenericError
func (h *HTTPServer) ProjectListRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &models.SearchProjectsOptions{}
	projects, _, err := models.SearchProjects(opts)
	if err != nil {
		return InternalError(err)
	}
	if projects == nil {
		return JSON(200, []*apiv1.Project{})
	}
	var apiProjects []*apiv1.Project
	for _, project := range projects {
		apiProjects = append(apiProjects, project.APIFormat())
	}

	return JSON(200, apiProjects)
}

// swagger:route POST /v1/project Project CreateProject
//
// Create operation
//
// 	Responses:
// 		201: Project
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) ProjectCreateRequest(_ http.ResponseWriter, r *http.Request) Response {
	opts := &apiv1.ProjectCreateOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}

	if err := validator.Validate(opts); err != nil {
		return ValidationError(err)
	}

	p := &models.Project{
		Name:          opts.Name,
		Description:   opts.Description,
		SlackURL:      opts.SlackURL,
		SlackChannel:  opts.SlackChannel,
		AckButton:     opts.AckButton,
		ResolveButton: opts.ResolveButton,
		SnoozeButton:  opts.SnoozeButton,
		CreatedUnix:   timeutil.TimeStampNow(),
	}

	err := models.CreateProject(p)
	if err != nil {
		if models.IsErrProjectAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return InternalError(err)
	}

	return JSON(201, p.APIFormat())
}

// swagger:route GET /v1/project/{project_id} Project GetProject
//
// Get operation
//
// 	Responses:
// 		200: Project
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ProjectGetRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	p, err := models.GetProjectByID(id)
	if err != nil {
		if models.IsErrProjectNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	return JSON(200, p.APIFormat())
}

// swagger:route POST /v1/project/{project_id} Project EditProject
//
// Edit operation
//
// 	Responses:
// 		200: Project
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ProjectEditRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	opts := &apiv1.ProjectEditOptions{}
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		return Error(400, "Failed to decode request", err)
	}
	if err := validator.Validate(opts); err != nil {
		return ValidationError(err)
	}

	ar, err := models.GetProjectByID(id)
	if err != nil {
		if models.IsErrProjectNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
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

	if opts.SlackChannel != nil {
		ar.SlackChannel = *opts.SlackChannel
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

	err = models.EditProject(ar)
	if err != nil {
		if models.IsErrProjectAlreadyExist(err) {
			return Error(400, "Already exists", err)
		}
		return InternalError(err)
	}

	return JSON(200, ar)
}

// swagger:route DELETE /v1/project/{project_id} Project DeleteProject
//
// Delete operation
//
// 	Responses:
// 		204:
// 		400: GenericError
// 		404: GenericError
// 		500: GenericError
func (h *HTTPServer) ProjectDeleteRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	p, err := models.GetProjectByID(id)
	if err != nil {
		if models.IsErrProjectNotExist(err) {
			return Error(404, "Not Found", nil)
		}
		return InternalError(err)
	}

	err = models.DeleteProject(p)
	if err != nil {
		return InternalError(err)
	}

	return Empty(204)
}

// swagger:route GET /v1/project/{project_id}/alerts Project ListProjectAlerts
//
// List alerts operation
//
// 	Responses:
// 		200: []Alert
// 		400: GenericError
// 		500: GenericError
func (h *HTTPServer) ProjectAlertsRequest(_ http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		return ValidationError(err)
	}

	opts := &models.SearchAlertsOptions{
		ProjectID: id,
	}
	alerts, _, err := models.SearchAlerts(opts)
	if err != nil {
		return InternalError(err)
	}

	if alerts == nil {
		return JSON(200, []*apiv1.Alert{})
	}

	var apiAlerts []*apiv1.Alert
	for _, alert := range alerts {
		apiAlerts = append(apiAlerts, alert.APIFormat())
	}

	return JSON(200, apiAlerts)
}
