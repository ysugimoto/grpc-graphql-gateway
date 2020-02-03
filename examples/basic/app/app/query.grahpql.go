package main

import (
	"errors"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"

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
		"title": &graphql.Field{
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

func createSchema(c *runtime.Connection) graphql.Schema {
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

						conn := c.Default
						if conn == nil {
							return nil, errors.New("failed to find default grpc connection")
						}
						client := author.NewAuthorServiceClient(conn)
						req := &author.GetAuthorRequest{}
						req.Name = p.Args["name"].(string)
						resp, err := client.GetAuthor(p.Context, req)
						if err != nil {
							return nil, err
						}
						return resp, nil
					},
				},
				"authors": &graphql.Field{
					Type: graphql.NewNonNull(graphql.NewList(gql_Type_Author)),

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {

						conn := c.Default
						if conn == nil {
							return nil, errors.New("failed to find default grpc connection")
						}
						client := author.NewAuthorServiceClient(conn)
						req := &author.ListAuthorsRequest{}
						resp, err := client.ListAuthors(p.Context, req)
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
						var keep bool
						conn := c.Find("localhost:8080")
						if conn == nil {
							if c, err := grpc.Dial("localhost:8080", grpc.WithInsecure()); err != nil {
								return nil, errors.New("failed to find grpc connection for 'localhost:8080'")
							} else {
								conn = c
							}
						} else {
							keep = true
						}
						if !keep {
							defer conn.Close()
						}
						client := book.NewBookServiceClient(conn)
						req := &book.GetBookRequest{}
						req.Id = int64(p.Args["id"].(int))
						resp, err := client.GetBook(p.Context, req)
						if err != nil {
							return nil, err
						}
						return resp, nil
					},
				},
				"books": &graphql.Field{
					Type: graphql.NewNonNull(graphql.NewList(gql_Type_Book)),

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var keep bool
						conn := c.Find("localhost:8080")
						if conn == nil {
							if c, err := grpc.Dial("localhost:8080", grpc.WithInsecure()); err != nil {
								return nil, errors.New("failed to find grpc connection for 'localhost:8080'")
							} else {
								conn = c
							}
						} else {
							keep = true
						}
						if !keep {
							defer conn.Close()
						}
						client := book.NewBookServiceClient(conn)
						req := &book.ListBooksRequest{}
						resp, err := client.ListBooks(p.Context, req)
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

func graphqlHandler(endpoint string, v interface{}) (runtime.GraphqlHandler, error) {
	var c *runtime.Connection
	if v == nil {
		c = runtime.NewConnection(nil)
	} else {
		switch t := v.(type) {
		case *grpc.ClientConn:
			c = runtime.NewConnection(t)
		case *runtime.Connection:
			c = t
		default:
			return nil, errors.New("invalid type conversion")
		}
	}

	schema := createSchema(c)

	return func(w http.ResponseWriter, r *http.Request) *graphql.Result {
		if r.URL.Path != endpoint {
			runtime.Respond(w, http.StatusNotFound, "endpoint not found")
			return nil
		}
		query, variables, err := runtime.ParseRequest(r)
		if err != nil {
			runtime.Respond(w, http.StatusBadRequest, err.Error())
			return nil
		}

		return graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  query,
			VariableValues: variables,
			Context:        r.Context(),
		})
	}, nil
}

func RegisterGraphqlHandler(mux *runtime.ServeMux, v interface{}, endpoint string) (err error) {
	mux.Handler, err = graphqlHandler(endpoint, v)
	return
}
