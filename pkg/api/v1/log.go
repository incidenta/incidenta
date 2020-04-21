package v1

import (
	"time"
)

// swagger:parameters GetLog DeleteLog
type swaggerLogOptions struct {
	// in: path
	LogID int `json:"log_id"`
}

// swagger:model
type Log struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	AlertID   int64     `json:"alert_id"`
	Username  string    `json:"username"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Logs struct {
	c *Client
}
