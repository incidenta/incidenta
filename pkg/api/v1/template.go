package v1

import (
	"fmt"
	"net/http"
	"time"
)

// swagger:parameters GetTemplate DeleteTemplate
type swaggerTemplateOptions struct {
	// in: path
	TemplateID int `json:"template_id"`
}

// swagger:model
type Template struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Content     string    `json:"content,omitempty"`
	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}

// swagger:parameters CreateTemplate
type swaggerTemplateCreateOptions struct {
	// in: body
	Body *TemplateCreateOptions
}

type TemplateCreateOptions struct {
	// required: true
	Name string `json:"name" validate:"required"`
	// required: true
	Content string `json:"content" validate:"required"`
}

// swagger:parameters EditTemplate
type swaggerTemplateEditOptions struct {
	swaggerTemplateOptions
	// in: body
	Body *TemplateEditOptions
}

type TemplateEditOptions struct {
	Name    *string `json:"name"`
	Content *string `json:"content"`
}

type Templates struct {
	c *Client
}

func (t *Templates) Get(id int64) (*Template, *Response, error) {
	req, err := t.c.newRequest(http.MethodGet, fmt.Sprintf("v1/template/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	template := &Template{}
	resp, err := t.c.doRequest(req, template)
	return template, resp, err
}

func (t *Templates) Delete(id int64) (*Response, error) {
	req, err := t.c.newRequest(http.MethodDelete, fmt.Sprintf("v1/template/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return t.c.doRequest(req, nil)
}

func (t *Templates) Create(opts *TemplateCreateOptions) (*Template, *Response, error) {
	req, err := t.c.newRequest(http.MethodPost, "v1/template", opts)
	if err != nil {
		return nil, nil, err
	}
	template := &Template{}
	resp, err := t.c.doRequest(req, template)
	return template, resp, err
}

func (t *Templates) Edit(id int64, opts *TemplateEditOptions) (*Template, *Response, error) {
	req, err := t.c.newRequest(http.MethodPost, fmt.Sprintf("v1/template/%d", id), opts)
	if err != nil {
		return nil, nil, err
	}
	template := &Template{}
	resp, err := t.c.doRequest(req, template)
	return template, resp, err
}

func (t *Templates) List() ([]*Template, *Response, error) {
	req, err := t.c.newRequest(http.MethodGet, "v1/templates", nil)
	if err != nil {
		return nil, nil, err
	}
	var templates []*Template
	resp, err := t.c.doRequest(req, &templates)
	return templates, resp, err
}
