package builder

import (
	"errors"
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Query builder generates query definition.
type Query struct {
	find     func(name string) *spec.Message
	findEnum func(name string) *spec.Enum
	ms       []*spec.Method
}

func NewQuery(
	f func(name string) *spec.Message,
	fe func(name string) *spec.Enum,
	ms []*spec.Method,
) *Query {
	return &Query{
		find:     f,
		findEnum: fe,
		ms:       ms,
	}
}

func (q *Query) BuildQuery() (string, error) {
	if len(q.ms) == 0 {
		return "", nil
	}

	lines := []string{`type Query {`}
	for _, method := range q.ms {
		i := method.Input()
		input := q.find(i)
		if input == nil {
			return "", errors.New("input " + i + " is not defined in " + method.Package())
		}
		o := method.Output()
		output := q.find(o)
		if output == nil {
			return "", errors.New("output " + o + " is not defined in " + method.Package())
		}

		var fieldName string
		if method.ExposeQuery() != "" {
			field := method.ExposeQueryFields(output)[0]
			fieldName = field.Name()
			if !field.IsOptional() {
				fieldName += "!"
			}
			if field.Label() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				fieldName = "[" + fieldName + "]"
				if resp := method.QueryResponse(); resp != nil {
					if !resp.GetOptional() {
						fieldName += "!"
					}
				}
			}
		} else {
			fieldName = output.Name()
			if resp := method.QueryResponse(); resp != nil {
				if resp.GetRepeated() {
					fieldName = "[" + fieldName + "]"
				}
				if !resp.GetOptional() {
					fieldName += "!"
				}
			}
		}

		if c := method.Comment(spec.GraphqlComment); c != "" {
			lines = append(lines, c)
		}

		lines = append(lines, fmt.Sprintf(
			"  %s(%s): %s",
			method.QueryName(),
			q.ExtractArguments(input),
			fieldName,
		))
	}
	return strings.Join(append(lines, "}\n"), "\n"), nil
}

func (q *Query) ExtractArguments(input *spec.Message) string {
	var args []string

	for _, f := range input.Fields() {
		sign := ""
		if !f.IsOptional() {
			sign = "!"
		}
		args = append(args, fmt.Sprintf(
			"%s: %s%s",
			f.Name(),
			f.GraphqlType(),
			sign,
		))
	}
	return strings.Join(args, ", ")
}

func (q *Query) BuildProgram() (string, error) {
	if len(q.ms) == 0 {
		return "", nil
	}
	fields := make([]string, len(q.ms))
	connections := make(map[string]string)

	for i, method := range q.ms {
		input := q.find(method.Input())
		if input == nil {
			return "", errors.New("failed to resolve input message: " + method.Input())
		}
		args, err := NewArgument(input).BuildProgram()
		if err != nil {
			return "", errors.New("failed to build program for argument message: " + method.Input())
		}
		if args != "" {
			args = strings.TrimSpace(fmt.Sprintf(`
				Args: graphql.FieldConfigArgument{
					%s
				},`,
				args,
			))
		}

		output := q.find(method.Output())
		if output == nil {
			return "", errors.New("failed to resolve output message: " + method.Output())
		}

		serviceName := method.ServiceName()
		if _, ok := connections[serviceName]; !ok {
			c, _ := NewConnection(method.Service).BuildProgram()
			connections[serviceName] = c
		}

		var typeName string
		if method.ExposeQuery() != "" {
			field := method.ExposeQueryFields(output)[0]
			typeName = field.GraphqlGoType()
			if field.Label() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				typeName = "graphql.NewList(" + typeName + ")"
			}
			if !field.IsOptional() {
				typeName = "graphql.NewNonNull(" + typeName + ")"
			}
		} else {
			typeName = spec.PrefixType(output.Name())
			if resp := method.QueryResponse(); resp != nil {
				if resp.GetRepeated() {
					typeName = "graphql.NewList(" + typeName + ")"
				}
				if !resp.GetOptional() {
					typeName = "graphql.NewNonNull(" + typeName + ")"
				}
			}
		}

		resolve, err := NewQueryResolver(q.find, q.findEnum, method).BuildProgram()
		if err != nil {
			return "", errors.New("failed to build resolver function: " + err.Error())
		}

		fields[i] = strings.TrimSpace(fmt.Sprintf(`
			%s
			"%s": &graphql.Field{
				Type: %s,
				%s
				Resolve: %s,
			},`,
			method.Comment(spec.GoComment),
			method.QueryName(),
			typeName,
			args,
			resolve,
		))
	}

	cns := make([]string, len(connections))
	for _, v := range connections {
		cns = append(cns, v)
	}

	return fmt.Sprintf(`
%s

// getQueryFields returns query target fields.
func getQueryFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		%s
	}
}`,
		strings.Join(cns, "\n"),
		strings.Join(fields, "\n"),
	), nil
}
