package builder

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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
		fieldType := ext.ConvertGraphqlType(f)
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			fieldType = "[" + fieldType + "]"
		}
		var optional bool
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			optional = opt.GetOptional()
		}
		if !optional {
			fieldType += "!"
		}
		lines = append(lines, fmt.Sprintf(
			"  %s: %s",
			f.GetName(),
			fieldType,
		))
	}
	lines = append(lines, "}\n")

	return strings.Join(lines, "\n")
}

func (b *Message) BuildProgram() string {
	fields := make([]string, len(b.m.Descriptor.GetField()))
	for i, f := range b.m.Descriptor.GetField() {
		fieldType := ext.ConvertGoType(f)
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			fieldType = "graphql.NewList(" + fieldType + ")"
		}

		var optional bool
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			optional = opt.GetOptional()
		}
		if !optional {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}

		fields[i] = strings.TrimSpace(fmt.Sprintf(`
			"%s": &graphql.Field{
				Type: %s,
			},`,
			f.GetName(),
			fieldType,
		))
	}

	return fmt.Sprintf(`
var %s = graphql.NewObject(graphql.ObjectConfig{
	Name: "%s",
	Fields: graphql.Fields{
		%s
	},
})`,
		ext.MessageName(b.m.Descriptor.GetName()),
		b.m.Descriptor.GetName(),
		strings.TrimSpace(strings.Join(fields, "\n")),
	)
}
