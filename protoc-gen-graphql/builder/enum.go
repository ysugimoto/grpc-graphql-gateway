package builder

import (
	"fmt"
	"strings"

	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/types"
)

type Enum struct {
	e *types.Enum
}

func NewEnum(e *types.Enum) *Enum {
	return &Enum{
		e: e,
	}
}

func (b *Enum) BuildQuery() string {
	enum := b.e

	lines := []string{
		fmt.Sprintf("# message %s in %s", enum.MessageName(), enum.Filename()),
		fmt.Sprintf(`enum %s {`, enum.Descriptor.GetName()),
	}

	for _, v := range enum.Descriptor.GetValue() {
		lines = append(lines, "  "+v.GetName())
	}

	lines = append(lines, "}\n")

	return strings.Join(lines, "\n")
}

func (b *Enum) BuildProgram() string {
	values := make([]string, len(b.e.Descriptor.GetValue()))
	for i, v := range b.e.Descriptor.GetValue() {
		values[i] = strings.TrimSpace(fmt.Sprintf(`
			"%s": &graphql.EnumValueConfig{
				Value: %d,
			},`,
			v.GetName(),
			v.GetNumber(),
		))
	}
	return fmt.Sprintf(`
var %s = graphql.NewEnum(graphql.EnumConfig{
	Name: "%s",
	Values: graphql.EnumValueConfigMap{
		%s
	},
})`,

		ext.EnumName(b.e.Descriptor.GetName()),
		b.e.Descriptor.GetName(),
		strings.TrimSpace(strings.Join(values, "\n")),
	)
}
