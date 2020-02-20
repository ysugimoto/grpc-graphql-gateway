package generator

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"go/format"
	"io/ioutil"
	"text/template"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
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

	files    []*spec.File
	args     *spec.Params
	messages map[string]*spec.Message
	enums    map[string]*spec.Enum
}

func New(gt GenerationType, files []*spec.File, args *spec.Params) *Generator {
	messages := make(map[string]*spec.Message)
	enums := make(map[string]*spec.Enum)

	for _, f := range files {
		for _, m := range f.Messages() {
			messages[m.FullPath()] = m
		}
		for _, e := range f.Enums() {
			enums[e.FullPath()] = e
		}
	}

	return &Generator{
		gt: gt,

		files:    files,
		args:     args,
		messages: messages,
		enums:    enums,
	}
}

func (g *Generator) Generate(tmpl string, fs []string) ([]*plugin.CodeGeneratorResponse_File, error) {

	queries, mutations, err := g.analyzeProto()
	if err != nil {
		return nil, err
	}

	for _, m := range g.messages {
		log.Printf("%s is depended from %v\n", m.FullPath(), m.GetDependendencies())
	}

	var outFiles []*plugin.CodeGeneratorResponse_File
	for _, f := range g.files {
		for _, v := range fs {
			if f.Filename() != v {
				continue
			}
			qs, ok := queries[f.Package()]
			if !ok {
				qs = []*spec.Query{}
			}
			ms, ok := mutations[f.Package()]
			if !ok {
				ms = []*spec.Mutation{}
			}
			if err := g.analyzeMessage(f); err != nil {
				return nil, err
			}

			file, err := g.generateFile(f, tmpl, qs, ms)
			if err != nil {
				return nil, err
			}
			outFiles = append(outFiles, file)
		}
	}
	return outFiles, nil
}

func (g *Generator) generateFile(
	file *spec.File,
	tmpl string,
	queries []*spec.Query,
	mutations []*spec.Mutation,
) (
	*plugin.CodeGeneratorResponse_File,
	error,
) {
	var types, inputs []*spec.Message
	var enums []*spec.Enum
	var packages []*spec.Package
	stack := make(map[string]struct{})

	for _, m := range g.messages {
		if m.IsDepended(spec.DependTypeMessage, file.Package()) {
			if file.Package() == m.Package() {
				types = append(types, m)
			} else if _, ok := stack[m.Package()]; !ok {
				packages = append(packages, spec.NewPackage(m.GoPackage()))
				stack[m.Package()] = struct{}{}
			}
		}
		if m.IsDepended(spec.DependTypeInput, file.Package()) {
			inputs = append(inputs, m)
		}
	}

	for _, e := range g.enums {
		if e.IsDepended(spec.DependTypeEnum, file.Package()) {
			if file.Package() == e.Package() {
				enums = append(enums, e)
			} else if _, ok := stack[e.Package()]; !ok {
				packages = append(packages, spec.NewPackage(e.GoPackage()))
				stack[e.Package()] = struct{}{}
			}
		}
	}

	root := spec.NewPackage(file.GoPackage())
	for _, t := range types {
		for _, f := range t.Fields() {
			log.Println(f.FieldType(root.Path))
		}
	}

	t := &Template{
		RootPackage: root,
		Packages:    packages,
		Types:       types,
		Enums:       enums,
		Inputs:      inputs,
		Queries:     queries,
		Mutations:   mutations,
	}

	buf := new(bytes.Buffer)
	if tmpl, err := template.New("go").Parse(tmpl); err != nil {
		return nil, err
	} else if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Println(buf.String())
		ioutil.WriteFile("/tmp/"+root.Name+".go", buf.Bytes(), 0666)
		return nil, err
	}
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s/%s.graphql.go", root.Path, root.Name)),
		Content: proto.String(string(out)),
	}, nil
}

func (g *Generator) analyzeProto() (
	map[string][]*spec.Query,
	map[string][]*spec.Mutation,
	error,
) {
	queries := make(map[string][]*spec.Query)
	mutations := make(map[string][]*spec.Mutation)

	for _, f := range g.files {
		for _, s := range f.Services() {
			for _, m := range s.Methods() {
				if m.Query == nil && m.Mutation == nil {
					continue
				}
				input, ok := g.messages[m.Input()]
				if !ok {
					return nil, nil, errors.New("failed to resolve input message: " + m.Input())
				}
				log.Printf("package %s depends on rpc request %s", f.Package(), input.FullPath())
				if m.Query != nil {
					input.Depend(spec.DependTypeMessage, f.Package())
					if err := g.analyzeFields(f.Package(), input.Fields(), false); err != nil {
						return nil, nil, err
					}
				}
				if m.Mutation != nil {
					input.Depend(spec.DependTypeInput, f.Package())
					if err := g.analyzeFields(f.Package(), input.Fields(), true); err != nil {
						return nil, nil, err
					}
				}

				output, ok := g.messages[m.Output()]
				if !ok {
					return nil, nil, errors.New("failed to resolve output message: " + m.Output())
				}

				output.Depend(spec.DependTypeMessage, f.Package())
				if err := g.analyzeFields(f.Package(), m.ExposeQueryFields(output), false); err != nil {
					return nil, nil, err
				}

				if m.Query != nil {
					if _, ok := queries[m.Package()]; !ok {
						queries[m.Package()] = make([]*spec.Query, 0)
					}
					queries[m.Package()] = append(queries[m.Package()], spec.NewQuery(m, input, output))
				}
				if m.Mutation != nil {
					if _, ok := mutations[m.Package()]; !ok {
						mutations[m.Package()] = make([]*spec.Mutation, 0)
					}
					mutations[m.Package()] = append(mutations[m.Package()], spec.NewMutation(m, input, output))
				}

			}
		}
		for _, e := range f.Enums() {
			e.Depend(spec.DependTypeEnum, f.Package())
		}
	}

	return queries, mutations, nil
}

func (g *Generator) analyzeFields(rootPkg string, fields []*spec.Field, asInput bool) error {
	for _, f := range fields {
		switch f.Type() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m, ok := g.messages[f.TypeName()]
			if !ok {
				return errors.New("failed to resolve field message type: " + f.TypeName())
			}
			f.TypeMessage = m
			if asInput {
				log.Printf("package %s depends on input %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeInput, rootPkg)
			} else {
				log.Printf("package %s depends on message %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeMessage, rootPkg)
			}
			g.analyzeFields(rootPkg, m.Fields(), asInput)
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			e, ok := g.enums[f.TypeName()]
			if !ok {
				return errors.New("failed to resolve field enum name: " + f.TypeName())
			}
			f.TypeEnum = e
			log.Printf("package %s depends on enum %s", rootPkg, e.FullPath())
			e.Depend(spec.DependTypeEnum, rootPkg)
		}
	}
	return nil
}
