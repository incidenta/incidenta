package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
)

func TestHTTPServer_ProjectCreateRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{})
	assert.Error(t, err)

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name: "name",
	})
	assert.Error(t, err)

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name",
		SlackChannel: "channel",
	})
	assert.Error(t, err)

	p, _, err := te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "name")
	assert.Greater(t, len(p.UID), 0)

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.Error(t, err)
}

func TestHTTPServer_ProjectGetRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	pCreated, _, err := te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)

	_, _, err = te.client.Projects.Get(100)
	assert.Error(t, err)

	pRequester, _, err := te.client.Projects.Get(pCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, pRequester.ID, pCreated.ID)
	assert.Equal(t, pRequester.Name, pCreated.Name)
}

func TestHTTPServer_ProjectListRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name1",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name2",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)

	projects, _, err := te.client.Projects.List()
	assert.NoError(t, err)
	assert.Equal(t, len(projects), 2)
}

func TestHTTPServer_ProjectDeleteRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	project, _, err := te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)

	_, err = te.client.Projects.Delete(project.ID)
	assert.NoError(t, err)

	_, err = te.client.Projects.Delete(project.ID)
	assert.Error(t, err)

	_, _, err = te.client.Projects.Get(project.ID)
	assert.Error(t, err)
}

func TestHTTPServer_ProjectEditRequest(t *testing.T) {
	te, err := newTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer te.Destroy()

	name := "name1"
	desc := "foobar"

	_, _, err = te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         name,
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)

	project, _, err := te.client.Projects.Create(&apiv1.ProjectCreateOptions{
		Name:         "name2",
		SlackChannel: "channel",
		SlackURL:     "url",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, project.Description, desc)

	_, _, err = te.client.Projects.Edit(project.ID+1, &apiv1.ProjectEditOptions{
		Description: &desc,
	})
	assert.Error(t, err)

	projectsPost, _, err := te.client.Projects.Edit(project.ID, &apiv1.ProjectEditOptions{
		Description: &desc,
	})
	assert.NoError(t, err)
	assert.Equal(t, projectsPost.Description, desc)

	_, _, err = te.client.Projects.Edit(project.ID, &apiv1.ProjectEditOptions{
		Name: &name,
	})
	assert.Error(t, err)
}
