package builder

import (
	"fmt"

	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Resolver struct {
	q *types.QuerySpec
}

func NewResolver(q *types.QuerySpec) *Resolver {
	return &Resolver{
		q: q,
	}
}

func (b *Resolver) BuildSchema() string {
	return ""
}

func (b *Resolver) BuildProgram() string {
	pkgName := filepath.Base(b.q.Input.GoPackageName())
	if pkgName == "main" {
		pkgName = ""
	} else {
		pkgName += "."
	}

	var response string
	if expose, _ := b.q.GetExposeField(); expose != nil {
		response = ".Get" + strcase.ToCamel(expose.GetName()) + "()"
	}

	// TODO: implement
	return fmt.Sprintf(`func(p graphql.ResolveParams) (interface{}, error) {
		client := %sNew%sClient(conn)
		resp,err := client.%s(
			p.Context,
			&%s{
				%s
			},
		)
		if err != nil {
			return nil, err
		}
		return resp%s, nil
	}`,
		pkgName,
		b.q.Service.GetName(),
		b.q.Method.GetName(),
		b.q.Input.StructName(false),
		"",
		response,
	)
}
