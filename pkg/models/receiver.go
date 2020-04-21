package models

import (
	"fmt"
	"strings"

	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Receiver struct {
	ID          int64  `xorm:"pk autoincr"`
	Name        string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Description string
	SlackURL    string

	TemplateID int64
	Template   *Template `xorm:"-"`

	AckButton     bool
	ResolveButton bool
	SnoozeButton  bool

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

type ErrReceiverNotExist struct {
	ID   int64
	Name string
}

func IsErrReceiverNotExist(err error) bool {
	_, ok := err.(ErrReceiverNotExist)
	return ok
}

func (e ErrReceiverNotExist) Error() string {
	return fmt.Sprintf("Receiver does not exist [id: %d, name: %s]", e.ID, e.Name)
}

func GetReceiverByName(name string) (*Receiver, error) {
	r := &Receiver{Name: strings.ToLower(name)}
	has, err := x.Get(r)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrReceiverNotExist{Name: name}
	}
	return r, nil
}

func GetReceiverByID(id int64) (*Receiver, error) {
	r := new(Receiver)
	has, err := x.ID(id).Get(r)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrReceiverNotExist{ID: id}
	}
	return r, nil
}

func DeleteReceiver(r *Receiver) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(r.ID).Delete(new(Receiver)); err != nil {
		return err
	}
	return sess.Commit()
}

func isReceiverExist(uid int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.
		Where("id!=?", uid).
		Get(&Receiver{Name: strings.ToLower(name)})
}

type ErrReceiverAlreadyExist struct {
	Name string
}

func IsErrReceiverAlreadyExist(err error) bool {
	_, ok := err.(ErrReceiverAlreadyExist)
	return ok
}

func (e ErrReceiverAlreadyExist) Error() string {
	return fmt.Sprintf("Receiver already exist [name: %s]", e.Name)
}

func EditReceiver(r *Receiver) error {
	isExist, err := isReceiverExist(r.ID, r.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrReceiverAlreadyExist{Name: r.Name}
	}
	_, err = x.ID(r.ID).AllCols().Update(r)
	return err
}

func CreateReceiver(r *Receiver) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	isExist, err := isReceiverExist(0, r.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrReceiverAlreadyExist{Name: r.Name}
	}
	r.Name = strings.ToLower(r.Name)
	if _, err = sess.Insert(r); err != nil {
		return err
	}
	return sess.Commit()
}

func (r *Receiver) APIFormat() *apiv1.Receiver {
	return &apiv1.Receiver{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description,
		SlackURL:      r.SlackURL,
		TemplateID:    r.TemplateID,
		AckButton:     r.AckButton,
		ResolveButton: r.ResolveButton,
		SnoozeButton:  r.SnoozeButton,
		CreatedUnix:   r.CreatedUnix.AsTime(),
		UpdatedUnix:   r.UpdatedUnix.AsTime(),
	}
}

type SearchReceiversOptions struct {
	TemplateID int64
	OrderBy    SearchOrderBy
}

func (o *SearchReceiversOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	if o.TemplateID >= 0 {
		cond = cond.And(builder.Eq{"template_id": o.TemplateID})
	}
	return cond
}

func SearchReceivers(opts *SearchReceiversOptions) ([]*Receiver, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Receiver))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByAlphabetically
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var receivers []*Receiver
	return receivers, count, sess.Find(&receivers)
}
