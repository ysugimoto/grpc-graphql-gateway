// This file is generated from proroc-gen-graphql, DO NOT EDIT!
package book

import (
	"github.com/graphql-go/graphql"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"

	author "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/author"
)

var gql__type_ListBooksRequest = graphql.NewObject(graphql.ObjectConfig{
	Name:   "ListBooksRequest",
	Fields: graphql.Fields{},
}) // message ListBooksRequest in book/book.proto

var gql__type_Book = graphql.NewObject(graphql.ObjectConfig{
	Name: "Book",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"title": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.Field{
			Type: graphql.NewNonNull(gql__enum_BookType),
		},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(gql__type_Author),
		},
	},
}) // message Book in book/book.proto

var gql__enum_BookType = graphql.NewEnum(graphql.EnumConfig{
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
}) // message BookType in book/book.proto

var gql__type_Author = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
}) // message Author in author/author.proto

var gql__type_ListBooksResponse = graphql.NewObject(graphql.ObjectConfig{
	Name: "ListBooksResponse",
	Fields: graphql.Fields{
		"books": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(gql__type_Book)),
		},
	},
}) // message ListBooksResponse in book/book.proto

var gql__type_GetBookRequest = graphql.NewObject(graphql.ObjectConfig{
	Name: "GetBookRequest",
	Fields: graphql.Fields{
		// this is example comment for id field
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
}) // message GetBookRequest in book/book.proto

var gql__type_CreateBookRequest = graphql.NewObject(graphql.ObjectConfig{
	Name: "CreateBookRequest",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.Field{
			Type: graphql.NewNonNull(gql__enum_BookType),
		},
		"author": &graphql.Field{
			Type: graphql.NewNonNull(gql__type_Author),
		},
	},
}) // message CreateBookRequest in book/book.proto

// Create gRPC connection to host which specified via Service directive.
// If you registered handler via ReegisterXXXGraphqlHandler with your *grpc.ClientConn,
// this function won't be called.
func createBookServiceConnection() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:8080", grpc.WithInsecure())
}

// getQueryFields returns query target fields.
func getQueryFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		"books": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(gql__type_Book)),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if c == nil {
					var err error
					if c, err = createBookServiceConnection(); err != nil {
						return nil, err
					}
					defer func() {
						c.Close()
						c = nil
					}()
				}
				client := NewBookServiceClient(c)
				req := &ListBooksRequest{}
				resp, err := client.ListBooks(p.Context, req)
				if err != nil {
					return nil, err
				}
				return resp.GetBooks(), nil
			},
		},
		"book": &graphql.Field{
			Type: gql__type_Book,
			Args: graphql.FieldConfigArgument{
				// this is example comment for id field
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if c == nil {
					var err error
					if c, err = createBookServiceConnection(); err != nil {
						return nil, err
					}
					defer func() {
						c.Close()
						c = nil
					}()
				}
				client := NewBookServiceClient(c)
				req := &GetBookRequest{}
				req.Id = int64(p.Args["id"].(int))
				resp, err := client.GetBook(p.Context, req)
				if err != nil {
					return nil, err
				}
				return resp, nil
			},
		},
	}
}

var gql__input_CreateBookRequest = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateBookRequest",
	Fields: graphql.InputObjectConfigFieldMap{
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(gql__enum_BookType),
		},
		"author": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(gql__type_Author),
		},
	},
}) // message CreateBookRequest in book/book.proto

// getMutationFields returns mutation target fields.
func getMutationFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		"": &graphql.Field{
			Type: gql__type_Book,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: gql__input_CreateBookRequest,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if c == nil {
					var err error
					if c, err = createBookServiceConnection(); err != nil {
						return nil, err
					}
					defer func() {
						c.Close()
						c = nil
					}()
				}
				client := NewBookServiceClient(c)
				req := &CreateBookRequest{}
				req.Title = p.Args["title"].(string)
				req.Type = p.Args["type"].(BookType)
				req.Author = p.Args["author"].(*author.Author)
				resp, err := client.CreateBook(p.Context, req)
				if err != nil {
					return nil, err
				}
				return resp, nil
			},
		},
	}
}

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you worried about open/close performance for each handling graphql request,
// then you can call RegisterBookHandler with *grpc.ClientConn manually.
func RegisterBookGraphql(mux *runtime.ServeMux) {
	RegisterBookGraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your client connection, so that we reuse that and never close connection inside.
// You need to close it maunally when appication will terminate.
func RegisterBookGraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) {
	mux.AddQueryField(getQueryFields(conn))
	mux.AddMutationField(getMutationFields(conn))
}
