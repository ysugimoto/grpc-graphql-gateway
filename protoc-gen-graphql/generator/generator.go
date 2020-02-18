package generator

import (
	"errors"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

type GenerationType int

const (
	GenerationTypeSchema GenerationType = iota
	GenerationTypeGo
)

// Generator is struct for analyzing protobuf definition
// and factory graphql definition in protobuf to generate,
// and collect builders with expected order.
type Generator struct {
	gt GenerationType
}

func New(gt GenerationType) *Generator {
	return &Generator{
		gt: gt,
	}
}

func (g *Generator) Generate(
	files []*spec.File,
	args *spec.Params,
) ([]*Template, error) {

	queries, mutations, err := g.analyzeMethods(files)
	if err != nil {
		return nil, err
	} else if len(queries) == 0 && len(mutations) == 0 {
		return nil, errors.New("nothing to generate queries")
	}

	r := NewResolver(files)
	switch g.gt {
	case GenerationTypeSchema:
		return []*Template{
			NewTemplate(g.gt, "", r, queries.Collect(), mutations.Collect()),
		}, nil
	case GenerationTypeGo:
		r.ResolveDependencies(queries, mutations)
		// Generate go program for each query definitions in package
		var ts []*Template
		for pkg, qs := range queries {
			ms := []*spec.Method{}
			if v, ok := mutations[pkg]; ok {
				ms = v
			}
			ts = append(ts, NewTemplate(g.gt, pkg, r, qs, ms))
		}
		return ts, nil
	default:
		return nil, errors.New("Unecpected GenerationType supplied")
	}
}

// analyzeMethods analyze all protobuf and find Query/Mutation definitions.
func (g *Generator) analyzeMethods(files []*spec.File) (Queries, Mutations, error) {
	queries := Queries{}
	mutations := Mutations{}

	for _, f := range files {
		pkgName := f.GoPackage()

		for _, s := range f.Services() {
			for _, m := range s.Methods() {
				if m.Query != nil {
					queries.Add(pkgName, m)
				}
				if m.Mutation != nil {
					mutations.Add(pkgName, m)
				}
			}
		}
	}
	return queries, mutations, nil
}
