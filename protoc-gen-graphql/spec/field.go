package spec

import (
	"log"
	"strings"

	// nolint: staticcheck
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Field spec wraps FieldDescriptorProto with keeping file info
type Field struct {
	descriptor *descriptor.FieldDescriptorProto
	Option     *graphql.GraphqlField
	*File

	paths []int

	DependType    interface{}
	IsCyclic      bool
	isCamel       bool
	forceRequired bool
}

func NewField(
	d *descriptor.FieldDescriptorProto,
	f *File,
	isCamel bool,
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
		isCamel:    isCamel,
	}
}

func (f *Field) Comment() string {
	if IsGooglePackage(f) {
		return ""
	}
	return f.File.getComment(f.paths)
}

func (f *Field) Name() string {
	return f.descriptor.GetName()
}

func (f *Field) setRequiredField() {
	f.forceRequired = true
}

func (f *Field) FieldName() string {
	if f.isCamel {
		return strcase.ToLowerCamel(f.Name())
	}
	return f.Name()
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

func (f *Field) IsRequired() bool {
	if f.forceRequired {
		return true
	}

	if f.Option == nil {
		return false
	}
	return f.Option.GetRequired()
}

func (f *Field) IsOmit() bool {
	if f.Option == nil {
		return false
	}
	return f.Option.GetOmit()
}

func (f *Field) IsRepeated() bool {
	return f.Label() == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

func (f *Field) FieldType(rootPackage string) string {
	pkg := NewGoPackageFromString(rootPackage)
	fieldType := f.GraphqlGoType(pkg.Name, false)
	if f.IsRequired() {
		fieldType = "graphql.NewNonNull(" + fieldType + ")"
	}
	if f.IsRepeated() {
		fieldType = "graphql.NewList(" + fieldType + ")"
		if f.IsRequired() {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}
	}
	return fieldType
}

func (f *Field) FieldTypeInput(rootPackage string) string {
	pkg := NewGoPackageFromString(rootPackage)
	fieldType := f.GraphqlGoType(pkg.Name, true)
	if f.IsRequired() {
		fieldType = "graphql.NewNonNull(" + fieldType + ")"
	}
	if f.IsRepeated() {
		fieldType = "graphql.NewList(" + fieldType + ")"
		if f.IsRequired() {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}
	}
	return fieldType
}

func (f *Field) SchemaType() string {
	fieldType := f.GraphqlType()
	if f.IsRepeated() {
		fieldType = "[" + fieldType + "]"
	}
	if f.IsRequired() {
		fieldType += "!"
	}
	return fieldType
}

func (f *Field) SchemaInputType() string {
	var prefix string
	if f.Type() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		m := f.DependType.(*Message) // nolint: errcheck
		if f.Package() == m.Package() || IsGooglePackage(f) {
			prefix = "Input_"
		}
	}

	fieldType := prefix + f.GraphqlType()
	if f.IsRepeated() {
		fieldType = "[" + fieldType + "]"
	}
	if f.IsRequired() {
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
		m := f.DependType.(*Message) // nolint: errcheck
		tn := strings.TrimPrefix(f.TypeName(), m.Package()+".")
		return strings.ReplaceAll(tn, ".", "_")
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		e := f.DependType.(*Enum) // nolint: errcheck
		tn := strings.TrimPrefix(f.TypeName(), e.Package()+".")
		return strings.ReplaceAll(tn, ".", "_")
	default:
		return "Unknown"
	}
}

// GraphqlGoType returns appropriate graphql-go type
func (f *Field) GraphqlGoType(rootPackage string, isInput bool) string {
	switch f.Type() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "graphql.Boolean"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "graphql.Float"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "graphql.Int"
	case descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "graphql.String"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		m := f.DependType.(*Message) // nolint: errcheck
		tn := strings.TrimPrefix(f.TypeName(), m.Package()+".")
		if f.IsCyclic {
			return PrefixInterface(strings.ReplaceAll(tn, ".", "_"))
		}
		var pkgPrefix string
		pkg := NewPackage(m)
		if IsGooglePackage(m) {
			// Case google.protobuf.XXX
			ptypeName, err := getImplementedPtypes(m)
			if err != nil {
				log.Fatalln("[PROTOC-GEN-GRAPHQL] Error:", err)
			}
			pkgPrefix = "gql_ptypes_" + ptypeName + "."
		} else if rootPackage != "." {
			// Case message is nested, also includes map_entry
			if pkg.Name != rootPackage {
				if !IsGooglePackage(m) {
					pkgPrefix = pkg.Name + "."
				}
			}
		}
		if isInput {
			return pkgPrefix + PrefixInput(strings.ReplaceAll(tn, ".", "_"))
		}
		return pkgPrefix + PrefixType(strings.ReplaceAll(tn, ".", "_"))
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		e := f.DependType.(*Enum) // nolint: errcheck
		var pkgPrefix string
		pkg := NewPackage(e)
		if pkg.Name != rootPackage {
			if !IsGooglePackage(e) {
				pkgPrefix = pkg.Name + "."
			}
		}
		tn := strings.TrimPrefix(f.TypeName(), e.Package()+".")
		return pkgPrefix + PrefixEnum(strings.ReplaceAll(tn, ".", "_"))
	default:
		return "graphql.SkipDirective"
	}
}

func (f *Field) IsResolve() bool {
	if f.Option == nil {
		return false
	}

	return f.Option.GetResolver() != ""
}

func (f *Field) ResolveSubField(services []*Service) *Query {
	name := f.Option.GetResolver()
	var query *Query
OUTER:
	for _, s := range services {
		for _, q := range s.Queries {
			if !q.IsResolver() {
				continue
			}
			if name == q.QueryName() {
				query = q
				break OUTER
			}
		}
	}

	if query == nil {
		panic("Could not find field resolve " + name + " in defined queries")
	}
	return query
}
