package render

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

var indent bool

func init() {
	indent, _ = strconv.ParseBool(
		os.Getenv("HTTP_JSON_INDENT"),
	)
}

func ErrorCode(w http.ResponseWriter, err error, status int) {
	JSON(w, map[string]string{"message": err.Error()}, status)
}

func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if indent {
		enc.SetIndent("", "  ")
	}
	enc.Encode(v)
}
