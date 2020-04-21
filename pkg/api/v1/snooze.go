package v1

import (
	"time"
)

type Snooze struct {
	ID           int64     `json:"id"`
	AlertID      int64     `json:"alert_id" required:"true"`
	Alert        *Alert    `json:"alert"`
	Username     string    `json:"username" required:"true"`
	DeadlineUnix time.Time `json:"deadline_at" required:"true"`

	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}

type Snoozes struct {
	c *Client
}
