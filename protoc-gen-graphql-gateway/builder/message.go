package builder

import (
	"fmt"
	"strings"

	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Message struct {
	m *types.Message
}

func NewMessage(m *types.Message) *Message {
	return &Message{
		m: m,
	}
}

func (b *Message) BuildQuery() string {
	m := b.m

	lines := []string{
		fmt.Sprintf("# message %s in %s", m.MessageName(), m.Filename()),
		fmt.Sprintf("type %s {", m.Descriptor.GetName()),
	}
	for _, f := range m.Descriptor.GetField() {
		sign := "!"
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			if opt.GetOptional() {
				sign = ""
			}
		}
		lines = append(lines, fmt.Sprintf(
			"  %s: %s%s",
			f.GetName(),
			ext.ConvertGraphqlType(f),
			sign,
		))
	}
	lines = append(lines, "}\n")

	return strings.Join(lines, "\n")
}

func (m *Message) BuildProgram() string {
	return ""
}
