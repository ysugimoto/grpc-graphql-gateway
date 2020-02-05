package builder

import (
	"errors"
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Muration builder generates muration query definition.
type Mutation struct {
	find     func(name string) *spec.Message
	findEnum func(name string) *spec.Enum
	ms       []*spec.Method
}

func NewMutation(
	f func(name string) *spec.Message,
	fe func(name string) *spec.Enum,
	ms []*spec.Method,
) *Mutation {
	return &Mutation{
		find:     f,
		findEnum: fe,
		ms:       ms,
	}
}

func (m *Mutation) BuildQuery() (string, error) {
	if len(m.ms) == 0 {
		return "", nil
	}
	inputTypes := make([]string, 0)

	lines := []string{`type Mutation {`}
	for _, method := range m.ms {
		i := method.Input()
		input := m.find(i)
		if input == nil {
			return "", errors.New("input " + i + " is not defined in " + method.Package())
		}
		o := method.Output()
		output := m.find(o)
		if output == nil {
			return "", errors.New("output " + o + " is not defined in " + method.Package())
		}
		inputType, err := NewInput(input).BuildQuery()
		if err != nil {
			return "", errors.New("failed to build program for input type: " + method.Input())
		}
		inputTypes = append(inputTypes, inputType)
		argName := input.SingleName()
		if req := method.MutationRequest(); req != nil {
			if n := req.GetName(); n != "" {
				argName = n
			}
		}
		argType := input.Name()

		var fieldName string
		if method.ExposeMutation() != "" {
			field := method.ExposeMutationFields(output)[0]
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
			if resp := method.MutationResponse(); resp != nil {
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
			"  %s(%s: %s): %s",
			method.MutationName(),
			argName,
			argType,
			fieldName,
		))
	}
	lines = append(lines, "}\n", strings.Join(inputTypes, "\n"))
	return strings.Join(lines, "\n"), nil
}

func (m *Mutation) BuildProgram() (string, error) {
	if len(m.ms) == 0 {
		return "", nil
	}

	fields := make([]string, len(m.ms))
	inputTypes := make([]string, len(m.ms))

	for i, method := range m.ms {
		input := m.find(method.Input())
		if input == nil {
			return "", errors.New("failed to resolve input message: " + method.Input())
		}
		inputType, err := NewInput(input).BuildProgram()
		if err != nil {
			return "", errors.New("failed to build program for input type: " + method.Input())
		}
		inputTypes[i] = inputType
		argTypeName := spec.PrefixInput(input.Name())
		argName := input.SingleName()
		if req := method.MutationRequest(); req != nil {
			if n := req.GetName(); n != "" {
				argName = n
			}
		}

		output := m.find(method.Output())
		if output == nil {
			return "", errors.New("failed to resolve output message: " + method.Output())
		}

		var typeName string
		if method.ExposeMutation() != "" {
			field := method.ExposeMutationFields(output)[0]
			typeName = field.GraphqlGoType()
			if field.Label() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				typeName = "graphql.NewList(" + typeName + ")"
			}
			if !field.IsOptional() {
				typeName = "graphql.NewNonNull(" + typeName + ")"
			}
		} else {
			typeName = spec.PrefixType(output.Name())
			if resp := method.MutationResponse(); resp != nil {
				if resp.GetRepeated() {
					typeName = "graphql.NewList(" + typeName + ")"
				}
				if !resp.GetOptional() {
					typeName = "graphql.NewNonNull(" + typeName + ")"
				}
			}
		}

		resolve, err := NewMutationResolver(m.find, m.findEnum, method).BuildProgram()
		if err != nil {
			return "", errors.New("failed to build resolve function: " + err.Error())
		}
		fields[i] = strings.TrimSpace(fmt.Sprintf(`
			%s
			"%s": &graphql.Field{
				Type: %s,
				Args: graphql.FieldConfigArgument{
				  "%s": &graphql.ArgumentConfig{
					  Type: %s,
				  },
				},
				Resolve: %s,
			},`,
			method.Comment(spec.GoComment),
			method.QueryName(),
			typeName,
			argName,
			argTypeName,
			resolve,
		))
	}

	return fmt.Sprintf(`
%s

// getMutationFields returns mutation target fields.
func getMutationFields(c *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		%s
	}
}`,
		strings.Join(inputTypes, "\n"),
		strings.Join(fields, "\n"),
	), nil
}
