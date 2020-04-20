package api

import (
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
	Name string `json:"name" required:"true"`
	// required: true
	Content string `json:"content" required:"true"`
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
