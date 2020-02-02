package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	author "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/author"
	book "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/book"
	"google.golang.org/grpc"
)

var gql_Type_Author = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var gql_Type_Book = graphql.NewObject(graphql.ObjectConfig{
	Name: "Book",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.Field{
			Type: graphql.NewNonNull(gql_Enum_BookType),
		},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(gql_Type_Author),
		},
	},
})

var gql_Enum_BookType = graphql.NewEnum(graphql.EnumConfig{
	Name: "BookType",
	Values: graphql.EnumValueConfigMap{
		"JAVASCRIPT": &graphql.EnumValueConfig{
			Value: 0,
		},
		"ECMASCRIPT": &graphql.EnumValueConfig{
			Value: 1,
		},
		"GIT": &graphql.EnumValueConfig{
			Value: 2,
		},
		"ASP_DOT_NET": &graphql.EnumValueConfig{
			Value: 3,
		},
	},
})

func createSchema(conn *grpc.ClientConn) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"author": &graphql.Field{
					Type: graphql.NewNonNull(gql_Type_Author),
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						client := author.NewAuthorServiceClient(conn)
						resp, err := client.GetAuthor(
							p.Context,
							&author.GetAuthorRequest{},
						)
						if err != nil {
							return nil, err
						}
						return resp, nil
					},
				},
				"authors": &graphql.Field{
					Type: graphql.NewNonNull(graphql.NewList(gql_Type_Author)),

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						client := author.NewAuthorServiceClient(conn)
						resp, err := client.ListAuthors(
							p.Context,
							&author.ListAuthorsRequest{},
						)
						if err != nil {
							return nil, err
						}
						return resp.GetAuthors(), nil
					},
				},
				"book": &graphql.Field{
					Type: graphql.NewNonNull(gql_Type_Book),
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						client := book.NewBookServiceClient(conn)
						resp, err := client.GetBook(
							p.Context,
							&book.GetBookRequest{},
						)
						if err != nil {
							return nil, err
						}
						return resp, nil
					},
				},
				"books": &graphql.Field{
					Type: graphql.NewNonNull(graphql.NewList(gql_Type_Book)),

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						client := book.NewBookServiceClient(conn)
						resp, err := client.ListBooks(
							p.Context,
							&book.ListBooksRequest{},
						)
						if err != nil {
							return nil, err
						}
						return resp.GetBooks(), nil
					},
				},
			},
		}),
	})
	return schema
}

type ErrorHandler func(errs []gqlerrors.FormattedError)

const (
	optNameErrorHandler = "errorhandler"
	optNameAllowCORS    = "allowcors"
	optNameGrpcOption   = "grpcoption"
)

type Option struct {
	name  string
	value interface{}
}

func WithErrorHandler(eh ErrorHandler) Option {
	return Option{
		name:  optNameErrorHandler,
		value: eh,
	}
}

func WithCORS() Option {
	return Option{
		name:  optNameAllowCORS,
		value: true,
	}
}

func WithGrpcOption(opts ...grpc.DialOption) Option {
	return Option{
		name:  optNameGrpcOption,
		value: opts,
	}
}

type GraphqlResolver struct {
	errorHandler ErrorHandler
	allowCORS    bool
	grpcOptions  []grpc.DialOption
}

func New(opts ...Option) *GraphqlResolver {
	var eh ErrorHandler
	var cors bool
	var grpcOptions []grpc.DialOption

	for _, o := range opts {
		switch o.name {
		case optNameErrorHandler:
			eh = o.value.(ErrorHandler)
		case optNameAllowCORS:
			cors = true
		case optNameGrpcOption:
			grpcOptions = value.([]grpc.DialOption)
		}
	}

	return &GraphqlResolver{
		errorHandler: eh,
		allowCORS:    cors,
		grpcOptions:  grpcOptions,
	}
}

func corsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1728000")
}

func respondError(w http.ResponseWriter, status int, message string) {
	m := []byte(message)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(m)))
	w.WriteHeader(status)
	if len(m) > 0 {
		w.Write(m)
	}
}

func (g *GraphqlResolver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if g.allowCORS {
		corsHeader(w, r)
	}
	var query string
	switch r.Method {
	case http.MethodOptions:
		respondError(w, http.StatusNoContent, "")
		return
	case http.MethodPost:
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusBadRequest, "malformed request body")
			return
		}
		query = string(buf)
	case http.MethodGet:
		query = r.URL.Query().Get("query")
	default:
		respondError(w, http.StatusBadRequest, "invalid request method: '"+r.Method+"'")
		return
	}

	conn, err := grpc.Dial("localhost:50051", g.grpcOptions...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()

	result := graphql.Do(graphql.Params{
		Schema:        createSchema(conn),
		RequestString: query,
		Context:       r.Context(),
	})
	if len(result.Errors) > 0 {
		if g.errorHandler != nil {
			g.errorHandler(result.Errors)
		}
	}
	out, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
