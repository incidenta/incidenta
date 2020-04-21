package v1

import (
	"fmt"
	"net/http"
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

func (l *Logs) Get(id int64) (*Log, *Response, error) {
	req, err := l.c.newRequest(http.MethodGet, fmt.Sprintf("v1/log/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	log := &Log{}
	resp, err := l.c.doRequest(req, log)
	return log, resp, err
}

func (l *Logs) Delete(id int64) (*Response, error) {
	req, err := l.c.newRequest(http.MethodDelete, fmt.Sprintf("v1/log/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return l.c.doRequest(req, nil)
}
