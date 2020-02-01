package main


import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	author "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/author"
book "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/book"
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

var gql_Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		
			"author": &graphql.Field{
				Type: gql_Type_GetAuthorRequest,
				
				Args: graphql.FieldConfigArgument{
					
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},

			"authors": &graphql.Field{
				Type: gql_Type_ListAuthorsRequest,
				
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},

			"book": &graphql.Field{
				Type: gql_Type_GetBookRequest,
				
				Args: graphql.FieldConfigArgument{
					
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},

			"books": &graphql.Field{
				Type: gql_Type_ListBooksRequest,
				
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
	},
})


type ErrorHandler func(errs []gqlerrors.FormattedError)

const (
	optNameErrorHandler = "errorhandler"
	optNameAllowCORS = "allowcors"
)

type Option struct {
	name string
	value interface{}
}

func WithErrorHandler(eh ErrorHandler) Option {
	return Option {
		name: optNameErrorHandler,
		value: eh,
	}
}

func WithCORS() Option {
	return Option {
		name: optNameAllowCORS,
		value: true,
	}
}

type GraphqlResolver struct {
	schema graphql.Schema
	errorHandler ErrorHandler
	allowCORS bool
}

func New(opts ...Option) *GraphqlResolver {
	var eh ErrorHandler
	var cors bool

	for _, o := range opts {
		switch o.name {
			case optNameErrorHandler:
				eh = o.value.(ErrorHandler)
			case optNameAllowCORS:
				cors = true
		}
	}

	return &GraphqlResolver {
		errorHandler: eh,
		allowCORS: cors,
	}
}

func (g *GraphqlResolver) SetCORSHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1728000")
}

func (g *GraphqlResolver) RespondError(w http.ResponseWriter, status int, message string) {
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
		g.SetCORSHeader(w, r)
	}
	var query string
	switch r.Method {
		case http.MethodOptions:
			g.RespondError(w, http.StatusNoContent, "")
			return
		case http.MethodPost:
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				g.RespondError(w, http.StatusBadRequest, "malformed request body")
				return
			}
			query = string(buf)
		case http.MethodGet:
			query = r.URL.Query().Get("query")
		default:
			g.RespondError(w, http.StatusBadRequest, "invalid request method: '" + r.Method + "'")
			return
	}

	result := graphql.Do(graphql.Params{
		Schema: g.schema,
		RequestString: query,
		Context: r.Context(),
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

