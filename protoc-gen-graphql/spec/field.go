package spec

import (
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
	"strings"
)

// Field spec wraps FieldDescriptorProto with keeping file info
type Field struct {
	descriptor *descriptor.FieldDescriptorProto
	Option     *graphql.GraphqlField
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
			if field, ok := ext.(*graphql.GraphqlField); !ok {
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

func (f *Field) Comment(t CommentType) string {
	return f.File.getComment(f.paths, t)
}

func (f *Field) Name() string {
	return f.descriptor.GetName()
}

func (f *Field) Type() descriptor.FieldDescriptorProto_Type {
	return f.descriptor.GetType()
}

func (f *Field) TypeName() string {
	return f.descriptor.GetTypeName()
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
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		spec := strings.Split(f.TypeName(), ".")
		name := spec[len(spec)-1]
		return name
	default:
		return "Unknown"
	}
}

// GraphqlGoType returns appropriate graphql-go type
func (f *Field) GraphqlGoType() string {
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
		spec := strings.Split(f.TypeName(), ".")
		name := spec[len(spec)-1]
		return PrefixType(name)
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		spec := strings.Split(f.TypeName(), ".")
		name := spec[len(spec)-1]
		return PrefixEnum(name)
	default:
		return "interface{}"
	}
}

// GoType returns appropriate go type
// but message and enum is special value -- it should be treated in each builder
func (f *Field) GoType() string {
	switch f.Type() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float64"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return "message"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "enum"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	default:
		return "interface{}"
	}
}
