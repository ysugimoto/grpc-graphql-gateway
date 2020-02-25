package generator

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

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

func (g *Generator) Generate(tmpl string) (
	[]*plugin.CodeGeneratorResponse_File,
	error,
) {
	services, err := g.analyzeService()
	if err != nil {
		return nil, err
	}

	var outFiles []*plugin.CodeGeneratorResponse_File
	var queries []*spec.Query
	var mutations []*spec.Mutation

	for _, f := range g.files {
		if _, ok := services[f.Package()]; !ok {
			continue
		}
		// mark as same package defininition in file
		g.analyzeEnum(f)
		if err := g.analyzeMessage(f); err != nil {
			return nil, err
		}
	}

	for _, s := range services {
		for _, v := range s {
			queries = append(queries, v.Queries...)
			mutations = append(mutations, v.Mutations...)
		}
	}

	file, err := g.generateFile(tmpl, queries, mutations)
	if err != nil {
		return nil, err
	}
	outFiles = append(outFiles, file)
	return outFiles, nil
}

func (g *Generator) generateFile(
	tmpl string,
	queries []*spec.Query,
	mutations []*spec.Mutation,
) (
	*plugin.CodeGeneratorResponse_File,
	error,
) {
	var types, inputs []*spec.Message
	var enums []*spec.Enum

	for _, m := range g.messages {
		deps := m.GetDependendencies()
		if len(deps["message"]) > 0 {
			types = append(types, m)
		}
		if len(deps["input"]) > 0 {
			inputs = append(inputs, m)
		}
	}

	for _, e := range g.enums {
		deps := e.GetDependendencies()
		if len(deps["enum"]) > 0 {
			enums = append(enums, e)
		}
	}

	t := &tpl.Template{
		Types:     types,
		Enums:     enums,
		Inputs:    inputs,
		Queries:   queries,
		Mutations: mutations,
	}
	buf := new(bytes.Buffer)
	if tmpl, err := template.New("schema").Parse(tmpl); err != nil {
		return nil, err
	} else if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s/schema.graphql", strings.TrimSuffix(g.args.QueryOut, "/"))),
		Content: proto.String(buf.String()),
	}, nil
}

func (g *Generator) analyzeMessage(file *spec.File) error {
	for _, m := range g.messages {
		if m.Package() != file.Package() {
			continue
		}
		m.Depend(spec.DependTypeMessage, file.Package())
		if err := g.analyzeFields(file.Package(), m.Fields(), false); err != nil {
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
			services[f.Package()] = []*spec.Service{}

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

			if len(s.Queries) > 0 && len(s.Mutations) > 0 {
				services[f.Package()] = append(services[f.Package()], s)
			}
		}
	}

	return services, nil
}

func (g *Generator) analyzeQuery(f *spec.File, q *spec.Query) error {
	log.Printf("package %s depends on query request %s", f.Package(), q.Input.FullPath())
	q.Input.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.PluckRequest(), false); err != nil {
		return err
	}

	q.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.PluckResponse(), false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeMutation(f *spec.File, m *spec.Mutation) error {
	log.Printf("package %s depends on mutation request %s", f.Package(), m.Input.FullPath())
	m.Input.Depend(spec.DependTypeInput, f.Package())
	if err := g.analyzeFields(f.Package(), m.PluckRequest(), true); err != nil {
		return err
	}
	m.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), m.PluckResponse(), false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeFields(rootPkg string, fields []*spec.Field, asInput bool) error {
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
			g.analyzeFields(rootPkg, m.Fields(), asInput)
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
