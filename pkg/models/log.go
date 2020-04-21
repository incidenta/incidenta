package models

import (
	"fmt"

	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Log struct {
	ID         int64 `xorm:"pk autoincr"`
	ReceiverID int64
	Receiver   *Receiver `xorm:"-"`
	AlertID    int64
	Alert      *Alert `xorm:"-"`
	Username   string
	Comment    string

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (l *Log) APIFormat() *apiv1.Log {
	return &apiv1.Log{
		ID:          l.ID,
		ReceiverID:  l.ReceiverID,
		AlertID:     l.AlertID,
		Username:    l.Username,
		Comment:     l.Comment,
		CreatedUnix: l.CreatedUnix.AsTime(),
		UpdatedUnix: l.UpdatedUnix.AsTime(),
	}
}

type ErrLogNotExist struct {
	ID int64
}

func IsErrLogNotExist(err error) bool {
	_, ok := err.(ErrLogNotExist)
	return ok
}

func (e ErrLogNotExist) Error() string {
	return fmt.Sprintf("Log does not exist [id: %d]", e.ID)
}

func DeleteLog(l *Log) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(l.ID).Delete(new(Log)); err != nil {
		return err
	}
	return sess.Commit()
}

func GetLogByID(id int64) (*Log, error) {
	u := new(Log)
	has, err := x.ID(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrLogNotExist{ID: id}
	}
	return u, nil
}

type SearchLogsOptions struct {
	ReceiverID int64
	AlertID    int64
	OrderBy    SearchOrderBy
}

func (o *SearchLogsOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	if o.ReceiverID >= 0 {
		cond = cond.And(builder.Eq{"receiver_id": o.ReceiverID})
	}
	if o.AlertID >= 0 {
		cond = cond.And(builder.Eq{"alert_id": o.AlertID})
	}
	return cond
}

func SearchLogs(opts *SearchLogsOptions) ([]*Log, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Log))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByNewest
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var logs []*Log
	return logs, count, sess.Find(&logs)
}
