package models

import (
	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Snooze struct {
	ID               int64  `xorm:"pk autoincr"`
	AlertFingerprint string `xorm:"INDEX fingerprint"`
	AlertID          int64
	Alert            *Alert `xorm:"-"`
	Username         string

	DeadlineUnix timeutil.TimeStamp `xorm:"INDEX deadline"`
	CreatedUnix  timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (s *Snooze) APIFormat() *apiv1.Snooze {
	return &apiv1.Snooze{
		ID:           s.ID,
		AlertID:      s.AlertID,
		Username:     s.Username,
		DeadlineUnix: s.DeadlineUnix.AsTime(),
		CreatedUnix:  s.CreatedUnix.AsTime(),
		UpdatedUnix:  s.UpdatedUnix.AsTime(),
	}
}
