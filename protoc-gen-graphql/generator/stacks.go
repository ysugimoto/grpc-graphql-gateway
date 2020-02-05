package generator

import (
	"sort"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Queries stacks unique methods by each package.
type Queries map[string][]*spec.Method

func (q Queries) Add(pkg string, m *spec.Method) {
	if _, ok := q[pkg]; !ok {
		q[pkg] = make([]*spec.Method, 0)
	}
	q[pkg] = append(q[pkg], m)
}

// Collect() returns flatten each methods
func (q Queries) Collect() []*spec.Method {
	var stack []*spec.Method
	for _, v := range q {
		stack = append(stack, v...)
	}
	sort.Slice(stack, func(i, j int) bool {
		return stack[i].QueryName() < stack[j].QueryName()
	})
	return stack
}

// Mutations stacks unique methods by each package.
type Mutations map[string][]*spec.Method

func (mu Mutations) Add(pkg string, m *spec.Method) {
	if _, ok := mu[pkg]; !ok {
		mu[pkg] = make([]*spec.Method, 0)
	}
	mu[pkg] = append(mu[pkg], m)
}

// Collect() returns flatten each methods
func (mu Mutations) Collect() []*spec.Method {
	var stack []*spec.Method
	for _, v := range mu {
		stack = append(stack, v...)
	}
	sort.Slice(stack, func(i, j int) bool {
		return stack[i].MutationName() < stack[j].MutationName()
	})
	return stack
}
