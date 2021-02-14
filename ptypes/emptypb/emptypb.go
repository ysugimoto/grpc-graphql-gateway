package emptypb

import (
	"github.com/graphql-go/graphql"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Expose Google defined ptypes as this package types
type Empty = emptypb.Empty

var (
	gql__type_Empty  *graphql.Object
	gql__input_Empty *graphql.InputObject
)

func Gql__type_Empty() *graphql.Object {
	if gql__type_Empty == nil {
		gql__type_Empty = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Empty",
			Fields: graphql.Fields{
				"_": &graphql.Field{
					Type: graphql.Boolean,
				},
			},
		})
	}
	return gql__type_Empty
}

func Gql__input_Empty() *graphql.InputObject {
	if gql__input_Empty == nil {
		gql__input_Empty = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Empty",
			Fields: graphql.InputObjectConfigFieldMap{
				"_": &graphql.InputObjectFieldConfig{
					Type: graphql.Boolean,
				},
			},
		})
	}
	return gql__input_Empty
}
