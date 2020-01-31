package builder

import (
	"fmt"
	"strings"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
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
	return ""
}
