package types

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
)

type MutationSpec struct {
	Input  *Message
	Output *Message
	Option *graphql.GraphqlMutation
}

func (m *MutationSpec) GetExposeField() (*descriptor.FieldDescriptorProto, error) {
	if m.Option.Response == nil {
		return nil, nil
	}
	expose := m.Option.Response.GetExpose()
	for _, f := range m.Output.Descriptor.GetField() {
		if f.GetName() == expose {
			return f, nil
		}
	}
	return nil, fmt.Errorf(
		"expose field %s not found in message %s",
		expose,
		m.Output.Descriptor.GetName(),
	)
}

func (m *MutationSpec) BuildQuery() string {
	format := "%s(%s): %s%s"
	name := m.Option.GetName()
	args := m.ExtractArguments(m.Input)
	returnType := m.Output.Descriptor.GetName()
	// TODO: check via option
	sign := "!"
	return fmt.Sprintf(format, name, args, returnType, sign)
}

func (m *MutationSpec) ExtractArguments(input *Message) string {
	var args []string
	for _, f := range input.Descriptor.GetField() {
		sign := "!"
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			if o := opt.GetOptional(); o {
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
