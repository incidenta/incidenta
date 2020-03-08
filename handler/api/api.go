package api

import (
	"fmt"
	"net/http"

	"github.com/incidenta/incidenta/alertmanager"
	"github.com/incidenta/incidenta/handler/api/render"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func Handler() http.Handler {
	c := cors.New(cors.Options{})
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m, err := alertmanager.ParseMessage(r)
		if err != nil {
			render.ErrorCode(w, fmt.Errorf("failed to parse"), http.StatusBadRequest)
			return
		}
		logrus.Infof("%#v", m)
	})
	return c.Handler(mux)
}
