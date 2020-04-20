package models

import (
	"fmt"
	"strings"

	"xorm.io/builder"

	"github.com/incidenta/incidenta/pkg/api"
	"github.com/incidenta/incidenta/pkg/timeutil"
)

type Template struct {
	ID      int64 `xorm:"pk autoincr"`
	Name    string
	Content string

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

func (t *Template) APIFormat(withContent bool) *api.Template {
	template := &api.Template{
		ID:          t.ID,
		Name:        t.Name,
		CreatedUnix: t.CreatedUnix.AsTime(),
		UpdatedUnix: t.UpdatedUnix.AsTime(),
	}
	if withContent {
		template.Content = t.Content
	}
	return template
}

type ErrTemplateNotExist struct {
	ID int64
}

func IsErrTemplateNotExist(err error) bool {
	_, ok := err.(ErrTemplateNotExist)
	return ok
}

func (e ErrTemplateNotExist) Error() string {
	return fmt.Sprintf("Template does not exist [id: %d]", e.ID)
}

func GetTemplateByID(id int64) (*Template, error) {
	t := new(Template)
	has, err := x.ID(id).Get(t)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrTemplateNotExist{ID: id}
	}
	return t, nil
}

func DeleteTemplate(t *Template) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := x.ID(t.ID).Delete(new(Template)); err != nil {
		return err
	}
	return sess.Commit()
}

func isTemplateExist(uid int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.
		Where("id!=?", uid).
		Get(&Template{Name: strings.ToLower(name)})
}

type ErrTemplateAlreadyExist struct {
	Name string
}

func IsErrTemplateAlreadyExist(err error) bool {
	_, ok := err.(ErrTemplateAlreadyExist)
	return ok
}

func (e ErrTemplateAlreadyExist) Error() string {
	return fmt.Sprintf("Template already exist [name: %s]", e.Name)
}

func EditTemplate(t *Template) error {
	isExist, err := isTemplateExist(t.ID, t.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrTemplateAlreadyExist{Name: t.Name}
	}
	_, err = x.ID(t.ID).AllCols().Update(t)
	return err
}

func CreateTemplate(t *Template) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	isExist, err := isTemplateExist(0, t.Name)
	if err != nil {
		return err
	} else if isExist {
		return ErrTemplateAlreadyExist{Name: t.Name}
	}
	t.Name = strings.ToLower(t.Name)
	if _, err = sess.Insert(t); err != nil {
		return err
	}
	return sess.Commit()
}

type SearchTemplatesOptions struct {
	OrderBy SearchOrderBy
}

func (o *SearchTemplatesOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	return cond
}

func SearchTemplates(opts *SearchTemplatesOptions) ([]*Template, int64, error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Template))
	if err != nil {
		return nil, 0, fmt.Errorf("count: %v", err)
	}

	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByAlphabetically
	}

	sess := x.Where(cond).OrderBy(opts.OrderBy.String())

	var templates []*Template
	return templates, count, sess.Find(&templates)
}
