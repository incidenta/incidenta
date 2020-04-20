package api

import (
	"time"
)

// swagger:parameters GetReceiver DeleteReceiver ListReceiverAlerts
type swaggerReceiverOptions struct {
	// in: path
	ReceiverID int `json:"receiver_id"`
}

// swagger:model
type Receiver struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	SlackURL      string    `json:"slack_url"`
	TemplateID    int64     `json:"template_id"`
	AckButton     bool      `json:"ack_button"`
	ResolveButton bool      `json:"resolve_button"`
	SnoozeButton  bool      `json:"snooze_button"`
	CreatedUnix   time.Time `json:"created_at"`
	UpdatedUnix   time.Time `json:"updated_at"`
}

// swagger:parameters CreateReceiver
type swaggerReceiverCreateOptions struct {
	// in: body
	Body *ReceiverCreateOptions
}

type ReceiverCreateOptions struct {
	// required: true
	Name        string `json:"name" required:"true"`
	Description string `json:"description"`
	// required: true
	SlackURL string `json:"slack_url" required:"true"`
	// required: true
	TemplateID    int64 `json:"template_id" required:"true"`
	AckButton     bool  `json:"ack_button"`
	ResolveButton bool  `json:"resolve_button"`
	SnoozeButton  bool  `json:"snooze_button"`
}

// swagger:parameters EditReceiver
type swaggerReceiverEditOptions struct {
	swaggerReceiverOptions
	// in: body
	Body *ReceiverEditOptions
}

type ReceiverEditOptions struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	SlackURL      *string `json:"slack_url"`
	TemplateID    *int64  `json:"template_id"`
	AckButton     *bool   `json:"ack_button"`
	ResolveButton *bool   `json:"resolve_button"`
	SnoozeButton  *bool   `json:"snooze_button"`
}
