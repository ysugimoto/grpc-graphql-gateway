package builder

import (
	"fmt"
	"strings"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Enum builder generates enum definition between GraphQL and Go pgoram.
type Enum struct {
	e *spec.Enum
}

func NewEnum(e *spec.Enum) *Enum {
	return &Enum{
		e: e,
	}
}

func (b *Enum) BuildQuery() (string, error) {
	lines := []string{}
	if c := b.e.Comment(spec.GraphqlComment); c != "" {
		lines = append(lines, c)
	}
	lines = append(lines, fmt.Sprintf(`enum %s {`, b.e.Name()))

	for _, v := range b.e.Values() {
		if c := v.Comment(spec.GraphqlComment); c != "" {
			lines = append(lines, c)
		}
		lines = append(lines, fmt.Sprintf("  %s", v.Name()))
	}

	lines = append(lines, fmt.Sprintf("} # message %s in %s\n", b.e.Name(), b.e.Filename()))
	return strings.Join(lines, "\n"), nil
}

func (b *Enum) BuildProgram() (string, error) {
	var values []string

	for _, v := range b.e.Values() {
		values = append(values, strings.TrimSpace(fmt.Sprintf(`
			%s
			"%s": &graphql.EnumValueConfig{
				Value: %d,
			},`,
			v.Comment(spec.GoComment),
			v.Name(),
			v.Number(),
		)))
	}
	return strings.TrimSpace(fmt.Sprintf(`
%s
var %s = graphql.NewEnum(graphql.EnumConfig{
	Name: "%s",
	Values: graphql.EnumValueConfigMap{
		%s
	},
}) // message %s in %s`,
		b.e.Comment(spec.GoComment),
		spec.PrefixEnum(b.e.Name()),
		b.e.Name(),
		strings.TrimSpace(strings.Join(values, "\n")),
		b.e.Name(),
		b.e.Filename(),
	)) + "\n", nil
}
