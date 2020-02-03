package middleware

import (
	"net/http"
)

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request) error
