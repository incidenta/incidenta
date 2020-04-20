package server

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func wrapper(handler func(w http.ResponseWriter, r *http.Request) Response) func(w http.ResponseWriter, r *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		reqURL := r.URL.String()
		start := time.Now()
		defer func() {
			logrus.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     reqURL,
				"duration": time.Since(start).String(),
			}).Debug("Request complete")
		}()

		obj := handler(w, r)
		if obj.Error() != nil {
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   reqURL,
				"code":   obj.Code(),
			}).WithError(obj.Error()).Error("Request failed")
		}

		err := obj.WriteTo(w)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   reqURL,
			}).WithError(err).Error("Error writing to response")
		}
	}
	return f
}
