package middleware

import (
	"net/http"
)

func Cors() MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		return nil
	}
}
