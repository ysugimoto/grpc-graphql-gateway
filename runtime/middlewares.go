package runtime

import (
	"context"
	"net/http"
)

// Cors is middelware function to provide CORS headers to response headers
func Cors() MiddlewareFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
		w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		return ctx, nil
	}
}
