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
	tpl "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/template"
)

// Generator is struct for analyzing protobuf definition
// and factory graphql definition in protobuf to generate.
type Generator struct {
	files    []*spec.File
	args     *spec.Params
	messages map[string]*spec.Message
	enums    map[string]*spec.Enum
}

func New(files []*spec.File, args *spec.Params) *Generator {
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
		files:    files,
		args:     args,
		messages: messages,
		enums:    enums,
	}
}

func (g *Generator) Generate(
	tmpl string,
	fs []string,
) (
	[]*plugin.CodeGeneratorResponse_File,
	error,
) {
	services, err := g.analyzeService()
	if err != nil {
		return nil, err
	}

	var outFiles []*plugin.CodeGeneratorResponse_File
	for _, f := range g.files {
		for _, v := range fs {
			if f.Filename() != v {
				continue
			}
			s, ok := services[f.Package()]
			if !ok {
				continue
			}
			// mark as same package defininition in file
			g.analyzeEnum(f)
			if err := g.analyzeMessage(f); err != nil {
				return nil, err
			}

			file, err := g.generateFile(f, tmpl, s)
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
	services []*spec.Service,
) (
	*plugin.CodeGeneratorResponse_File,
	error,
) {
	var types, inputs, interfaces []*spec.Message
	var enums []*spec.Enum
	var packages []*spec.Package
	stack := make(map[string]struct{})

	for _, m := range g.messages {
		if m.IsDepended(spec.DependTypeMessage, file.Package()) {
			if file.Package() == m.Package() || spec.IsGooglePackage(m) {
				types = append(types, m)
			} else if _, ok := stack[m.Package()]; !ok {
				packages = append(packages, spec.NewPackage(m.GoPackage()))
				stack[m.Package()] = struct{}{}
			}
		}
		if m.IsDepended(spec.DependTypeInput, file.Package()) {
			inputs = append(inputs, m)
		}
		if m.IsDepended(spec.DependTypeInterface, file.Package()) {
			interfaces = append(interfaces, m)
		}
	}

	for _, e := range g.enums {
		if e.IsDepended(spec.DependTypeEnum, file.Package()) {
			if file.Package() == e.Package() || spec.IsGooglePackage(e) {
				enums = append(enums, e)
			} else if _, ok := stack[e.Package()]; !ok {
				packages = append(packages, spec.NewPackage(e.GoPackage()))
				stack[e.Package()] = struct{}{}
			}
		}
	}

	root := spec.NewPackage(file.GoPackage())
	t := &tpl.Template{
		RootPackage: root,
		Packages:    packages,
		Types:       types,
		Enums:       enums,
		Inputs:      inputs,
		Interfaces:  interfaces,
		Services:    services,
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

func (g *Generator) analyzeMessage(file *spec.File) error {
	for _, m := range g.messages {
		if m.Package() != file.Package() {
			continue
		}
		m.Depend(spec.DependTypeMessage, file.Package())
		if err := g.analyzeFields(file.Package(), m, m.Fields(), false); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) analyzeEnum(file *spec.File) error {
	for _, e := range g.enums {
		if e.Package() != file.Package() {
			continue
		}
		e.Depend(spec.DependTypeEnum, file.Package())
	}
	return nil
}

func (g *Generator) analyzeService() (
	map[string][]*spec.Service,
	error,
) {
	services := make(map[string][]*spec.Service)

	for _, f := range g.files {
		for _, s := range f.Services() {
			for _, m := range s.Methods() {
				if m.Query == nil && m.Mutation == nil {
					continue
				}
				var input, output *spec.Message
				var ok bool

				if input, ok = g.messages[m.Input()]; !ok {
					return nil, errors.New("failed to resolve input message: " + m.Input())
				}
				if output, ok = g.messages[m.Output()]; !ok {
					return nil, errors.New("failed to resolve output message: " + m.Output())
				}

				if m.Query != nil {
					q := spec.NewQuery(m, input, output)
					if err := g.analyzeQuery(f, q); err != nil {
						return nil, err
					}
					s.Queries = append(s.Queries, q)
				}
				if m.Mutation != nil {
					mu := spec.NewMutation(m, input, output)
					if err := g.analyzeMutation(f, mu); err != nil {
						return nil, err
					}
					s.Mutations = append(s.Mutations, mu)
				}
			}

			if len(s.Queries) > 0 || len(s.Mutations) > 0 {
				services[f.Package()] = append(services[f.Package()], s)
			}
		}
	}
	log.Println("services: ", services)

	return services, nil
}

func (g *Generator) analyzeQuery(f *spec.File, q *spec.Query) error {
	log.Printf("package %s depends on query request %s", f.Package(), q.Input.FullPath())
	q.Input.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Input, q.PluckRequest(), false); err != nil {
		return err
	}

	q.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Output, q.PluckResponse(), false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeMutation(f *spec.File, m *spec.Mutation) error {
	log.Printf("package %s depends on mutation request %s", f.Package(), m.Input.FullPath())
	m.Input.Depend(spec.DependTypeInput, f.Package())
	if err := g.analyzeFields(f.Package(), m.Input, m.PluckRequest(), true); err != nil {
		return err
	}
	m.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), m.Output, m.PluckResponse(), false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeFields(rootPkg string, orig *spec.Message, fields []*spec.Field, asInput bool) error {
	for _, f := range fields {
		switch f.Type() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m, ok := g.messages[f.TypeName()]
			if !ok {
				return errors.New("failed to resolve field message type: " + f.TypeName())
			}
			f.DependType = m
			if asInput {
				log.Printf("package %s depends on input %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeInput, rootPkg)
			} else {
				log.Printf("package %s depends on message %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeMessage, rootPkg)
			}
			// Original message has cyclic dependency
			if m == orig {
				log.Printf("%s has cyclic dependencies of field %s\n", m.Name(), f.Name())
				f.IsCyclic = true
				m.Depend(spec.DependTypeInterface, rootPkg)
			} else {
				// Guard from recursive with infinite loop
				g.analyzeFields(rootPkg, m, m.Fields(), asInput)
			}
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			e, ok := g.enums[f.TypeName()]
			if !ok {
				return errors.New("failed to resolve field enum name: " + f.TypeName())
			}
			f.DependType = e
			log.Printf("package %s depends on enum %s", rootPkg, e.FullPath())
			e.Depend(spec.DependTypeEnum, rootPkg)
		}
	}
	return nil
}
