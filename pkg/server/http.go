package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

// swagger:response GenericError
type swaggerGenericError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type HTTPServer struct {
	config Config
	router *mux.Router
}

func NewHTTPServer(c Config) *HTTPServer {
	h := &HTTPServer{
		config: c,
		router: mux.NewRouter(),
	}
	h.addRoutes()
	return h
}

func (h *HTTPServer) addRoutes() {
	r := h.router.PathPrefix("/v1").Subrouter()

	// Project
	r.HandleFunc("/projects", wrapper(h.ProjectListRequest)).Methods("GET")
	r.HandleFunc("/project", wrapper(h.ProjectCreateRequest)).Methods("POST")
	r.HandleFunc("/project/{project_id}", wrapper(h.ProjectGetRequest)).Methods("GET")
	r.HandleFunc("/project/{project_id}", wrapper(h.ProjectEditRequest)).Methods("POST")
	r.HandleFunc("/project/{project_id}", wrapper(h.ProjectDeleteRequest)).Methods("DELETE")
	r.HandleFunc("/project/{project_id}/alerts", wrapper(h.ProjectAlertsRequest)).Methods("GET")

	// Event
	r.HandleFunc("/event/{event_id}", wrapper(h.EventGetRequest)).Methods("GET")
	r.HandleFunc("/event/{event_id}", wrapper(h.EventDeleteRequest)).Methods("DELETE")

	// Alert
	r.HandleFunc("/alerts", wrapper(h.AlertListRequest)).Methods("GET")
	r.HandleFunc("/alert/{alert_id}", wrapper(h.AlertGetRequest)).Methods("GET")
	r.HandleFunc("/alert/{alert_id}", wrapper(h.AlertDeleteRequest)).Methods("DELETE")
	r.HandleFunc("/alert/{alert_id}/events", wrapper(h.AlertEventsRequest)).Methods("GET")

	// Alertmanager
	r.HandleFunc("/integrations/alertmanager/{project_uid}", wrapper(h.AlertmanagerRequest))

	// Telemetry
	h.router.Handle("/metrics", promhttp.Handler())
	h.router.HandleFunc("/healthz", wrapper(h.HealthzRequest))

	// Static
	h.router.PathPrefix("/").
		Handler(http.FileServer(http.Dir(h.config.StaticAssets)))
}

func (h *HTTPServer) Serve() error {
	return http.ListenAndServe(h.config.Addr, h.router)
}
