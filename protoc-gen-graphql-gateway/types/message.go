package types

import (
	"fmt"
	"strings"

	"path/filepath"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
)

type Message struct {
	Descriptor *descriptor.DescriptorProto
	File       *descriptor.FileDescriptorProto
}

func NewMessage(
	m *descriptor.DescriptorProto,
	f *descriptor.FileDescriptorProto,
) *Message {
	return &Message{
		Descriptor: m,
		File:       f,
	}
}

func (m *Message) MessageName() string {
	spec := strings.Split(m.Descriptor.GetName(), ".")
	return spec[len(spec)-1]
}

func (m *Message) GoPackageName() string {
	n := m.ProtoPackageName()
	if n == "" {
		n = "main"
	}
	if opts := m.File.GetOptions(); opts == nil {
		return n
	} else if gopkg := opts.GetGoPackage(); gopkg == "" {
		return n
	} else {
		return gopkg
	}
}

func (m *Message) ProtoPackageName() string {
	return m.File.GetPackage()
}

func (m *Message) Filename() string {
	return m.File.GetName()
}

func (m *Message) StructName(ptr bool) string {
	gopkg := m.GoPackageName()
	if gopkg == "main" {
		gopkg = ""
	} else {
		gopkg = filepath.Base(gopkg) + "."
	}
	var sign string
	if ptr {
		sign = "*"
	}
	return sign + gopkg + m.Descriptor.GetName()
}

func (m *Message) BuildQuery() string {
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
			convertProtoEnumToGraphqlType(f),
			sign,
		))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}
