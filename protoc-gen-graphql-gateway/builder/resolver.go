package builder

import (
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
	// TODO: implement
	return `func(p graphql.ResolveParams) (interface{}, error) {
		return nil, nil
	}`
}
