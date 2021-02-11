package timestamppb

import (
	"github.com/graphql-go/graphql"
)

var (
	gql__type_Timestamp  *graphql.Object
	gql__input_Timestamp *graphql.InputObject
)

func Gql__type_Timestamp() *graphql.Object {
	if gql__type_Timestamp == nil {
		gql__type_Timestamp = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Timestamp",
			Fields: graphql.Fields{
				"seconds": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"nanos": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__type_Timestamp
}

func Gql__input_Timestamp() *graphql.InputObject {
	if gql__input_Timestamp == nil {
		gql__input_Timestamp = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Timestamp",
			Fields: graphql.InputObjectConfigFieldMap{
				"seconds": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"nanos": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__input_Timestamp
}
