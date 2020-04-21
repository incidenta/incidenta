package api

import (
	"time"
)

// swagger:parameters GetAlert DeleteAlert ListAlertLogs
type swaggerAlertOptions struct {
	// in: path
	AlertID int `json:"alert_id"`
}

// swagger:model
type Alert struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ReceiverID  int64     `json:"receiver_id"`
	Fingerprint string    `json:"fingerprint"`
	Body        string    `json:"body"`
	CreatedUnix time.Time `json:"created_at"`
	UpdatedUnix time.Time `json:"updated_at"`
}
