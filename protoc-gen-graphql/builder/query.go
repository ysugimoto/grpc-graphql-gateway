package builder

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/types"
)

type Query struct {
	qs []*types.QuerySpec
}

func NewQuery(qs []*types.QuerySpec) *Query {
	return &Query{
		qs: qs,
	}
}

func (q *Query) BuildQuery() string {
	if len(q.qs) == 0 {
		return ""
	}
	lines := []string{`type Query {`}

	for _, v := range q.qs {
		var fieldName, sign string
		if f, _ := v.GetExposeField(); f != nil {
			fieldName = ext.ConvertGraphqlType(f)
			var optional bool
			if opts := ext.GraphqlFieldExtension(f); opts != nil {
				optional = opts.GetOptional()
			}
			if !optional {
				fieldName += "!"
			}
			if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				fieldName = "[" + fieldName + "]"
				if !v.IsOutputOptional() {
					sign = "!"
				}
			}
		} else {
			fieldName = v.Output.Descriptor.GetName()
			if !v.IsOutputOptional() {
				sign = "!"
			}
		}

		lines = append(lines, fmt.Sprintf(
			"  %s(%s): %s%s",
			v.Option.GetName(),
			q.ExtractArguments(v.Input),
			fieldName,
			sign,
		))
	}

	lines = append(lines, "}\n")

	return strings.Join(lines, "\n")
}

func (q *Query) ExtractArguments(input *types.Message) string {
	var args []string
	for _, f := range input.Descriptor.GetField() {
		sign := "!"
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			if opt.GetOptional() {
				sign = ""
			}
		}
		args = append(args, fmt.Sprintf(
			"%s: %s%s",
			f.GetName(),
			ext.ConvertGraphqlType(f),
			sign,
		))
	}
	return strings.Join(args, ", ")
}

func (q *Query) BuildProgram() string {
	fields := make([]string, len(q.qs))

	for i, v := range q.qs {
		args := NewArgument(&types.ArgumentSpec{
			Message: v.Input,
		}).BuildProgram()

		var argField string
		if args != "" {
			argField = strings.TrimSpace(fmt.Sprintf(`
				Args: graphql.FieldConfigArgument{
					%s
				},`,
				args,
			))
		}
		var typeName string
		expose, _ := v.GetExposeField()
		if expose != nil {
			typeName = ext.ConvertGoType(expose)
			if expose.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				typeName = "graphql.NewList(" + typeName + ")"
			}
			var optional bool
			if opt := ext.GraphqlFieldExtension(expose); opt != nil {
				optional = opt.GetOptional()
			}
			if !optional {
				typeName = "graphql.NewNonNull(" + typeName + ")"
			}
		} else {
			typeName = ext.MessageName(v.Output.Descriptor.GetName())
			if v.IsOutputRepeated() {
				typeName = "graphql.NewList(" + typeName + ")"
			}
			if !v.IsOutputOptional() {
				typeName = "graphql.NewNonNull(" + typeName + ")"
			}
		}

		fields[i] = strings.TrimSpace(fmt.Sprintf(`
			"%s": &graphql.Field{
				Type: %s,
				%s
				Resolve: %s,
			},`,
			v.Option.GetName(),
			typeName,
			argField,
			NewResolver(v).BuildProgram(),
		))
	}

	return fmt.Sprintf(`
func createSchema(c *runtime.Connection) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				%s
			},
		}),
	})
	return schema
}`,
		strings.Join(fields, "\n"),
	)
}
