package server

import (
	"context"
	"io/ioutil"
	"os"

	apiv1 "github.com/incidenta/incidenta/pkg/api/v1"
	"github.com/incidenta/incidenta/pkg/models"
)

type testEnv struct {
	s        *Server
	filename string
	client   *apiv1.Client
}

func newTestEnv() (*testEnv, error) {
	t := &testEnv{
		s: New(Config{}),
	}
	f, err := ioutil.TempFile("", "incidenta")
	if err != nil {
		return nil, err
	}
	t.filename = f.Name()
	t.client = apiv1.NewClient(nil, "").WithRouter(t.s.httpServer.router)
	err = models.NewEngine(context.Background(), models.Config{
		Type:            "sqlite3",
		Schema:          "public",
		Path:            f.Name(),
		MaxIdleConn:     2,
		MaxOpenConn:     0,
		ConnMaxLifetime: 0,
	})
	return t, err
}

func (t *testEnv) Destroy() {
	models.Close()
	os.Remove(t.filename)
}
