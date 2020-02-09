package generator

import (
	"errors"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/resolver"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Generator is struct for analyzing protobuf definition
// and factory graphql definition in protobuf to generate,
// and collect builders with expected order.
type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(
	files []*spec.File,
	args *spec.Params,
) ([]*plugin.CodeGeneratorResponse_File, error) {

	var genFiles []*plugin.CodeGeneratorResponse_File

	queries, mutations, err := g.analyzeMethods(files)
	if err != nil {
		return nil, err
	} else if len(queries) == 0 && len(mutations) == 0 {
		return nil, errors.New("nothing to generate queries")
	}

	r := resolver.New(files)

	// to work this line, query=[outdir] argument is required
	if args.QueryOut != "" {
		// file, err := g.generateSchema(
		// 	r,
		// 	args.QueryOut,
		// 	queries.Collect(),
		// 	mutations.Collect(),
		// )
		// if err != nil {
		// 	return nil, err
		// }
		// genFiles = append(genFiles, file)
	}

	// Generate go program for each query definitions in package
	for pkg, qs := range queries {
		var ms []*spec.Method
		if v, ok := mutations[pkg]; ok {
			ms = v
		}
		file, err := NewProgram(pkg).Generate(r, qs, ms)
		if err != nil {
			return nil, err
		}
		genFiles = append(genFiles, file)
	}

	return genFiles, nil
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

/*
// generateSchema generates GraphQL schema definition.
func (g *Generator) generateSchema(
	r *resolver.Resolver,
	outDir string,
	qs []*spec.Method,
	ms []*spec.Method,
) (*plugin.CodeGeneratorResponse_File, error) {
	types, err := r.ResolveTypes(qs, ms)
	if err != nil {
		return nil, err
	}
	builders := []builder.Builder{
		builder.NewQuery(r.Find, r.FindEnum, qs),
		builder.NewMutation(r.Find, r.FindEnum, ms),
	}
	builders = append(builders, types...)
	schema := NewSchema(builders)
	return schema.Format(filepath.Join(outDir, "./schema.graphql"))
}

// generateProgram generates Go code.
func (g *Generator) generateProgram(
	r *resolver.Resolver,
	pkgName string,
	qs []*spec.Method,
	ms []*spec.Method,
) (*plugin.CodeGeneratorResponse_File, error) {
	types, err := r.ResolveTypes(qs, ms)
	if err != nil {
		return nil, err
	}

	var imports []string
	packages, err := r.ResolvePackages(qs, ms)
	if err != nil {
		return nil, err
	}
	for _, p := range packages {
		if p == pkgName {
			continue
		}
		imports = append(imports, p)
	}

	basename := filepath.Base(pkgName)
	builders := []builder.Builder{
		builder.NewPackage(basename),
		builder.NewImport(imports),
	}
	builders = append(builders, types...)
	builders = append(
		builders,
		builder.NewQuery(r.Find, r.FindEnum, qs),
		builder.NewMutation(r.Find, r.FindEnum, ms),
		builder.NewHandler(basename),
	)
	program := NewProgram(builders)
	return program.Format(
		fmt.Sprintf("%s/%s.graphql.go", pkgName, basename),
	)
}
*/
