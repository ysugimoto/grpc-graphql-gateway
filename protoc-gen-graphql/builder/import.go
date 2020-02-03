package builder

import (
	"fmt"
	"strings"

	"path/filepath"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/types"
)

type Import struct {
	qs []*types.QuerySpec
	ms []*types.MutationSpec
}

func NewImport(qs []*types.QuerySpec, ms []*types.MutationSpec) *Import {
	return &Import{
		qs: qs,
		ms: ms,
	}
}

func (b *Import) BuildQuery() string {
	return ""
}

func (b *Import) BuildProgram() string {
	var dendencies []string
	stack := map[string]struct{}{}

	for _, v := range b.qs {
		input := v.Input.GoPackageName()
		if _, ok := stack[input]; !ok {
			dendencies = append(dendencies, fmt.Sprintf(
				`%s "%s"`,
				filepath.Base(input),
				input,
			))
			stack[input] = struct{}{}
		}
		output := v.Output.GoPackageName()
		if _, ok := stack[output]; !ok {
			dendencies = append(dendencies, fmt.Sprintf(
				`%s "%s"`,
				filepath.Base(output),
				output,
			))
			dendencies = append(dendencies, output)
			stack[output] = struct{}{}
		}
	}
	for _, v := range b.ms {
		input := v.Input.GoPackageName()
		if _, ok := stack[input]; !ok {
			dendencies = append(dendencies, fmt.Sprintf(
				`%s "%s"`,
				filepath.Base(input),
				input,
			))
			stack[input] = struct{}{}
		}
		output := v.Output.GoPackageName()
		if _, ok := stack[output]; !ok {
			dendencies = append(dendencies, fmt.Sprintf(
				`%s "%s"`,
				filepath.Base(output),
				output,
			))
			stack[output] = struct{}{}
		}
	}

	return fmt.Sprintf(`
import (
	"net/http"

	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"

	%s
)
`,
		strings.Join(dendencies, "\n"),
	)
}
