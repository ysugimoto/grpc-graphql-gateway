package types

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
)

type QuerySpec struct {
	Input  *Message
	Output *Message
	Option *graphql.GraphqlQuery
}

func (q *QuerySpec) IsOutputOptional() bool {
	if q.Option.Response == nil {
		return false
	}
	return q.Option.Response.GetOptional()
}

func (q *QuerySpec) IsOutputRepeated() bool {
	if q.Option.Response == nil {
		return false
	}
	return q.Option.Response.GetRepeated()
}

func (q *QuerySpec) GetExposeField() (*descriptor.FieldDescriptorProto, error) {
	if q.Option.Response == nil {
		return nil, nil
	}
	expose := q.Option.Response.GetExpose()
	for _, f := range q.Output.Descriptor.GetField() {
		if f.GetName() == expose {
			return f, nil
		}
	}
	return nil, fmt.Errorf(
		"expose field %s not found in message %s",
		expose,
		q.Output.Descriptor.GetName(),
	)
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
