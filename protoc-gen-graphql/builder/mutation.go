package builder

import (
	"fmt"
	"strings"

	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/types"
)

type Mutation struct {
	ms []*types.MutationSpec
}

func NewMutation(ms []*types.MutationSpec) *Mutation {
	return &Mutation{
		ms: ms,
	}
}

func (m *Mutation) BuildQuery() string {
	if len(m.ms) == 0 {
		return ""
	}
	lines := []string{`type Mutation {`}

	for _, v := range m.ms {
		lines = append(lines, fmt.Sprintf(
			"  %s(%s): %s%s",
			v.Option.GetName(),
			m.ExtractArguments(v.Input),
			v.Output.Descriptor.GetName(),
			// TODO: check via option
			"!",
		))
	}

	lines = append(lines, "}\n")

	return strings.Join(lines, "\n")
}

func (m *Mutation) ExtractArguments(input *types.Message) string {
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

func (m *Mutation) BuildProgram() string {
	return ""
}
