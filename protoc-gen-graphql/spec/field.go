package spec

import (
	"strings"

	"path/filepath"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Field spec wraps FieldDescriptorProto with keeping file info
type Field struct {
	descriptor  *descriptor.FieldDescriptorProto
	Option      *graphql.GraphqlField
	TypeMessage *Message
	TypeEnum    *Enum
	*File

	paths []int
}

func NewField(
	d *descriptor.FieldDescriptorProto,
	f *File,
	paths ...int,
) *Field {
	var o *graphql.GraphqlField
	if opts := d.GetOptions(); opts != nil {
		if ext, err := proto.GetExtension(opts, graphql.E_Field); err == nil {
			if field, ok := ext.(*graphql.GraphqlField); ok {
				o = field
			}
		}
	}

	return &Field{
		descriptor: d,
		Option:     o,
		File:       f,
		paths:      paths,
	}
}

func (f *Field) Comment() string {
	if strings.HasPrefix(f.Package(), "google.protobuf") {
		return ""
	}
	return f.File.getComment(f.paths)
}

func (f *Field) Name() string {
	return f.descriptor.GetName()
}

func (f *Field) Type() descriptor.FieldDescriptorProto_Type {
	return f.descriptor.GetType()
}

func (f *Field) TypeName() string {
	return strings.TrimPrefix(f.descriptor.GetTypeName(), ".")
}

func (f *Field) Label() descriptor.FieldDescriptorProto_Label {
	return f.descriptor.GetLabel()
}

func (f *Field) IsOptional() bool {
	if f.Option == nil {
		return false
	}
	return f.Option.GetOptional()
}

func (f *Field) IsRepeated() bool {
	return f.Label() == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

func (f *Field) FieldType(rootPackage string) string {
	fieldType := f.GraphqlGoType(filepath.Base(rootPackage), false)
	if f.IsRepeated() {
		fieldType = "graphql.NewList(" + fieldType + ")"
	}
	if !f.IsOptional() {
		fieldType = "graphql.NewNonNull(" + fieldType + ")"
	}
	return fieldType
}

func (f *Field) FieldTypeInput(rootPackage string) string {
	fieldType := f.GraphqlGoType(filepath.Base(rootPackage), true)
	if f.IsRepeated() {
		fieldType = "graphql.NewList(" + fieldType + ")"
	}
	if !f.IsOptional() {
		fieldType = "graphql.NewNonNull(" + fieldType + ")"
	}
	return fieldType
}

func (f *Field) SchemaType() string {
	fieldType := f.GraphqlType()
	if f.IsRepeated() {
		fieldType = "[" + fieldType + "]"
	}
	if !f.IsOptional() {
		fieldType += "!"
	}
	return fieldType
}

func (f *Field) DefaultValue() string {
	if f.Option == nil {
		return ""
	}
	switch f.Type() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		return f.Option.GetDefault()
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return `"` + f.Option.GetDefault() + `"`
	default:
		return ""
	}
}

// GraphqlType returns appropriate GraphQL type
func (f *Field) GraphqlType() string {
	switch f.Type() {
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
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		tn := strings.TrimPrefix(f.TypeName(), f.TypeMessage.Package()+".")
		return strings.ReplaceAll(tn, ".", "_")
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		tn := strings.TrimPrefix(f.TypeName(), f.TypeEnum.Package()+".")
		return strings.ReplaceAll(tn, ".", "_")
	default:
		return "Unknown"
	}
}

// GraphqlGoType returns appropriate graphql-go type
func (f *Field) GraphqlGoType(rootPackage string, isInput bool) string {
	switch f.Type() {
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
		var pkgPrefix string
		tn := strings.TrimPrefix(f.TypeName(), f.TypeMessage.Package()+".")
		if filepath.Base(f.TypeMessage.GoPackage()) != rootPackage {
			if !strings.HasPrefix(f.TypeMessage.Package(), "google.protobuf") {
				pkgPrefix = filepath.Base(f.TypeMessage.GoPackage()) + "."
			}
		}
		if isInput {
			return pkgPrefix + PrefixInput(strings.ReplaceAll(tn, ".", "_"))
		} else {
			return pkgPrefix + PrefixType(strings.ReplaceAll(tn, ".", "_"))
		}
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		var pkgPrefix string
		tn := strings.TrimPrefix(f.TypeName(), f.TypeEnum.Package()+".")
		if filepath.Base(f.TypeEnum.GoPackage()) != rootPackage {
			if !strings.HasPrefix(f.TypeEnum.Package(), "google.protobuf") {
				pkgPrefix = filepath.Base(f.TypeEnum.GoPackage()) + "."
			}
		}
		return pkgPrefix + PrefixEnum(strings.ReplaceAll(tn, ".", "_"))
	default:
		return "interface{}"
	}
}
