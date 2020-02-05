package runtime

import (
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
)

// ServeMux is struct can execute graphql request via incoming HTTP request.
// This is inspired from grpc-gateway implementation, thanks!
type ServeMux struct {
	middlewares  []MiddlewareFunc
	ErrorHandler GraphqlErrorHandler

	queries   graphql.Fields
	mutations graphql.Fields
}

func NewServeMux(ms ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,

		queries:   graphql.Fields{},
		mutations: graphql.Fields{},
	}
}

func (s *ServeMux) AddQueryField(fields graphql.Fields) {
	for k, v := range fields {
		s.queries[k] = v
	}
}

func (s *ServeMux) AddMutationField(fields graphql.Fields) {
	for k, v := range fields {
		s.mutations[k] = v
	}
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

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: s.queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: s.mutations,
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
