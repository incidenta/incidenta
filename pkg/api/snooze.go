package api

import (
	"time"
)

type Snooze struct {
	ID           int64     `json:"id"`
	AlertID      int64     `json:"alert_id"`
	Alert        *Alert    `json:"alert"`
	Username     string    `json:"username"`
	DeadlineUnix time.Time `json:"deadline_at"`

	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}
