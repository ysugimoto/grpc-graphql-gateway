package types

import (
	"fmt"
	"strings"

	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/graphql"
)

type QuerySpec struct {
	Input  *Message
	Output *Message
	Option *graphql.GraphqlQuery
}

func (q *QuerySpec) GetExposeField() string {
	if q.Option.Response == nil {
		return ""
	}
	return q.Option.Response.GetExpose()
}

func (q *QuerySpec) BuildQuery() string {
	format := "%s(%s): %s%s"
	name := q.Option.GetName()
	args := q.ExtractArguments(q.Input)
	returnType := q.Output.Descriptor.GetName()
	// TODO: check via option
	sign := "!"
	return fmt.Sprintf(format, name, args, returnType, sign)
}

func (q *QuerySpec) ExtractArguments(input *Message) string {
	var args []string
	for _, f := range input.Descriptor.GetField() {
		sign := "!"
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			if opt.GetOptional() {
				sign = ""
			}
		}
		args = append(args, fmt.Sprintf(
			"$%s: %s%s",
			f.GetName(),
			convertProtoEnumToGraphqlType(f),
			sign,
		))
	}
	return strings.Join(args, ", ")
}
