package v1

import (
	"time"
)

type Snooze struct {
	ID           int64     `json:"id"`
	AlertID      int64     `json:"alert_id" validate:"required"`
	Username     string    `json:"username" validate:"required"`
	DeadlineUnix time.Time `json:"deadline_at" validate:"required"`

	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}

type Snoozes struct {
	c *Client
}
