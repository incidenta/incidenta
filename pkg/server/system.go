package server

import "net/http"

func (h *HTTPServer) HealthzRequest(_ http.ResponseWriter, r *http.Request) Response {
	return Empty(204)
}
