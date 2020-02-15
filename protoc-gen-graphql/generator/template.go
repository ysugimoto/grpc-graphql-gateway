package generator

import (
	"bytes"
	"fmt"
	"strings"

	"go/format"
	"io/ioutil"
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
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

	r  *Resolver
	qs []*spec.Method
	ms []*spec.Method
	gt GenerationType
}

func NewTemplate(
	gt GenerationType,
	pkgName string,
	r *Resolver,
	qs, ms []*spec.Method,
) *Template {
	return &Template{
		RootPackage: spec.NewPackage(pkgName),
		Queries:     make([]*spec.Query, 0),
		Mutations:   make([]*spec.Mutation, 0),
		Packages:    make([]*spec.Package, 0),
		Types:       make([]*spec.Message, 0),
		Enums:       make([]*spec.Enum, 0),
		Inputs:      make([]*spec.Message, 0),

		r:  r,
		qs: qs,
		ms: ms,
		gt: gt,
	}
}

func (t *Template) Execute(
	tmplString string,
) (*plugin.CodeGeneratorResponse_File, error) {

	if err := t.generateTemplateParams(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch t.gt {
	case GenerationTypeGo:
		t.filterSamePackages()
		if tmpl, err := template.New("program").Parse(tmplString); err != nil {
			return nil, err
		} else if err := tmpl.Execute(buf, t); err != nil {
			return nil, err
		}
		out, err := format.Source(buf.Bytes())
		if err != nil {
			ioutil.WriteFile("/tmp/"+t.RootPackage.Name+".go", buf.Bytes(), 0666)
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
	case GenerationTypeSchema:
		if tmpl, err := template.New("schema").Parse(tmplString); err != nil {
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

func (t *Template) generateTemplateParams() (err error) {
	var pkgs []*spec.Package
	t.Types, t.Enums, t.Inputs, pkgs, err = t.r.ResolveTypes(t.qs, t.ms)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if pkg.Path == t.RootPackage.Path {
			continue
		}
		t.Packages = append(t.Packages, pkg)
	}

	if len(t.qs) > 0 {
		m := t.qs[0]
		t.Service = m.Service
	} else if len(t.ms) > 0 {
		m := t.ms[0]
		t.Service = m.Service
	}

	for _, q := range t.qs {
		t.Queries = append(t.Queries, spec.NewQuery(
			q,
			t.r.Find(q.Input()),
			t.r.Find(q.Output()),
		))
	}

	for _, m := range t.ms {
		t.Mutations = append(t.Mutations, spec.NewMutation(
			m,
			t.r.Find(m.Input()),
			t.r.Find(m.Output()),
		))
	}

	return nil
}

func (t *Template) filterSamePackages() {
	var types, inputs []*spec.Message
	var enums []*spec.Enum

	for _, v := range t.Types {
		if v.GoPackage() == t.RootPackage.Path ||
			strings.HasPrefix(v.Package(), "google.protobuf") {
			types = append(types, v)
		}
	}
	t.Types = types

	for _, v := range t.Enums {
		if v.GoPackage() == t.RootPackage.Path ||
			strings.HasPrefix(v.Package(), "google.protobuf") {
			enums = append(enums, v)
		}
	}
	t.Enums = enums

	for _, v := range t.Inputs {
		if v.GoPackage() == t.RootPackage.Path ||
			strings.HasPrefix(v.Package(), "google.protobuf") {
			inputs = append(inputs, v)
		}
	}
	t.Inputs = inputs
}
