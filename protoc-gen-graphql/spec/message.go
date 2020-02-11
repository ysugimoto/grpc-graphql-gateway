package spec

import (
	"strings"

	"path/filepath"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Message spec wraps DescriptorProto
type Message struct {
	descriptor *descriptor.DescriptorProto
	*File

	prefix []string
	paths  []int
}

func NewMessage(
	m *descriptor.DescriptorProto,
	f *File,
	prefix []string,
	paths ...int,
) *Message {
	return &Message{
		descriptor: m,
		File:       f,
		prefix:     prefix,
		paths:      paths,
	}
}

func (m *Message) Fields() []*Field {
	var fields []*Field
	for i, f := range m.descriptor.GetField() {
		paths := make([]int, len(m.paths))
		copy(paths, m.paths)
		fields = append(fields, NewField(f, m.File, append(paths, 2, i)...))
	}
	return fields
}

func (m *Message) Comment() string {
	return m.File.getComment(m.paths)
}

func (m *Message) Name() string {
	var p string
	if len(m.prefix) > 0 {
		p = strings.Join(m.prefix, ".") + "."
	}
	return p + m.descriptor.GetName()
}

func (m *Message) TypeName() string {
	var p string
	if len(m.prefix) > 0 {
		p = strings.Join(m.prefix, "_") + "_"
	}
	return p + m.descriptor.GetName()
}

func (m *Message) SingleName() string {
	spl := strings.Split(m.Name(), ".")
	return spl[len(spl)-1]
}

func (m *Message) StructName(ptr bool) string {
	gopkg := m.GoPackage()
	if gopkg == "main" {
		gopkg = ""
	} else {
		gopkg = filepath.Base(gopkg) + "."
	}
	var sign string
	if ptr {
		sign = "*"
	}
	return sign + gopkg + m.Name()
}

func (m *Message) FullPath() string {
	return m.File.Package() + "." + m.Name()
}
