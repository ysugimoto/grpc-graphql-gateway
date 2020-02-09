package runtime

import (
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
)

type SchemaBuilder interface {
	GetMutations() graphql.Fields
	GetQueries() graphql.Fields
}

// ServeMux is struct can execute graphql request via incoming HTTP request.
// This is inspired from grpc-gateway implementation, thanks!
type ServeMux struct {
	middlewares  []MiddlewareFunc
	ErrorHandler GraphqlErrorHandler

	handlers []SchemaBuilder
}

func NewServeMux(ms ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,
		handlers:    []SchemaBuilder{},
	}
}

func (s *ServeMux) AddHandler(h SchemaBuilder) {
	s.handlers = append(s.handlers, h)
}

// Use adds more middlwares which user defined
func (s *ServeMux) Use(ms ...MiddlewareFunc) *ServeMux {
	s.middlewares = append(s.middlewares, ms...)
	return s
}

// ServeHTTP implements http.Handler
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.middlewares {
		if err := m(w, r); err != nil {
			Respond(w, http.StatusBadRequest, "middleware error occured: "+err.Error())
			return
		}
	}

	queries := graphql.Fields{}
	mutations := graphql.Fields{}
	for _, h := range s.handlers {
		for k, v := range h.GetQueries() {
			queries[k] = v
		}
		for k, v := range h.GetMutations() {
			mutations[k] = v
		}
	}

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		}),
	})

	// TODO: assign variables
	query, variables, err := parseRequest(r)
	if err != nil {
		Respond(w, http.StatusBadRequest, err.Error())
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

	out, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
