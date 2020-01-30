package types

import (
	"fmt"
	"strings"

	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/graphql"
)

type MutationSpec struct {
	Input  *Message
	Output *Message
	Option *graphql.GraphqlMutation
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
