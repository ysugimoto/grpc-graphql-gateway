package runtime

import (
	"context"
	"errors"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareError(t *testing.T) {
	t.Run("Set message field on common error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := NewServeMux()
			mux.Use(func(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
				return ctx, errors.New("error")
			})
			mux.ServeHTTP(w, r)
		}))
		defer srv.Close()
		resp, err := http.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var r graphql.Result
		err = json.NewDecoder(resp.Body).Decode(&r)
		assert.NoError(t, err)
		assert.Nil(t, r.Data)
		if assert.Len(t, r.Errors, 1) {
			e := r.Errors[0]
			assert.Equal(t, "error", e.Message)
			if assert.Len(t, e.Extensions, 1) {
				assert.Equal(t, "MIDDLEWARE_ERROR", e.Extensions["code"])
			}
		}
	})

	t.Run("Set message field on MiddewareError", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := NewServeMux()
			mux.Use(func(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
				return ctx, NewMiddlewareError("CUSTOM_CODE", "CUSTOM_MESSAGE")
			})
			mux.ServeHTTP(w, r)
		}))
		defer srv.Close()
		resp, err := http.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var r graphql.Result
		err = json.NewDecoder(resp.Body).Decode(&r)
		assert.NoError(t, err)
		assert.Nil(t, r.Data)
		if assert.Len(t, r.Errors, 1) {
			e := r.Errors[0]
			assert.Equal(t, "CUSTOM_MESSAGE", e.Message)
			if assert.Len(t, e.Extensions, 1) {
				assert.Equal(t, "CUSTOM_CODE", e.Extensions["code"])
			}
		}
	})
}
