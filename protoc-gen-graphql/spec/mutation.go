package spec

import (
	"path/filepath"

	"github.com/iancoleman/strcase"
)

// Mutation spec wraps MethodDescriptorProto.
type Mutation struct {
	*Method
	Input  *Message
	Output *Message
}

func NewMutation(m *Method, input, output *Message) *Mutation {
	return &Mutation{
		Method: m,
		Input:  input,
		Output: output,
	}
}

func (m *Mutation) MutationType() string {
	if m.Method.ExposeMutation() != "" {
		field := m.Method.ExposeMutationFields(m.Output)[0]
		fieldName := field.GraphqlGoType()
		if field.IsRepeated() {
			fieldName = "graphql.NewList(" + fieldName + ")"
		}
		if field.IsOptional() {
			fieldName = "graphql.NewNonNull(" + fieldName + ")"
		}
		return fieldName
	}

	typeName := PrefixType(m.Output.Name())
	if resp := m.Method.MutationResponse(); resp != nil {
		if resp.GetRepeated() {
			typeName = "graphql.NewList(" + typeName + ")"
		}
		if !resp.GetOptional() {
			typeName = "graphql.NewNonNull(" + typeName + ")"
		}
	}
	return typeName
}

func (m *Mutation) OutputName() string {
	if m.Method.ExposeMutation() != "" {
		field := m.Method.ExposeMutationFields(m.Output)[0]
		fieldType := field.GraphqlType()
		if field.IsRepeated() {
			fieldType = "[" + fieldType + "]"
		}
		if !field.IsOptional() {
			fieldType += "!"
		}
		return fieldType
	}

	typeName := m.Output.Name()
	if resp := m.Method.MutationResponse(); resp != nil {
		if resp.GetRepeated() {
			typeName = "[" + typeName + "]"
		}
		if !resp.GetOptional() {
			typeName += "!"
		}
	} else {
		typeName += "!"
	}
	return typeName
}

func (m *Mutation) InputName() string {
	inputName := m.Input.SingleName()
	if req := m.Method.MutationRequest(); req != nil {
		if n := req.GetName(); n != "" {
			inputName = n
		}
	}
	return inputName
}

func (m *Mutation) RequestType() string {
	if m.Method.GoPackage() != m.Input.GoPackage() {
		return m.Input.StructName(false)
	}
	return m.Input.Name()
}

func (m *Mutation) Package() string {
	var pkgName string
	if m.Method.GoPackage() != m.Input.GoPackage() {
		pkgName = filepath.Base(m.Input.GoPackage())
		if pkgName != "main" {
			pkgName += "."
		}
	}
	return pkgName
}

func (m *Mutation) Expose() string {
	if m.Method.ExposeMutation() == "" {
		return ""
	}
	field := m.Method.ExposeMutationFields(m.Output)[0]
	return strcase.ToCamel(field.Name())
}
