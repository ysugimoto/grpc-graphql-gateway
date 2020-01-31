package builder

import (
	"fmt"
	"strings"

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
		lines = append(lines, fmt.Sprintf(
			"  %s(%s): %s%s",
			v.Option.GetName(),
			q.ExtractArguments(v.Input),
			v.Output.Descriptor.GetName(),
			// TODO: check via option
			"!",
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
	return ""
}
