package v1

import (
	"fmt"
	"net/http"
	"time"
)

// swagger:parameters GetAlert DeleteAlert ListAlertEvents
type swaggerAlertOptions struct {
	// in: path
	AlertID int `json:"alert_id"`
}

// swagger:model
type Alert struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	ProjectID    int64             `json:"project_id"`
	Fingerprint  string            `json:"fingerprint"`
	Labels       map[string]string `json:"labels"`
	Snoozed      bool              `json:"snoozed"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	GeneratorURL string            `json:"generator_url"`
}

type Alerts struct {
	c *Client
}

func (a *Alerts) List() ([]*Alert, *Response, error) {
	req, err := a.c.newRequest(http.MethodGet, "v1/alerts", nil)
	if err != nil {
		return nil, nil, err
	}
	var alerts []*Alert
	resp, err := a.c.doRequest(req, &alerts)
	return alerts, resp, err
}

func (a *Alerts) Get(id int64) (*Alert, *Response, error) {
	req, err := a.c.newRequest(http.MethodGet, fmt.Sprintf("v1/alert/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	alert := &Alert{}
	resp, err := a.c.doRequest(req, alert)
	return alert, resp, err
}

func (a *Alerts) Delete(id int64) (*Response, error) {
	req, err := a.c.newRequest(http.MethodDelete, fmt.Sprintf("v1/alert/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return a.c.doRequest(req, nil)
}

func (a *Alerts) Events(id int64) ([]*Event, *Response, error) {
	req, err := a.c.newRequest(http.MethodGet, fmt.Sprintf("v1/alert/%d/events", id), nil)
	if err != nil {
		return nil, nil, err
	}
	var events []*Event
	resp, err := a.c.doRequest(req, &events)
	return events, resp, err
}
