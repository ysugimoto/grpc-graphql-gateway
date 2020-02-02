package generator

import (
	"sort"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Queries map[string][]*types.QuerySpec

func (q Queries) Add(pkg string, qs *types.QuerySpec) {
	if _, ok := q[pkg]; !ok {
		q[pkg] = make([]*types.QuerySpec, 0)
	}
	q[pkg] = append(q[pkg], qs)
}

func (q Queries) Concat() []*types.QuerySpec {
	var stack []*types.QuerySpec
	for _, v := range q {
		stack = append(stack, v...)
	}
	sort.Slice(stack, func(i, j int) bool {
		return stack[i].Option.GetName() < stack[j].Option.GetName()
	})
	return stack
}

type Mutations map[string][]*types.MutationSpec

func (m Mutations) Add(pkg string, ms *types.MutationSpec) {
	if _, ok := m[pkg]; !ok {
		m[pkg] = make([]*types.MutationSpec, 0)
	}
	m[pkg] = append(m[pkg], ms)
}

func (m Mutations) Concat() []*types.MutationSpec {
	var stack []*types.MutationSpec
	for _, v := range m {
		stack = append(stack, v...)
	}
	sort.Slice(stack, func(i, j int) bool {
		return stack[i].Option.GetName() < stack[j].Option.GetName()
	})
	return stack
}

type Types []*types.Message

func (t Types) Sort() []*types.Message {
	sort.Slice(t, func(i, j int) bool {
		return t[i].MessageName() < t[j].MessageName()
	})
	return t
}
