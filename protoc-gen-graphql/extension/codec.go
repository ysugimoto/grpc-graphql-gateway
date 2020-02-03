package extension

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func ConvertGraphqlType(f *descriptor.FieldDescriptorProto) string {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "Boolean"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "Float"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "Int"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "String"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		spec := strings.Split(f.GetTypeName(), ".")
		name := spec[len(spec)-1]
		return name
	default:
		return "Unknown"
	}
}

func ConvertGoType(f *descriptor.FieldDescriptorProto) string {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "graphql.Bool"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "graphql.Float"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "graphql.Int"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "graphql.String"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		spec := strings.Split(f.GetTypeName(), ".")
		name := spec[len(spec)-1]
		return MessageName(name)
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		spec := strings.Split(f.GetTypeName(), ".")
		name := spec[len(spec)-1]
		return EnumName(name)
	default:
		return "interface{}"
	}
}

func ConvertGoPrimitiveType(f *descriptor.FieldDescriptorProto) string {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	default:
		return "interface{}"
	}
}
