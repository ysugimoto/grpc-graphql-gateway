package types

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type Enum struct {
	Descriptor *descriptor.EnumDescriptorProto
	File       *descriptor.FileDescriptorProto
}

func (e *Enum) MessageName() string {
	spec := strings.Split(e.Descriptor.GetName(), ".")
	return spec[len(spec)-1]
}

func (e *Enum) Filename() string {
	return e.File.GetName()
}
