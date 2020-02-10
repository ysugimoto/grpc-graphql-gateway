package generator

import (
	"bytes"
	"fmt"

	"go/format"
	_ "io/ioutil"
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

type TemplateType int

const (
	TemplateTypeSchema TemplateType = iota
	TemplateTypeGo
)

// Program generator is used for generating Go code.
type Template struct {
	RootPackage *spec.Package
	Service     *spec.Service

	Packages  []*spec.Package
	Types     []*spec.Message
	Enums     []*spec.Enum
	Inputs    []*spec.Message
	Queries   []*spec.Query
	Mutations []*spec.Mutation
}

func NewTemplate(pkgName string) *Template {
	return &Template{
		RootPackage: spec.NewPackage(pkgName),
		Queries:     make([]*spec.Query, 0),
		Mutations:   make([]*spec.Mutation, 0),
		Packages:    make([]*spec.Package, 0),
		Types:       make([]*spec.Message, 0),
		Enums:       make([]*spec.Enum, 0),
		Inputs:      make([]*spec.Message, 0),
	}
}

func (t *Template) Generate(
	outputType TemplateType,
	r *Resolver,
	qs, ms []*spec.Method,
) (*plugin.CodeGeneratorResponse_File, error) {

	if err := t.generateTemplateParams(r, qs, ms); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch outputType {
	case TemplateTypeGo:
		if tmpl, err := template.New("program").Parse(goTemplate); err != nil {
			return nil, err
		} else if err := tmpl.Execute(buf, t); err != nil {
			return nil, err
		}
		// ioutil.WriteFile("/tmp/"+.RootPackage.Name+".go", buf.Bytes(), 0666)
		out, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, err
		}
		return &plugin.CodeGeneratorResponse_File{
			Name: proto.String(fmt.Sprintf(
				"%s/%s.graphql.go",
				t.RootPackage.Path,
				t.RootPackage.Name,
			)),
			Content: proto.String(string(out)),
		}, nil
	case TemplateTypeSchema:
		if tmpl, err := template.New("schema").Parse(schemaTemplate); err != nil {
			return nil, err
		} else if err := tmpl.Execute(buf, t); err != nil {
			return nil, err
		}
		return &plugin.CodeGeneratorResponse_File{
			Name: proto.String(fmt.Sprintf(
				"%s/schema.graphql",
				t.RootPackage.Name,
			)),
			Content: proto.String(buf.String()),
		}, nil
	default:
		return nil, fmt.Errorf("unexpected template type provided")
	}

}

func (t *Template) generateTemplateParams(
	r *Resolver,
	qs, ms []*spec.Method,
) (err error) {
	var pkgs []*spec.Package
	t.Types, t.Enums, t.Inputs, pkgs, err = r.ResolveTypes(qs, ms)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if pkg.Path == t.RootPackage.Path {
			continue
		}
		t.Packages = append(t.Packages, pkg)
	}

	if len(qs) > 0 {
		m := qs[0]
		t.Service = m.Service
	} else if len(ms) > 0 {
		m := ms[0]
		t.Service = m.Service
	}

	for _, q := range qs {
		t.Queries = append(t.Queries, spec.NewQuery(
			q,
			r.Find(q.Input()),
			r.Find(q.Output()),
		))
	}

	for _, m := range ms {
		t.Mutations = append(t.Mutations, spec.NewMutation(
			m,
			r.Find(m.Input()),
			r.Find(m.Output()),
		))
	}

	return nil
}
