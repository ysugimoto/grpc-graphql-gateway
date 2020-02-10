package spec

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Enum spec wraps EnumDescriptorProto with keeping file definiton.
type Enum struct {
	descriptor *descriptor.EnumDescriptorProto
	*File

	paths []int
}

func NewEnum(
	e *descriptor.EnumDescriptorProto,
	f *File,
	paths ...int,
) *Enum {
	return &Enum{
		descriptor: e,
		File:       f,
		paths:      paths,
	}
}

func (e *Enum) Comment() string {
	return e.File.getComment(e.paths)
}

func (e *Enum) Name() string {
	return e.descriptor.GetName()
}

func (e *Enum) SingleName() string {
	spl := strings.Split(e.Name(), ".")
	return spl[len(spl)-1]
}

func (e *Enum) Values() []*EnumValue {
	var values []*EnumValue
	for i, v := range e.descriptor.GetValue() {
		paths := make([]int, len(e.paths))
		copy(paths, e.paths)
		values = append(values, NewEnumValue(v, e.File, append(paths, 2, i)...))
	}
	return values
}

func (e *Enum) IsSamePackage(rootPackage string) bool {
	return e.GoPackage() == rootPackage
}
