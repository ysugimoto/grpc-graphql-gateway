package types

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type ArgumentSpec struct {
	Descriptor *descriptor.FieldDescriptorProto
}
