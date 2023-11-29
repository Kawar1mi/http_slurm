package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Printf("request: %s method: %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
