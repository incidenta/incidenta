package v1

import (
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
