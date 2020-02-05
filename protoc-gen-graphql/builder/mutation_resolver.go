package builder

import (
	"errors"
	"fmt"
	"strings"

	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// MutationResolver builder generates Resolve section for Go.
// this builder generates Go code which open gRPC connection and close,
// collect resolve arguments, send request to gRPC service with appropriate request type,
// and return extact response type.
type MutationResolver struct {
	find     func(name string) *spec.Message
	findEnum func(name string) *spec.Enum
	m        *spec.Method
}

func NewMutationResolver(
	f func(name string) *spec.Message,
	fe func(name string) *spec.Enum,
	m *spec.Method,
) *MutationResolver {
	return &MutationResolver{
		find:     f,
		findEnum: fe,
		m:        m,
	}
}

func (b *MutationResolver) BuildSchema() (string, error) {
	return "", nil
}

func (b *MutationResolver) BuildProgram() (string, error) {
	input := b.find(b.m.Input())
	if input == nil {
		return "", errors.New("input " + b.m.Input() + " is not defined")
	}
	var pkgName string
	if b.m.GoPackage() != input.GoPackage() {
		if p := filepath.Base(input.GoPackage()); p != "main" {
			pkgName += "."
		}
	}

	output := b.find(b.m.Output())
	if output == nil {
		return "", errors.New("input " + b.m.Output() + " is not defined")
	}
	var response string
	if b.m.ExposeMutation() != "" {
		field := b.m.ExposeMutationFields(output)[0]
		response = ".Get" + strcase.ToCamel(field.Name()) + "()"
	}

	typeName := input.Name()
	if b.m.GoPackage() != input.GoPackage() {
		typeName = input.StructName(false)
	}

	request := []string{
		fmt.Sprintf("req := &%s{}", typeName),
	}

	for _, f := range input.Fields() {
		if !f.IsOptional() {
			if param, err := b.buildRequiredParams(f); err != nil {
				return "", err
			} else {
				request = append(request, param)
			}
		} else {
			if param, err := b.buildOptionalParams(f); err != nil {
				return "", err
			} else {
				request = append(request, param)
			}
		}
	}

	// TODO: implement
	return fmt.Sprintf(`func(p graphql.ResolveParams) (interface{}, error) {
		if c == nil {
			var err error
			if c, err = create%sConnection(); err != nil {
				return nil, err
			}
			defer func() {
				c.Close()
				c = nil
			}()
		}
		client := %sNew%sClient(c)
		%s
		resp,err := client.%s(p.Context, req)
		if err != nil {
			return nil, err
		}
		return resp%s, nil }`,
		b.m.ServiceName(),
		pkgName,
		b.m.Service.Name(),
		strings.Join(request, "\n"),
		b.m.Name(),
		response,
	), nil
}

func (b *MutationResolver) buildRequiredParams(f *spec.Field) (string, error) {
	t := f.GoType()
	name := f.Name()
	camelName := strcase.ToCamel(name)

	switch t {
	case "int64":
		// In graphql, argument comes with int, not in64
		// so we need to cast from int to int64 internaly
		return fmt.Sprintf(`req.%s = int64(p.Args["%s"].(int))`, camelName, name), nil
	case "message":
		m := b.find(strings.TrimPrefix(f.TypeName(), "."))
		if m == nil {
			return "", errors.New("failed to find struct type for field: " + name)
		}
		sName := m.StructName(true)
		if f.GoPackage() == m.GoPackage() {
			sName = "*" + m.SingleName()
		}
		return fmt.Sprintf(`req.%s = p.Args["%s"].(%s)`, camelName, name, sName), nil
	case "enum":
		e := b.findEnum(strings.TrimPrefix(f.TypeName(), "."))
		if e == nil {
			return "", errors.New("failed to find enum type for field: " + name)
		}
		eName := e.Name()
		if f.GoPackage() == e.GoPackage() {
			eName = e.SingleName()
		}
		return fmt.Sprintf(`req.%s = p.Args["%s"].(%s)`, camelName, name, eName), nil
	default:
		return fmt.Sprintf(`req.%s = p.Args["%s"].(%s)`, strcase.ToCamel(name), name, t), nil
	}
}

func (b *MutationResolver) buildOptionalParams(f *spec.Field) (string, error) {
	t := f.GoType()
	name := f.Name()
	camelName := strcase.ToCamel(name)

	var assignLine, castType string
	switch t {
	case "int64":
		assignLine = fmt.Sprintf(`req.%s = int64(arg)`, camelName)
		castType = "int"
	case "message":
		m := b.find(strings.TrimPrefix(f.TypeName(), "."))
		if m == nil {
			return "", errors.New("failed to find struct type for field: " + name)
		}
		sName := m.StructName(true)
		if f.GoPackage() == m.GoPackage() {
			sName = "*" + m.SingleName()
		}
		assignLine = fmt.Sprintf(`req.%s = arg`, camelName)
		castType = sName
	case "enum":
		e := b.findEnum(strings.TrimPrefix(f.TypeName(), "."))
		if e == nil {
			return "", errors.New("failed to find enum type for field: " + name)
		}
		eName := e.Name()
		if f.GoPackage() == e.GoPackage() {
			eName = e.SingleName()
		}
		assignLine = fmt.Sprintf(`req.%s = arg`, camelName)
		castType = eName
	default:
		assignLine = fmt.Sprintf(`req.%s = arg`, camelName)
		castType = t
	}

	return strings.TrimSpace(fmt.Sprintf(`
		if v, ok := p.Args["%s"]; !ok {
			return nil, errors.New("%s is not found in parameter")
		} else if arg, ok := v.(%s); !ok {
			return nil, errors.New("failed to do type conversion to int for field %s")
		} else {
			req.%s = int64(arg)
		}`,
		name,
		name,
		castType,
		name,
		assignLine,
	)), nil
}
