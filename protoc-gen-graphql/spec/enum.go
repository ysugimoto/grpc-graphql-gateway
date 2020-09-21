package spec

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Enum spec wraps EnumDescriptorProto with keeping file definition.
type Enum struct {
	descriptor *descriptor.EnumDescriptorProto
	*File

	prefix []string
	paths  []int
	values []*EnumValue

	*Dependencies
}

func NewEnum(
	d *descriptor.EnumDescriptorProto,
	f *File,
	prefix []string,
	paths ...int,
) *Enum {

	e := &Enum{
		descriptor: d,
		File:       f,
		prefix:     prefix,
		paths:      paths,
		values:     make([]*EnumValue, 0),

		Dependencies: NewDependencies(),
	}
	for i, v := range d.GetValue() {
		ps := make([]int, len(paths))
		copy(ps, paths)
		e.values = append(e.values, NewEnumValue(v, f, append(ps, 2, i)...))
	}
	return e
}

func (e *Enum) Comment() string {
	return e.File.getComment(e.paths)
}

// Get name for definition generation
func (e *Enum) Name() string {
	var p string
	if len(e.prefix) > 0 {
		p = strings.Join(e.prefix, "_") + "_"
	}
	return p + e.descriptor.GetName()
}

// Get path name with prefix
func (e *Enum) PathName() string {
	var p string
	if len(e.prefix) > 0 {
		p = strings.Join(e.prefix, ".") + "."
	}
	return p + e.descriptor.GetName()
}

func (e *Enum) SingleName() string {
	spl := strings.Split(e.Name(), ".")
	return spl[len(spl)-1]
}

func (e *Enum) Values() []*EnumValue {
	return e.values
}

func (e *Enum) FullPath() string {
	return e.File.Package() + "." + e.PathName()
}
