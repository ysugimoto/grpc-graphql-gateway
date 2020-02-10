package spec

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"path/filepath"
)

// Query spec wraps MethodDescriptorProto.
type Query struct {
	*Method
	Input  *Message
	Output *Message
}

func NewQuery(m *Method, input, output *Message) *Query {
	return &Query{
		Method: m,
		Input:  input,
		Output: output,
	}
}

func (q *Query) QueryType() string {
	if q.Method.ExposeQuery() != "" {
		field := q.Method.ExposeQueryFields(q.Output)[0]
		return field.FieldType()
	}

	typeName := PrefixType(q.Output.Name())
	if resp := q.Method.QueryResponse(); resp != nil {
		if resp.GetRepeated() {
			typeName = "graphql.NewList(" + typeName + ")"
		}
		if !resp.GetOptional() {
			typeName = "graphql.NewNonNull(" + typeName + ")"
		}
	}
	return typeName
}

func (q *Query) Args() []*Field {
	return q.Input.Fields()
}

func (q *Query) SchemaArgs() string {
	args := make([]string, len(q.Input.Fields()))
	for i, v := range q.Input.Fields() {
		var defValue string
		if d := v.DefaultValue(); d != "" {
			defValue = " = " + d
		}
		args[i] = fmt.Sprintf("%s: %s%s", v.Name(), v.SchemaType(), defValue)
	}
	return strings.Join(args, ", ")
}

func (q *Query) OutputName() string {
	if q.Method.ExposeQuery() != "" {
		field := q.Method.ExposeQueryFields(q.Output)[0]
		fieldType := field.GraphqlType()
		if field.IsRepeated() {
			fieldType = "[" + fieldType + "]"
		}
		if !field.IsOptional() {
			fieldType += "!"
		}
		return fieldType
	}

	typeName := q.Output.Name()
	if resp := q.Method.QueryResponse(); resp != nil {
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

func (q *Query) RequestType() string {
	if q.Method.GoPackage() != q.Input.GoPackage() {
		return q.Input.StructName(false)
	}
	return q.Input.Name()
}

func (q *Query) Package() string {
	var pkgName string
	if q.Method.GoPackage() != q.Input.GoPackage() {
		pkgName = filepath.Base(q.Input.GoPackage())
		if pkgName != "main" {
			pkgName += "."
		}
	}
	return pkgName
}

func (q *Query) Expose() string {
	if q.Method.ExposeQuery() == "" {
		return ""
	}
	field := q.Method.ExposeQueryFields(q.Output)[0]
	return strcase.ToCamel(field.Name())
}
