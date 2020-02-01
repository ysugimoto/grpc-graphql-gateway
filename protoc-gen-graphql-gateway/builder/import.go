package builder

import (
	"fmt"
	"strings"

	"path/filepath"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
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
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	%s
)`,
		strings.Join(dendencies, "\n"),
	)
}
