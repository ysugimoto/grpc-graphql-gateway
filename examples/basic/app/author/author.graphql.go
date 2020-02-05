// This file is generated from proroc-gen-graphql, DO NOT EDIT!
package author

import (
	"github.com/graphql-go/graphql"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
)

var gql__type_ListAuthorsRequest = graphql.NewObject(graphql.ObjectConfig{
	Name:   "ListAuthorsRequest",
	Fields: graphql.Fields{},
}) // message ListAuthorsRequest in author/author.proto

var gql__type_Author = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
}) // message Author in author/author.proto

var gql__type_ListAuthorsResponse = graphql.NewObject(graphql.ObjectConfig{
	Name: "ListAuthorsResponse",
	Fields: graphql.Fields{
		"authors": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(gql__type_Author)),
		},
	},
}) // message ListAuthorsResponse in author/author.proto

var gql__type_GetAuthorRequest = graphql.NewObject(graphql.ObjectConfig{
	Name: "GetAuthorRequest",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
}) // message GetAuthorRequest in author/author.proto

// Create gRPC connection to host which specified via Service directive.
// If you registered handler via ReegisterXXXGraphqlHandler with your *grpc.ClientConn,
// this function won't be called.
func createAuthorServiceConnection() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:8080", grpc.WithInsecure())
}

// getQueryFields returns query target fields.
func getQueryFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		"authors": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(gql__type_Author)),

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if c == nil {
					var err error
					if c, err = createAuthorServiceConnection(); err != nil {
						return nil, err
					}
					defer func() {
						c.Close()
						c = nil
					}()
				}
				client := NewAuthorServiceClient(c)
				req := &ListAuthorsRequest{}
				resp, err := client.ListAuthors(p.Context, req)
				if err != nil {
					return nil, err
				}
				return resp.GetAuthors(), nil
			},
		},
		"author": &graphql.Field{
			Type: gql__type_Author,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if c == nil {
					var err error
					if c, err = createAuthorServiceConnection(); err != nil {
						return nil, err
					}
					defer func() {
						c.Close()
						c = nil
					}()
				}
				client := NewAuthorServiceClient(c)
				req := &GetAuthorRequest{}
				req.Name = p.Args["name"].(string)
				resp, err := client.GetAuthor(p.Context, req)
				if err != nil {
					return nil, err
				}
				return resp, nil
			},
		},
	}
}

// getMutationFields returns mutation target fields.
func getMutationFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{}
}

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you worried about open/close performance for each handling graphql request,
// then you can call RegisterAuthorHandler with *grpc.ClientConn manually.
func RegisterAuthorGraphql(mux *runtime.ServeMux) {
	RegisterAuthorGraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your client connection, so that we reuse that and never close connection inside.
// You need to close it maunally when appication will terminate.
func RegisterAuthorGraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) {
	mux.AddQueryField(getQueryFields(conn))
	mux.AddMutationField(getMutationFields(conn))
}
