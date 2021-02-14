package spec

import (
	"fmt"
	"log"
	"strings"

	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Query spec wraps MethodDescriptorProto.
type Query struct {
	*Method
	Input  *Message
	Output *Message

	isCamel bool
}

func NewQuery(m *Method, input, output *Message, isCamel bool) *Query {
	return &Query{
		Method:  m,
		Input:   input,
		Output:  output,
		isCamel: isCamel,
	}
}

func (q *Query) IsResolver() bool {
	return q.Schema.GetType() == graphql.GraphqlType_RESOLVER
}

func (q *Query) IsCamel() bool {
	return q.isCamel
}

func (q *Query) QueryName() string {
	return q.Schema.GetName()
}

func (q *Query) Request() *graphql.GraphqlRequest {
	return q.Schema.GetRequest()
}

func (q *Query) Response() *graphql.GraphqlResponse {
	return q.Schema.GetResponse()
}

func (q *Query) IsPluckRequest() bool {
	req := q.Request()
	if req == nil {
		return false
	}
	return len(req.GetPlucks()) > 0
}

func (q *Query) IsPluckResponse() bool {
	resp := q.Response()
	if resp == nil {
		return false
	}
	return resp.GetPluck() != ""
}

func (q *Query) PluckRequest() []*Field {
	var plucks []string
	if req := q.Request(); req != nil {
		plucks = req.GetPlucks()
	}

	if len(plucks) == 0 {
		return q.Input.Fields()
	}
	var fields []*Field
	for _, f := range q.Input.Fields() {
		for _, p := range plucks {
			if p != f.Name() {
				continue
			}
			fields = append(fields, f)
		}
	}
	return fields
}

func (q *Query) PluckResponse() []*Field {
	var pluck string
	if resp := q.Response(); resp != nil {
		pluck = resp.GetPluck()
	}

	if pluck == "" {
		return q.Output.Fields()
	}
	var fields []*Field
	for _, f := range q.Output.Fields() {
		if f.Name() != pluck {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (q *Query) QueryType() string {
	if q.IsPluckResponse() {
		field := q.PluckResponse()[0]
		return field.FieldType(q.GoPackage())
	}

	var pkgPrefix string
	if q.GoPackage() != q.Output.GoPackage() {
		if IsGooglePackage(q.Output) {
			ptypeName, err := getImplementedPtypes(q.Output)
			if err != nil {
				log.Fatalf("[PROTOC-GEN-GRAPHQL] Error: %s\n", err)
			}
			pkgPrefix = "gql_ptypes_" + ptypeName + "."
		} else {
			pkgPrefix = filepath.Base(q.GoPackage())
			if index := strings.Index(pkgPrefix, ";"); index > -1 {
				pkgPrefix = pkgPrefix[index+1:]
			}
			if pkgPrefix != mainPackage {
				pkgPrefix += "."
			}
		}
	}

	typeName := pkgPrefix + PrefixType(q.Output.Name())
	if resp := q.Response(); resp != nil {
		if resp.GetRequired() {
			typeName = "graphql.NewNonNull(" + typeName + ")"
		}
	}
	return typeName
}

func (q *Query) Args() []*Field {
	return q.PluckRequest()
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
	if fields := q.PluckResponse(); len(fields) > 0 {
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

	typeName := q.Output.Name()
	if resp := q.Response(); resp != nil {
		if resp.GetRequired() {
			typeName += "!"
		}
	}
	return typeName
}

func (q *Query) InputType() string {
	if q.Method.GoPackage() != q.Input.GoPackage() {
		if IsGooglePackage(q.Input) {
			ptypeName, err := getImplementedPtypes(q.Input)
			if err != nil {
				log.Fatalln("[PROTOC-GEN-GRAPHQL] Error:", err)
			}
			return "gql_ptypes_" + ptypeName + "." + q.Input.Name()
		}
		return q.Input.StructName(false)
	}
	return q.Input.Name()
}

func (q *Query) PluckResponseFieldName() string {
	fields := q.PluckResponse()
	return strcase.ToCamel(fields[0].Name())
}

func (q *Query) Package() string {
	var pkgName string
	if q.GoPackage() != q.Input.GoPackage() {
		pkgName = filepath.Base(q.Input.GoPackage())
		if pkgName != mainPackage {
			pkgName += "."
		}
	}
	return pkgName
}
