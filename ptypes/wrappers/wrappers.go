package wrappers

import (
	"github.com/graphql-go/graphql"
)

var (
	gql__type_DoubleValue  *graphql.Object
	gql__type_FloatValue   *graphql.Object
	gql__type_Int64Value   *graphql.Object
	gql__type_Uint64Value  *graphql.Object
	gql__type_Int32Value   *graphql.Object
	gql__type_Uint32Value  *graphql.Object
	gql__type_BoolValue    *graphql.Object
	gql__type_StringValue  *graphql.Object
	gql__input_DoubleValue *graphql.InputObject
	gql__input_FloatValue  *graphql.InputObject
	gql__input_Int64Value  *graphql.InputObject
	gql__input_Uint64Value *graphql.InputObject
	gql__input_Int32Value  *graphql.InputObject
	gql__input_Uint32Value *graphql.InputObject
	gql__input_BoolValue   *graphql.InputObject
	gql__input_StringValue *graphql.InputObject
)

func Gql__type_DoubleValue() *graphql.Object {
	if gql__type_DoubleValue == nil {
		gql__type_DoubleValue = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_DoubleValue",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
		})
	}
	return gql__type_DoubleValue
}

func Gql__type_FloatValue() *graphql.Object {
	if gql__type_FloatValue == nil {
		gql__type_FloatValue = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_FloatValue",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
		})
	}
	return gql__type_FloatValue
}

func Gql__type_Int64Value() *graphql.Object {
	if gql__type_Int64Value == nil {
		gql__type_Int64Value = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_Int64Value",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__type_Int64Value
}

func Gql__type_Uint64Value() *graphql.Object {
	if gql__type_Uint64Value == nil {
		gql__type_Uint64Value = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_Uint64Value",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__type_Uint64Value
}

func Gql__type_Int32Value() *graphql.Object {
	if gql__type_Int32Value == nil {
		gql__type_Int32Value = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_Int32Value",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__type_Int32Value
}

func Gql__type_Uint32Value() *graphql.Object {
	if gql__type_Uint32Value == nil {
		gql__type_Uint64Value = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_Uint32Value",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__type_Uint32Value
}

func Gql__type_BoolValue() *graphql.Object {
	if gql__type_BoolValue == nil {
		gql__type_BoolValue = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_BoolValue",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
			},
		})
	}
	return gql__type_BoolValue
}

func Gql__type_StringValue() *graphql.Object {
	if gql__type_StringValue == nil {
		gql__type_BoolValue = graphql.NewObject(graphql.ObjectConfig{
			Name: "Google_type_Wrappers_StringValue",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		})
	}
	return gql__type_StringValue
}

func Gql__input_DoubleValue() *graphql.InputObject {
	if gql__input_DoubleValue == nil {
		gql__input_DoubleValue = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_DoubleValue",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
		})
	}
	return gql__input_DoubleValue
}

func Gql__input_FloatValue() *graphql.InputObject {
	if gql__input_FloatValue == nil {
		gql__input_FloatValue = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_FloatValue",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
		})
	}
	return gql__input_FloatValue
}

func Gql__input_Int64Value() *graphql.InputObject {
	if gql__input_Int64Value == nil {
		gql__input_Int64Value = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_Int64Value",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__input_Int64Value
}

func Gql__input_Uint64Value() *graphql.InputObject {
	if gql__input_Uint64Value == nil {
		gql__input_Uint64Value = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_Uint64Value",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__input_Uint64Value
}

func Gql__input_Int32Value() *graphql.InputObject {
	if gql__input_Int32Value == nil {
		gql__input_Int32Value = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_Int32Value",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__input_Int32Value
}

func Gql__input_Uint32Value() *graphql.InputObject {
	if gql__input_Uint32Value == nil {
		gql__input_Uint32Value = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_Uint32Value",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		})
	}
	return gql__input_Uint32Value
}

func Gql__input_BoolValue() *graphql.InputObject {
	if gql__input_BoolValue == nil {
		gql__input_BoolValue = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_BoolValue",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
			},
		})
	}
	return gql__input_BoolValue
}

func Gql__input_StringValue() *graphql.InputObject {
	if gql__input_StringValue == nil {
		gql__input_StringValue = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "Google_input_Wrappers_StringValue",
			Fields: graphql.InputObjectConfigFieldMap{
				"value": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		})
	}
	return gql__input_StringValue
}
