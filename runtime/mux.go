package runtime

import (
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
)

type GraphqlHandler interface {
	CreateConnection() (*grpc.ClientConn, func(), error)
	GetMutations(*grpc.ClientConn) graphql.Fields
	GetQueries(*grpc.ClientConn) graphql.Fields
}

// ServeMux is struct can execute graphql request via incoming HTTP request.
// This is inspired from grpc-gateway implementation, thanks!
type ServeMux struct {
	middlewares  []MiddlewareFunc
	ErrorHandler GraphqlErrorHandler

	handlers []GraphqlHandler
}

func NewServeMux(ms ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,
		handlers:    make([]GraphqlHandler, 0),
	}
}

func (s *ServeMux) AddHandler(h GraphqlHandler) {
	s.handlers = append(s.handlers, h)
}

// Use adds more middlwares
func (s *ServeMux) Use(ms ...MiddlewareFunc) *ServeMux {
	s.middlewares = append(s.middlewares, ms...)
	return s
}

// ServeHTTP implements http.Handler
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.middlewares {
		if err := m(w, r); err != nil {
			http.Error(w, "middleware error occured: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	queries := graphql.Fields{}
	mutations := graphql.Fields{}
	for _, h := range s.handlers {
		c, closer, err := h.CreateConnection()
		if err != nil {
			http.Error(w, "failed to create grpc connection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer closer()

		for k, v := range h.GetQueries(c) {
			queries[k] = v
		}
		for k, v := range h.GetMutations(c) {
			mutations[k] = v
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		}),
	})
	if err != nil {
		http.Error(w, "failed to build graphql schema: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: assign variables
	query, variables, err := parseRequest(r)
	if err != nil {
		http.Error(w, "failed to parse request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        r.Context(),
	})

	if len(result.Errors) > 0 {
		if s.ErrorHandler != nil {
			s.ErrorHandler(result.Errors)
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "failed to marshal response JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
