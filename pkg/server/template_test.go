package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
)

func TestHTTPServer_TemplateCreateRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name: "only-name",
	})
	assert.Error(t, err)

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Content: "only-content",
	})
	assert.Error(t, err)

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name",
		Content: "content",
	})
	assert.NoError(t, err)

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name",
		Content: "content",
	})
	assert.Error(t, err)
}

func TestHTTPServer_TemplateGetRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	templateOnCreate, _, err := te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name",
		Content: "content",
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, templateOnCreate.ID, int64(1))

	template, _, err := te.client.Templates.Get(templateOnCreate.ID)
	assert.NoError(t, err)
	assert.Equal(t, template.ID, templateOnCreate.ID)
	assert.Equal(t, template.Name, templateOnCreate.Name)
}

func TestHTTPServer_TemplateDeleteRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	template, _, err := te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name",
		Content: "content",
	})
	assert.NoError(t, err)

	_, err = te.client.Templates.Delete(template.ID)
	assert.NoError(t, err)

	_, err = te.client.Templates.Delete(template.ID)
	assert.Error(t, err)
}

func TestHTTPServer_TemplateListRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name1",
		Content: "content",
	})
	assert.NoError(t, err)

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name2",
		Content: "content",
	})
	assert.NoError(t, err)

	templates, _, err := te.client.Templates.List()
	assert.NoError(t, err)
	assert.Equal(t, len(templates), 2)
}

func TestHTTPServer_TemplateEditRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	_, _, err = te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name1",
		Content: "content",
	})
	assert.NoError(t, err)

	templatePre, _, err := te.client.Templates.Create(&apiv1.TemplateCreateOptions{
		Name:    "name2",
		Content: "content",
	})
	assert.NoError(t, err)

	name := "name3"
	content := "content3"
	templatePost, _, err := te.client.Templates.Edit(templatePre.ID, &apiv1.TemplateEditOptions{
		Name:    &name,
		Content: &content,
	})
	assert.NoError(t, err)
	assert.Equal(t, templatePost.ID, templatePre.ID)
	assert.Equal(t, templatePost.Name, name)
	assert.Equal(t, templatePost.Content, content)

	name = "name1"
	_, _, err = te.client.Templates.Edit(templatePre.ID, &apiv1.TemplateEditOptions{
		Name: &name,
	})
	assert.Error(t, err)
}
