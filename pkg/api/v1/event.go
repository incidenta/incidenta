package v1

import (
	"fmt"
	"net/http"
	"time"
)

// swagger:parameters GetEvent DeleteEvent
type swaggerEventOptions struct {
	// in: path
	EventID int `json:"event_id"`
}

// swagger:model
type Event struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	AlertID     int64     `json:"alert_id"`
	AlertStatus string    `json:"alert_status"`
	Username    string    `json:"username"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Events struct {
	c *Client
}

func (e *Events) Get(id int64) (*Event, *Response, error) {
	req, err := e.c.newRequest(http.MethodGet, fmt.Sprintf("v1/event/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	event := &Event{}
	resp, err := e.c.doRequest(req, event)
	return event, resp, err
}

func (e *Events) Delete(id int64) (*Response, error) {
	req, err := e.c.newRequest(http.MethodDelete, fmt.Sprintf("v1/event/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return e.c.doRequest(req, nil)
}
