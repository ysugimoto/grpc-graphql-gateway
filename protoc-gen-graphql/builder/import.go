package builder

import (
	"fmt"
	"strings"

	"path/filepath"
)

// Import Builder is used only for Go program generation.
// This builder generates import package section in Go program.
// If we need to import more libraries, it should be added here.
type Import struct {
	imports []string
}

func NewImport(imports []string) *Import {
	return &Import{
		imports: imports,
	}
}

func (b *Import) BuildQuery() (string, error) {
	return "", nil
}

func (b *Import) BuildProgram() (string, error) {
	var deps []string

	for _, i := range b.imports {
		if idx := strings.Index(i, "/"); idx > -1 {
			deps = append(deps, fmt.Sprintf(`%s "%s"`, filepath.Base(i), i))
		} else {
			deps = append(deps, i)
		}
	}

	return fmt.Sprintf(`
import (
	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"

	%s
)
`,
		strings.Join(deps, "\n"),
	), nil
}
