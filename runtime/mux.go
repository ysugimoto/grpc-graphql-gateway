package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	wsgraphql "github.com/eientei/wsgraphql/v1"
	"github.com/eientei/wsgraphql/v1/compat/gorillaws"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"google.golang.org/grpc"
)

type (
	// MiddlewareFunc type definition
	MiddlewareFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error)
)

type GraphqlHandler interface {
	CreateConnection(context.Context) (*grpc.ClientConn, func(), error)
	GetMutations(*grpc.ClientConn) graphql.Fields
	GetQueries(*grpc.ClientConn) graphql.Fields
	GetSubscriptions(*grpc.ClientConn) graphql.Fields
}

// ServeMux is struct can execute graphql request via incoming HTTP request.
// This is inspired from grpc-gateway implementation, thanks!
// It also supports GraphQL subscription over WebSocket.
type ServeMux struct {
	middlewares  []MiddlewareFunc
	ErrorHandler GraphqlErrorHandler

	handlers []GraphqlHandler
}

// NewServeMux creates ServeMux pointer
func NewServeMux(ms ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,
		handlers:    make([]GraphqlHandler, 0),
	}
}

// Validate handler definition
func (s *ServeMux) validateHandler(h GraphqlHandler) error {
	queries := h.GetQueries(nil)
	mutations := h.GetMutations(nil)
	subscriptions := h.GetSubscriptions(nil)

	// If handler doesn't have any definitions, pass
	if len(queries) == 0 && len(mutations) == 0 && len(subscriptions) == 0 {
		return nil
	}

	schemaConfig := graphql.SchemaConfig{}
	if len(queries) > 0 {
		schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		})
	}
	if len(mutations) > 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		})
	}

	if len(subscriptions) > 0 {
		schemaConfig.Subscription = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Subscription",
			Fields: subscriptions,
		})
	}

	// Try to generate Schema and check error
	if _, err := graphql.NewSchema(schemaConfig); err != nil {
		return fmt.Errorf("Schema validation error: %s", err)
	}
	return nil
}

// AddHandler registers graphql handler which is built via plugin
func (s *ServeMux) AddHandler(h GraphqlHandler) error {
	if err := s.validateHandler(h); err != nil {
		return err
	}
	s.handlers = append(s.handlers, h)
	return nil
}

// ServeWs handles the GraphQL WebSocket upgrade and subscription traffic.
func (s *ServeMux) ServeWs(w http.ResponseWriter, r *http.Request, schema graphql.Schema) {
	// wrap your upgrader so it will negotiate *both* the old and the new sub-protocols
	upgrader := &websocket.Upgrader{
		Subprotocols: []string{
			wsgraphql.WebsocketSubprotocolGraphqlWS.String(),          // "graphql-ws"
			wsgraphql.WebsocketSubprotocolGraphqlTransportWS.String(), // "graphql-transport-ws"
		},
		CheckOrigin: func(r *http.Request) bool {
			// your origin logic here (or just return true)
			return true
		},
	}

	// create a single http.Handler that will do the upgrade & the message loop
	srv, err := wsgraphql.NewServer(
		schema,
		wsgraphql.WithUpgrader(gorillaws.Wrap(upgrader)), // hook into gorilla's upgrader
		// you can also add interceptors, logging, etc:
		// wsgraphql.WithInterceptor(myAuthInterceptor),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to initialize websocket server: %v", err), http.StatusInternalServerError)
		return
	}

	// hand off to wsgraphql — it will do the upgrade, subprotocol negotiation,
	// and then run your schema over whichever protocol the client picked
	srv.ServeHTTP(w, r)
}

// ServeHTTP implements http.Handler
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	for _, m := range s.middlewares {
		var err error
		ctx, err = m(ctx, w, r)
		if err != nil {
			ge := GraphqlError{}
			if me, ok := err.(*MiddlewareError); ok {
				ge.Message = me.Message
				ge.Extensions = map[string]interface{}{
					"code": me.Code,
				}
			} else {
				ge.Message = err.Error()
				ge.Extensions = map[string]interface{}{
					"code": "MIDDLEWARE_ERROR",
				}
			}
			respondResult(w, &graphql.Result{
				Errors: []GraphqlError{ge},
			})
			return
		}
	}
	// Build root schema from all handlers
	queries := graphql.Fields{}
	mutations := graphql.Fields{}
	subs := graphql.Fields{}

	for _, h := range s.handlers {
		c, closer, err := h.CreateConnection(ctx)
		if err != nil {
			respondResult(w, &graphql.Result{
				Errors: []gqlerrors.FormattedError{{
					Message:    "Failed to parse request: " + err.Error(),
					Extensions: map[string]interface{}{"code": "REQUEST_PARSE_ERROR"},
				},
				},
			})
			return
		}
		defer closer()

		for k, v := range h.GetQueries(c) {
			queries[k] = v
		}
		for k, v := range h.GetMutations(c) {
			mutations[k] = v
		}
		for k, v := range h.GetSubscriptions(c) {
			subs[k] = v
		}
	}

	schemaConfig := graphql.SchemaConfig{Query: buildObject("Query", queries), Mutation: buildObject("Mutation", mutations), Subscription: buildObject("Subscription", subs)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		respondResult(w, &graphql.Result{
			Errors: []gqlerrors.FormattedError{
				{
					Message:    "Failed to build schema: " + err.Error(),
					Extensions: map[string]interface{}{"code": "SCHEMA_GENERATION_ERROR"},
				},
			},
		})

		return
	}

	// Handle WebSocket upgrade for subscriptions
	if websocket.IsWebSocketUpgrade(r) && len(subs) > 0 {
		// Create subscription manager and WS handler
		s.ServeWs(w, r, schema)
		return
	}

	// Fallback to HTTP for queries and mutations
	req, err := parseRequest(r)
	if err != nil {
		respondResult(w, &graphql.Result{
			Errors: []gqlerrors.FormattedError{
				{
					Message:    "Failed to parse request: " + err.Error(),
					Extensions: map[string]interface{}{"code": "REQUEST_PARSE_ERROR"},
				},
			},
		})

		return
	}

	result := graphql.Do(graphql.Params{Schema: schema, RequestString: req.Query, VariableValues: req.Variables, Context: r.Context()})

	if len(result.Errors) > 0 {
		if s.ErrorHandler != nil {
			s.ErrorHandler(result.Errors)
		} else {
			defaultGraphqlErrorHandler(result.Errors)
		}
	}
	respondResult(w, result)
}

// Use adds more middlwares
func (s *ServeMux) Use(ms ...MiddlewareFunc) *ServeMux {
	s.middlewares = append(s.middlewares, ms...)
	return s
}

// helper to avoid duplication when building schema objects
func buildObject(name string, fields graphql.Fields) *graphql.Object {
	if len(fields) == 0 {
		return nil
	}
	return graphql.NewObject(graphql.ObjectConfig{Name: name, Fields: fields})
}

func respondResult(w http.ResponseWriter, result *graphql.Result) {
	out, _ := json.Marshal(result) // nolint: errcheck

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out) // nolint: errcheck
}
