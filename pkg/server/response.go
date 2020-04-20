package server

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Code() int
	Error() error
	WriteTo(w http.ResponseWriter) error
}

type NormalResponse struct {
	status  int
	body    []byte
	header  http.Header
	errCode string
	err     error
}

func (r *NormalResponse) Error() error {
	return r.err
}

func (r *NormalResponse) Code() int {
	return r.status
}

func (r *NormalResponse) WriteTo(w http.ResponseWriter) error {
	header := w.Header()
	for k, v := range r.header {
		header[k] = v
	}
	w.WriteHeader(r.status)
	// Fix: http: request method or response status code does not allow body
	if r.status == http.StatusNoContent {
		return nil
	}
	_, err := w.Write(r.body)
	return err
}

func (r *NormalResponse) Header(key, value string) *NormalResponse {
	r.header.Set(key, value)
	return r
}

// Empty пустой ответ
func Empty(status int) *NormalResponse {
	return Respond(status, nil)
}

// JSON ответ с JSON документом
func JSON(status int, body interface{}) *NormalResponse {
	return Respond(status, body).Header("Content-Type", "application/json")
}

// Error ответ с ошибкой
func Error(status int, code string, err error) *NormalResponse {
	data := make(map[string]interface{})

	switch status {
	case http.StatusForbidden:
		data["error"] = "Forbidden"
	case http.StatusNotFound:
		data["error"] = "Not Found"
	case http.StatusInternalServerError:
		data["error"] = "Internal Server Error"
	}

	if code != "" {
		data["error"] = code
	}

	if err != nil {
		data["error_description"] = err.Error()
	}

	resp := JSON(status, data)

	if err != nil {
		resp.errCode = code
		resp.err = err
	}

	return resp
}

func Respond(status int, body interface{}) *NormalResponse {
	var b []byte
	var err error
	switch t := body.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			return Error(http.StatusInternalServerError, "Response encode error", err)
		}
	}
	return &NormalResponse{
		body:   b,
		status: status,
		header: make(http.Header),
	}
}
