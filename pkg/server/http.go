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

	// Receiver
	r.HandleFunc("/receivers", wrapper(h.ReceiverListRequest)).Methods("GET")
	r.HandleFunc("/receiver", wrapper(h.ReceiverCreateRequest)).Methods("POST")
	r.HandleFunc("/receiver/{{ receiver_id }}", wrapper(h.ReceiverGetRequest)).Methods("GET")
	r.HandleFunc("/receiver/{{ receiver_id }}", wrapper(h.ReceiverEditRequest)).Methods("POST")
	r.HandleFunc("/receiver/{{ receiver_id }}", wrapper(h.ReceiverDeleteRequest)).Methods("DELETE")
	r.HandleFunc("/receiver/{{ receiver_id }}/alerts", wrapper(h.ReceiverAlertsRequest)).Methods("GET")

	// Log
	r.HandleFunc("/log/{{ log_id }}", wrapper(h.LogGetRequest)).Methods("GET")
	r.HandleFunc("/log/{{ log_id }}", wrapper(h.LogDeleteRequest)).Methods("DELETE")

	// Alert
	r.HandleFunc("/alerts", wrapper(h.AlertListRequest)).Methods("GET")
	r.HandleFunc("/alert/{{ alert_id }}", wrapper(h.AlertGetRequest)).Methods("GET")
	r.HandleFunc("/alert/{{ alert_id }}", wrapper(h.AlertDeleteRequest)).Methods("DELETE")
	r.HandleFunc("/alert/{{ alert_id }}/logs", wrapper(h.AlertLogsRequest)).Methods("GET")

	// Template
	r.HandleFunc("/templates", wrapper(h.TemplateListRequest)).Methods("GET")
	r.HandleFunc("/template", wrapper(h.TemplateCreateRequest)).Methods("POST")
	r.HandleFunc("/template/{{ template_id }}", wrapper(h.TemplateGetRequest)).Methods("GET")
	r.HandleFunc("/template/{{ template_id }}", wrapper(h.TemplateEditRequest)).Methods("POST")
	r.HandleFunc("/template/{{ template_id }}", wrapper(h.TemplateDeleteRequest)).Methods("DELETE")

	// Alertmanager
	r.HandleFunc("/alertmanager/webhook", wrapper(h.AlertmanagerWebhookRequest))

	// Telemetry
	h.router.Handle("/metrics", promhttp.Handler())

	// Static
	h.router.PathPrefix("/").
		Handler(http.FileServer(http.Dir("./static")))
}

func (h *HTTPServer) Serve() error {
	return http.ListenAndServe(h.config.Addr, h.router)
}
