package v1

import (
	"fmt"
	"net/http"

	"github.com/prometheus/alertmanager/notify/webhook"
)

type Integrations struct {
	c *Client
}

// swagger:model
type AlertmanagerResponse struct {
	Errors []string `json:"errors"`
}

func (i *Integrations) AlertmanagerEvent(projectUID string, payload *webhook.Message) (*AlertmanagerResponse, *Response, error) {
	req, err := i.c.newRequest(http.MethodPost, fmt.Sprintf("v1/integrations/alertmanager/%s", projectUID), payload)
	if err != nil {
		return nil, nil, err
	}
	managerResp := &AlertmanagerResponse{}
	resp, err := i.c.doRequest(req, managerResp)
	return managerResp, resp, err
}
