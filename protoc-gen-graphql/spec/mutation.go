package spec

import (
	"log"

	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Mutation spec wraps MethodDescriptorProto.
type Mutation struct {
	*Method
	Input  *Message
	Output *Message

	isCamel bool
}

func NewMutation(m *Method, input, output *Message, isCamel bool) *Mutation {
	return &Mutation{
		Method:  m,
		Input:   input,
		Output:  output,
		isCamel: isCamel,
	}
}

func (m *Mutation) IsCamel() bool {
	return m.isCamel
}

func (m *Mutation) MutationName() string {
	return m.Schema.GetName()
}

func (m *Mutation) Request() *graphql.GraphqlRequest {
	return m.Schema.GetRequest()
}

func (m *Mutation) IsPluckRequest() bool {
	req := m.Request()
	if req == nil {
		return false
	}
	return len(req.GetPlucks()) > 0
}

func (m *Mutation) Response() *graphql.GraphqlResponse {
	return m.Schema.GetResponse()
}

func (m *Mutation) IsPluckResponse() bool {
	resp := m.Response()
	if resp == nil {
		return false
	}
	return resp.GetPluck() != ""
}

func (m *Mutation) InputName() string {
	if req := m.Request(); req != nil {
		return req.GetName()
	}
	return ""
}

func (m *Mutation) PluckRequest() []*Field {
	var plucks []string
	if req := m.Request(); req != nil {
		plucks = req.GetPlucks()
	}

	if len(plucks) == 0 {
		return m.Input.Fields()
	}

	var fields []*Field
	for _, f := range m.Input.Fields() {
		for _, p := range plucks {
			if p != f.Name() {
				continue
			}
			fields = append(fields, f)
		}
	}
	return fields
}

func (m *Mutation) PluckResponse() []*Field {
	var pluck string
	if resp := m.Response(); resp != nil {
		pluck = resp.GetPluck()
	}

	if pluck == "" {
		return m.Output.Fields()
	}

	var fields []*Field
	for _, f := range m.Output.Fields() {
		if pluck != f.Name() {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (m *Mutation) Args() []*Field {
	return m.PluckRequest()
}

func (m *Mutation) MutationType() string {
	var pkgPrefix string
	if m.GoPackage() != m.Output.GoPackage() {
		if IsGooglePackage(m.Output) {
			ptypeName, err := getImplementedPtypes(m.Output)
			if err != nil {
				log.Fatalln("[PROTOC-GEN-GRAPHQL] Error:", err)
			}
			pkgPrefix = "gql_ptypes_" + ptypeName + "."
		} else {
			pkgPrefix = m.Output.GoPackage()
			if pkgPrefix != mainPackage {
				pkgPrefix += "."
			}
		}
	}
	typeName := pkgPrefix + PrefixType(m.Output.Name())
	if resp := m.Response(); resp != nil {
		if resp.GetRequired() {
			typeName = "graphql.NewNonNull(" + typeName + ")"
		}
	}
	return typeName
}

func (m *Mutation) OutputName() string {
	if fields := m.PluckResponse(); len(fields) > 0 {
		field := fields[0]
		fieldType := field.GraphqlType()
		if field.IsRepeated() {
			fieldType = "[" + fieldType + "]"
		}
		if field.IsRequired() {
			fieldType += "!"
		}
		return fieldType
	}

	typeName := m.Output.Name()
	if resp := m.Response(); resp != nil {
		if resp.GetRequired() {
			typeName += "!"
		}
	}
	return typeName
}

//
func (m *Mutation) InputType() string {
	if m.Method.GoPackage() != m.Input.GoPackage() {
		return m.Input.StructName(false)
	}
	return m.Input.Name()
}

func (m *Mutation) PluckResponseFieldName() string {
	fields := m.PluckResponse()
	return strcase.ToCamel(fields[0].Name())
}

func (m *Mutation) Package() string {
	var pkgName string
	if m.GoPackage() != m.Input.GoPackage() {
		pkgName = filepath.Base(m.Input.GoPackage())
		if pkgName != mainPackage {
			pkgName += "."
		}
	}
	return pkgName
}
