package models

import (
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Snooze struct {
	ID       int64 `xorm:"pk autoincr"`
	AlertID  int64
	Alert    *Alert `xorm:"-"`
	Username string

	DeadlineUnix timeutil.TimeStamp `xorm:"INDEX deadline"`
	CreatedUnix  timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
}
