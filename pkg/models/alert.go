package models

import (
	"fmt"

	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Alert struct {
	ID          int64 `xorm:"pk autoincr"`
	Name        string
	ReceiverID  int64
	Receiver    *Receiver `xorm:"-"`
	Fingerprint string    `xorm:"INDEX NOT NULL"`
	Body        string

	Snoozed     bool
	SnoozedUnix timeutil.TimeStamp

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (a *Alert) APIFormat() *apiv1.Alert {
	return &apiv1.Alert{
		ID:           a.ID,
		Name:         a.Name,
		ReceiverID:   a.ReceiverID,
		Fingerprint:  a.Fingerprint,
		Body:         a.Body,
		Snoozed:      a.Snoozed,
		SnoozedUntil: a.SnoozedUnix.AsTime(),
		CreatedAt:    a.CreatedUnix.AsTime(),
		UpdatedAt:    a.UpdatedUnix.AsTime(),
	}
}

type ErrAlertNotExist struct {
	ID int64
}

func IsErrAlertNotExist(err error) bool {
	_, ok := err.(ErrAlertNotExist)
	return ok
}

func (e ErrAlertNotExist) Error() string {
	return fmt.Sprintf("Alert does not exist [id: %d]", e.ID)
}

func GetAlertByID(id int64) (*Alert, error) {
	a := new(Alert)
	has, err := x.ID(id).Get(a)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAlertNotExist{ID: id}
	}
	return a, nil
}

func DeleteAlert(a *Alert) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(a.ID).Delete(new(Alert)); err != nil {
		return err
	}
	return sess.Commit()
}

type SearchAlertsOptions struct {
	ReceiverID int64
	OrderBy    SearchOrderBy
}

func (o *SearchAlertsOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	if o.ReceiverID >= 0 {
		cond = cond.And(builder.Eq{"receiver_id": o.ReceiverID})
	}
	return cond
}

func SearchAlerts(opts *SearchAlertsOptions) ([]*Alert, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Alert))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByNewest
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var alerts []*Alert
	return alerts, count, sess.Find(&alerts)
}
