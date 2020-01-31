package builder

import (
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Argument struct {
	a *types.ArgumentSpec
}

func NewArgument(a *types.ArgumentSpec) *Argument {
	return &Argument{
		a: a,
	}
}

func (a *Argument) BuildQuery() string {
	return ""
}

func (a *Argument) BuildProgram() string {
	return ""
}
