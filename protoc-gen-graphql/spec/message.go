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
	fields []*Field

	*Dependencies
	PluckFields []*Field
}

func NewMessage(
	d *descriptor.DescriptorProto,
	f *File,
	prefix []string,
	paths ...int,
) *Message {
	m := &Message{
		descriptor: d,
		File:       f,
		prefix:     prefix,
		paths:      paths,
		fields:     make([]*Field, 0),

		Dependencies: NewDependencies(),
	}
	for i, field := range d.GetField() {
		ps := make([]int, len(paths))
		copy(ps, paths)
		m.fields = append(m.fields, NewField(field, f, append(ps, 2, i)...))
	}
	return m
}

func (m *Message) Fields() []*Field {
	return m.fields
}

func (m *Message) TypeFields() []*Field {
	if m.PluckFields == nil {
		return m.Fields()
	}
	return m.PluckFields
}

func (m *Message) Comment() string {
	if strings.HasPrefix(m.Package(), "google.protobuf") {
		return ""
	}
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

func (m *Message) Interfaces() (ifs []*Message) {
	for _, f := range m.fields {
		if !f.IsCyclic {
			continue
		}
		ifs = append(ifs, f.DependType.(*Message))
	}
	return ifs
}
