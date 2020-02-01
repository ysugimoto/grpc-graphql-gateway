package builder

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
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
			argField = fmt.Sprintf(`
				Args: graphql.FieldConfigArgument{
					%s
				},`,
				args,
			)
		}

		fields[i] = fmt.Sprintf(`
			"%s": &graphql.Field{
				Type: %s,
				%s
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},`,
			v.Option.GetName(),
			ext.MessageName(v.Input.Descriptor.GetName()),
			argField,
		)
	}

	return fmt.Sprintf(`
var %s = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		%s
	},
})`,
		ext.QueryName(),
		strings.Join(fields, "\n"),
	)
}
