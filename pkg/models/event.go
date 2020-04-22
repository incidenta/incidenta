package models

import (
	"fmt"

	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Event struct {
	ID          int64 `xorm:"pk autoincr"`
	ProjectID   int64
	Project     *Project `xorm:"-"`
	AlertID     int64
	Alert       *Alert `xorm:"-"`
	AlertStatus string
	Username    string
	Comment     string

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (e *Event) APIFormat() *apiv1.Event {
	return &apiv1.Event{
		ID:          e.ID,
		ProjectID:   e.ProjectID,
		AlertID:     e.AlertID,
		AlertStatus: e.AlertStatus,
		Username:    e.Username,
		Comment:     e.Comment,
		CreatedAt:   e.CreatedUnix.AsTime(),
		UpdatedAt:   e.UpdatedUnix.AsTime(),
	}
}

type ErrEventNotExist struct {
	ID int64
}

func IsErrEventNotExist(err error) bool {
	_, ok := err.(ErrEventNotExist)
	return ok
}

func (e ErrEventNotExist) Error() string {
	return fmt.Sprintf("Event does not exist [id: %d]", e.ID)
}

func DeleteEvent(e *Event) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(e.ID).Delete(new(Event)); err != nil {
		return err
	}
	return sess.Commit()
}

func GetEventByID(id int64) (*Event, error) {
	u := new(Event)
	has, err := x.ID(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrEventNotExist{ID: id}
	}
	return u, nil
}

type SearchEventsOptions struct {
	ProjectID int64
	AlertID   int64
	OrderBy   SearchOrderBy
}

func (o *SearchEventsOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	if o.ProjectID > 0 {
		cond = cond.And(builder.Eq{"project_id": o.ProjectID})
	}
	if o.AlertID > 0 {
		cond = cond.And(builder.Eq{"alert_id": o.AlertID})
	}
	return cond
}

func SearchEvents(opts *SearchEventsOptions) ([]*Event, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Event))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByNewest
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var events []*Event
	return events, count, sess.Find(&events)
}

func CreateEvent(e *Event) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Insert(e); err != nil {
		return err
	}
	return sess.Commit()
}
