package api

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
	ID          int64     `json:"id"`
	ReceiverID  int64     `json:"receiver_id"`
	AlertID     int64     `json:"alert_id"`
	Username    string    `json:"username"`
	Comment     string    `json:"comment"`
	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}
