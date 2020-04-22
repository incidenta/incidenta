package models

import (
	"fmt"
	"time"

	"github.com/prometheus/alertmanager/template"
	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Alert struct {
	ID           int64 `xorm:"pk autoincr"`
	Name         string
	ProjectID    int64
	Project      *Project `xorm:"-"`
	Fingerprint  string   `xorm:"INDEX NOT NULL"`
	Labels       map[string]string
	GeneratorURL string

	SnoozedUnix timeutil.TimeStamp

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (a *Alert) APIFormat() *apiv1.Alert {
	return &apiv1.Alert{
		ID:           a.ID,
		Name:         a.Name,
		ProjectID:    a.ProjectID,
		Fingerprint:  a.Fingerprint,
		Labels:       a.Labels,
		GeneratorURL: a.GeneratorURL,
		Snoozed:      a.IsSnoozed(),
		CreatedAt:    a.CreatedUnix.AsTime(),
		UpdatedAt:    a.UpdatedUnix.AsTime(),
	}
}

func (a *Alert) IsSnoozed() bool {
	return !a.SnoozedUnix.IsZero() && time.Since(a.SnoozedUnix.AsTime()) > 0
}

type ErrAlertNotExist struct {
	ID          int64
	Fingerprint string
}

func IsErrAlertNotExist(err error) bool {
	_, ok := err.(ErrAlertNotExist)
	return ok
}

func (e ErrAlertNotExist) Error() string {
	return fmt.Sprintf("Alert does not exist [id: %d, fingerprint: %s]", e.ID, e.Fingerprint)
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

func GetAlertByFingerprint(fingerprint string) (*Alert, error) {
	if len(fingerprint) == 0 {
		return nil, ErrAlertNotExist{Fingerprint: fingerprint}
	}
	a := &Alert{Fingerprint: fingerprint}
	has, err := x.Get(a)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAlertNotExist{Fingerprint: fingerprint}
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
	ProjectID int64
	OrderBy   SearchOrderBy
}

func (o *SearchAlertsOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	if o.ProjectID > 0 {
		cond = cond.And(builder.Eq{"project_id": o.ProjectID})
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

func isAlertExist(id int64, fingerprint string) (bool, error) {
	if len(fingerprint) == 0 {
		return false, nil
	}
	return x.
		Where("id!=?", id).
		Get(&Alert{Fingerprint: fingerprint})
}

type ErrAlertAlreadyExist struct {
	Fingerprint string
}

func IsErrAlertAlreadyExist(err error) bool {
	_, ok := err.(ErrAlertAlreadyExist)
	return ok
}

func (e ErrAlertAlreadyExist) Error() string {
	return fmt.Sprintf("Alert already exist [fingerprint: %s]", e.Fingerprint)
}

func CreateAlert(a *Alert) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	isExist, err := isAlertExist(0, a.Fingerprint)
	if err != nil {
		return err
	} else if isExist {
		return ErrAlertAlreadyExist{Fingerprint: a.Fingerprint}
	}
	if _, err = sess.Insert(a); err != nil {
		return err
	}
	return sess.Commit()
}

func CreateAlertFromAlertmanagerAlert(project *Project, alert template.Alert) (*Alert, error) {
	name, ok := alert.Labels["alertname"]
	if !ok {
		name = "unnamed"
	}
	a := &Alert{
		Name:         name,
		ProjectID:    project.ID,
		Fingerprint:  alert.Fingerprint,
		Labels:       alert.Labels,
		CreatedUnix:  timeutil.TimeStampNow(),
		GeneratorURL: alert.GeneratorURL,
	}
	if err := CreateAlert(a); err != nil {
		return nil, err
	}
	_ = a.SysLog(fmt.Sprintf("Created with %s status", alert.Status))
	return a, nil
}

func (a *Alert) SysLog(comment string) error {
	return a.Log("IncidentaBot", comment)
}

func (a *Alert) Log(username, comment string) error {
	return CreateLog(&Log{
		ProjectID:   a.ProjectID,
		AlertID:     a.ID,
		Username:    "IncidentaBot",
		Comment:     comment,
		CreatedUnix: timeutil.TimeStampNow(),
	})
}
