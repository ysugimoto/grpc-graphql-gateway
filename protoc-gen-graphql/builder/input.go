package builder

import (
	"fmt"
	"strings"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Input builder generates input object between GraphQL and Go program.
// On GraphQL, generates input {...} definition, and on Go,
// genearates InputObject in order to use on some Muration query.
type Input struct {
	m *spec.Message
}

func NewInput(m *spec.Message) *Input {
	return &Input{
		m: m,
	}
}

func (b *Input) BuildQuery() (string, error) {
	lines := []string{}
	if c := b.m.Comment(spec.GraphqlComment); c != "" {
		lines = append(lines, c)
	}
	lines = append(lines, fmt.Sprintf("input %s {", b.m.Name()))

	for _, f := range b.m.Fields() {
		fieldType := f.GraphqlType()
		if f.IsRepeated() {
			fieldType = "[" + fieldType + "]"
		}
		if !f.IsOptional() {
			fieldType += "!"
		}

		if c := f.Comment(spec.GraphqlComment); c != "" {
			lines = append(lines, c)
		}
		lines = append(lines, fmt.Sprintf(
			"  %s: %s",
			f.Name(),
			fieldType,
		))
	}

	lines = append(lines, fmt.Sprintf(
		"} # message %s in %s\n",
		b.m.SingleName(),
		b.m.Filename(),
	))
	return strings.Join(lines, "\n"), nil
}

func (b *Input) BuildProgram() (string, error) {
	var fields []string

	for _, f := range b.m.Fields() {
		fieldType := f.GraphqlGoType()
		if f.IsRepeated() {
			fieldType = "graphql.NewList(" + fieldType + ")"
		}
		if !f.IsOptional() {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}

		fields = append(fields, strings.TrimSpace(fmt.Sprintf(`
			%s
			"%s": &graphql.InputObjectFieldConfig{
				Type: %s,
			},`,
			f.Comment(spec.GoComment),
			f.Name(),
			fieldType,
		)))
	}

	return strings.TrimSpace(fmt.Sprintf(`
%s
var %s = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "%s",
	Fields: graphql.InputObjectConfigFieldMap{
		%s
	},
}) // message %s in %s`,
		b.m.Comment(spec.GoComment),
		spec.PrefixInput(b.m.Name()),
		b.m.Name(),
		strings.TrimSpace(strings.Join(fields, "\n")),
		b.m.SingleName(),
		b.m.Filename(),
	)) + "\n", nil
}
