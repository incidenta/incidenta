package models

import (
	"fmt"
	"strings"

	"xorm.io/builder"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/generate"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Project struct {
	ID           int64  `xorm:"pk autoincr"`
	UID          string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Name         string
	Description  string
	SlackURL     string
	SlackChannel string

	AckButton     bool
	ResolveButton bool
	SnoozeButton  bool

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

type ErrProjectNotExist struct {
	ID   int64
	UID  string
	Name string
}

func IsErrProjectNotExist(err error) bool {
	_, ok := err.(ErrProjectNotExist)
	return ok
}

func (e ErrProjectNotExist) Error() string {
	return fmt.Sprintf("Project does not exist [id: %d, uid: %s, name: %s]", e.ID, e.UID, e.Name)
}

func GetProjectByUID(uid string) (*Project, error) {
	p := &Project{UID: strings.ToLower(uid)}
	has, err := x.Get(p)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrProjectNotExist{UID: uid}
	}
	return p, nil
}

func GetProjectByID(id int64) (*Project, error) {
	p := new(Project)
	has, err := x.ID(id).Get(p)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrProjectNotExist{ID: id}
	}
	return p, nil
}

func DeleteProject(p *Project) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(p.ID).Delete(new(Project)); err != nil {
		return err
	}
	return sess.Commit()
}

func isProjectExist(id int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.
		Where("id!=?", id).
		Get(&Project{Name: name})
}

type ErrProjectAlreadyExist struct {
	Name string
}

func IsErrProjectAlreadyExist(err error) bool {
	_, ok := err.(ErrProjectAlreadyExist)
	return ok
}

func (e ErrProjectAlreadyExist) Error() string {
	return fmt.Sprintf("Project already exist [name: %s]", e.Name)
}

func EditProject(p *Project) error {
	isExist, err := isProjectExist(p.ID, p.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrProjectNotExist{Name: p.Name}
	}
	_, err = x.ID(p.ID).AllCols().Update(p)
	return err
}

func GetProjectUID() (string, error) {
	return generate.GetRandomString(10)
}

func CreateProject(p *Project) error {
	uid, err := GetProjectUID()
	if err != nil {
		return err
	}
	p.UID = uid

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	isExist, err := isProjectExist(0, p.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrProjectNotExist{Name: p.Name}
	}

	if _, err = sess.Insert(p); err != nil {
		return err
	}
	return sess.Commit()
}

func (p *Project) APIFormat() *apiv1.Project {
	return &apiv1.Project{
		ID:            p.ID,
		UID:           p.UID,
		Name:          p.Name,
		Description:   p.Description,
		SlackURL:      p.SlackURL,
		SlackChannel:  p.SlackChannel,
		AckButton:     p.AckButton,
		ResolveButton: p.ResolveButton,
		SnoozeButton:  p.SnoozeButton,
		CreatedAt:     p.CreatedUnix.AsTime(),
		UpdatedAt:     p.UpdatedUnix.AsTime(),
	}
}

type SearchProjectsOptions struct {
	OrderBy SearchOrderBy
}

func (o *SearchProjectsOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	return cond
}

func SearchProjects(opts *SearchProjectsOptions) ([]*Project, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Project))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByAlphabetically
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var projects []*Project
	return projects, count, sess.Find(&projects)
}
