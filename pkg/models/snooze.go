package models

import (
	"fmt"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
	"xorm.io/xorm"
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

func GetSnoozeByAlertFingerprint(fingerprint string) (*Snooze, error) {
	return getSnoozeByAlertFingerprint(x, fingerprint)
}

type ErrSnoozeNotExist struct {
	ID          int64
	Fingerprint string
}

func IsErrSnoozeNotExist(err error) bool {
	_, ok := err.(ErrSnoozeNotExist)
	return ok
}

func (e ErrSnoozeNotExist) Error() string {
	return fmt.Sprintf("Snooze does not exist [id: %d, fingerprint: %s]", e.ID, e.Fingerprint)
}

func getSnoozeByAlertFingerprint(x *xorm.Engine, fingerprint string) (*Snooze, error) {
	if len(fingerprint) == 0 {
		return nil, ErrSnoozeNotExist{Fingerprint: fingerprint}
	}
	s := &Snooze{AlertFingerprint: fingerprint}
	has, err := x.Get(s)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrSnoozeNotExist{Fingerprint: fingerprint}
	}
	return s, nil
}
