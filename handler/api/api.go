package api

import (
	"net/http"

	"github.com/rs/cors"
)

func Handler() http.Handler {
	c := cors.New(cors.Options{})
	mux := http.NewServeMux()
	return c.Handler(mux)
}
