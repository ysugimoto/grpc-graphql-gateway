package builder

import (
	"fmt"
	"strings"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Argument builder is builder for graphql argument definition.
type Argument struct {
	m *spec.Message
}

func NewArgument(m *spec.Message) *Argument {
	return &Argument{
		m: m,
	}
}

// On Graphql schema generation, this builder nothing to geneate
func (b *Argument) BuildQuery() (string, error) {
	return "", nil
}

// Generate Go program for field arguments
func (b *Argument) BuildProgram() (string, error) {
	var args []string

	for _, f := range b.m.Fields() {
		fieldType := f.GraphqlGoType()
		if f.IsRepeated() {
			fieldType = "graphql.NewList(" + fieldType + ")"
		}
		if !f.IsOptional() {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}

		args = append(args, strings.TrimSpace(fmt.Sprintf(`
			%s
			"%s": &graphql.ArgumentConfig{
				Type: %s,
			},`,
			f.Comment(spec.GoComment),
			f.Name(),
			fieldType,
		)))
	}

	return strings.TrimSpace(strings.Join(args, "\n")), nil
}
