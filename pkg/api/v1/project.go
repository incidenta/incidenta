package v1

import (
	"fmt"
	"net/http"
	"time"
)

// swagger:parameters GetProject DeleteProject ListProjectAlerts
type swaggerProjectOptions struct {
	// in: path
	ProjectID int `json:"project_id"`
}

// swagger:model
type Project struct {
	ID            int64     `json:"id"`
	UID           string    `json:"uid"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	SlackURL      string    `json:"slack_url"`
	SlackChannel  string    `json:"slack_channel"`
	AckButton     bool      `json:"ack_button"`
	ResolveButton bool      `json:"resolve_button"`
	SnoozeButton  bool      `json:"snooze_button"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// swagger:parameters CreateProject
type swaggerProjectCreateOptions struct {
	// in: body
	Body *ProjectCreateOptions
}

type ProjectCreateOptions struct {
	// required: true
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// required: true
	SlackURL string `json:"slack_url" validate:"required"`
	// required: true
	SlackChannel  string `json:"slack_channel" validate:"required"`
	AckButton     bool   `json:"ack_button"`
	ResolveButton bool   `json:"resolve_button"`
	SnoozeButton  bool   `json:"snooze_button"`
}

// swagger:parameters EditProject
type swaggerProjectEditOptions struct {
	swaggerProjectOptions
	// in: body
	Body *ProjectEditOptions
}

type ProjectEditOptions struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	SlackURL      *string `json:"slack_url"`
	SlackChannel  *string `json:"slack_channel"`
	AckButton     *bool   `json:"ack_button"`
	ResolveButton *bool   `json:"resolve_button"`
	SnoozeButton  *bool   `json:"snooze_button"`
}

type Projects struct {
	c *Client
}

func (p *Projects) Delete(id int64) (*Response, error) {
	req, err := p.c.newRequest(http.MethodDelete, fmt.Sprintf("v1/project/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return p.c.doRequest(req, nil)
}

func (p *Projects) Get(id int64) (*Project, *Response, error) {
	req, err := p.c.newRequest(http.MethodGet, fmt.Sprintf("v1/project/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	project := &Project{}
	resp, err := p.c.doRequest(req, project)
	return project, resp, err
}

func (p *Projects) Create(opts *ProjectCreateOptions) (*Project, *Response, error) {
	req, err := p.c.newRequest(http.MethodPost, "v1/project", opts)
	if err != nil {
		return nil, nil, err
	}
	project := &Project{}
	resp, err := p.c.doRequest(req, project)
	return project, resp, err
}

func (p *Projects) Edit(id int64, opts *ProjectEditOptions) (*Project, *Response, error) {
	req, err := p.c.newRequest(http.MethodPost, fmt.Sprintf("v1/project/%d", id), opts)
	if err != nil {
		return nil, nil, err
	}
	project := &Project{}
	resp, err := p.c.doRequest(req, project)
	return project, resp, err
}

func (p *Projects) List() ([]*Project, *Response, error) {
	req, err := p.c.newRequest(http.MethodGet, "v1/projects", nil)
	if err != nil {
		return nil, nil, err
	}
	var projects []*Project
	resp, err := p.c.doRequest(req, &projects)
	return projects, resp, err
}
