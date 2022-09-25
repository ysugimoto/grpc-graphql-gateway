package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"

	"go/format"
	"io/ioutil"
	"text/template"

	// nolint: staticcheck
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

type Template struct {
	RootPackage *spec.Package

	Packages   []*spec.Package
	Types      []*spec.Message
	Interfaces []*spec.Message
	Enums      []*spec.Enum
	Inputs     []*spec.Message
	Services   []*spec.Service
}

// Generator is struct for analyzing protobuf definition
// and factory graphql definition in protobuf to generate.
type Generator struct {
	files    []*spec.File
	args     *spec.Params
	messages map[string]*spec.Message
	enums    map[string]*spec.Enum
	logger   *Logger
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

	w := ioutil.Discard
	if args.Verbose {
		w = os.Stderr
	}

	return &Generator{
		files:    files,
		args:     args,
		messages: messages,
		enums:    enums,
		logger:   NewLogger(w),
	}
}

func (g *Generator) Generate(tmpl string, fs []string) ([]*plugin.CodeGeneratorResponse_File, error) {
	services, err := g.analyzeServices()
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

			// mark as same package definition in file
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

// nolint: gocognit, funlen, gocyclo
func (g *Generator) generateFile(file *spec.File, tmpl string, services []*spec.Service) (
	*plugin.CodeGeneratorResponse_File,
	error,
) {

	var types, inputs, interfaces []*spec.Message
	var enums []*spec.Enum
	var packages []*spec.Package

	for _, m := range g.messages {
		// skip empty field message, otherwise graphql-go raise error
		if len(m.Fields()) == 0 {
			continue
		}
		if m.IsDepended(spec.DependTypeMessage, file.Package()) {
			switch {
			case file.Package() == m.Package():
				types = append(types, m)
			case spec.IsGooglePackage(m):
				packages = append(packages, spec.NewGooglePackage(m))
			default:
				packages = append(packages, spec.NewPackage(m))
			}
		}
		if m.IsDepended(spec.DependTypeInput, file.Package()) {
			if !spec.IsGooglePackage(m) {
				inputs = append(inputs, m)
			}
		}
		if m.IsDepended(spec.DependTypeInterface, file.Package()) {
			interfaces = append(interfaces, m)
		}
	}

	for _, s := range services {
		for _, q := range s.Queries {
			input, output := q.Input, q.Output
			if input.Package() != file.Package() {
				if spec.IsGooglePackage(input) {
					packages = append(packages, spec.NewGooglePackage(input))
				} else {
					packages = append(packages, spec.NewPackage(input))
				}
			}
			if output.Package() != file.Package() {
				if spec.IsGooglePackage(output) {
					packages = append(packages, spec.NewGooglePackage(output))
				} else {
					packages = append(packages, spec.NewPackage(output))
				}
			}
		}
		for _, m := range s.Mutations {
			input, output := m.Input, m.Output
			if input.Package() != file.Package() {
				if spec.IsGooglePackage(input) {
					packages = append(packages, spec.NewGooglePackage(input))
				} else {
					packages = append(packages, spec.NewPackage(input))
				}
			}
			if output.Package() != file.Package() {
				if spec.IsGooglePackage(output) {
					packages = append(packages, spec.NewGooglePackage(output))
				} else {
					packages = append(packages, spec.NewPackage(output))
				}
			}
		}
	}

	for _, e := range g.enums {
		// skip empty values enum, otherwise graphql-go raise error
		if len(e.Values()) == 0 {
			continue
		}
		if e.IsDepended(spec.DependTypeEnum, file.Package()) {
			if file.Package() == e.Package() || spec.IsGooglePackage(e) {
				enums = append(enums, e)
			} else {
				packages = append(packages, spec.NewPackage(e))
			}
		}
	}

	// drop duplicate packages
	uniquePackages := make([]*spec.Package, 0)
	stack := make(map[string]struct{})
	for _, p := range packages {
		if _, ok := stack[p.Path]; ok {
			continue
		}
		uniquePackages = append(uniquePackages, p)
		stack[p.Path] = struct{}{}
	}

	// Sort by name to avoid to appear some diff on each generation
	sort.Slice(uniquePackages, func(i, j int) bool {
		return uniquePackages[i].Name > uniquePackages[j].Name
	})
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() > types[j].Name()
	})
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name() > enums[j].Name()
	})
	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].Name() > inputs[j].Name()
	})
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Name() > interfaces[j].Name()
	})
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name() > services[j].Name()
	})

	root := spec.NewPackage(file)
	t := &Template{
		RootPackage: root,
		Packages:    uniquePackages,
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
		ioutil.WriteFile("/tmp/"+root.Name+".go", buf.Bytes(), 0666) // nolint: errcheck
		return nil, err
	}

	// If paths=source_relative option is provided, put generated file relatively
	if g.args.IsSourceRelative() {
		return &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.graphql.go", root.GeneratedFilenamePrefix)),
			Content: proto.String(string(out)),
		}, nil
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s/%s.graphql.go", root.Path, root.FileName)),
		Content: proto.String(string(out)),
	}, nil
}

func (g *Generator) getMessage(name string) *spec.Message {
	if v, ok := g.messages[name]; ok {
		return v
	} else if v, ok := g.messages["."+name]; ok {
		return v
	}
	return nil
}

func (g *Generator) getEnum(name string) *spec.Enum {
	if v, ok := g.enums[name]; ok {
		return v
	} else if v, ok := g.enums["."+name]; ok {
		return v
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeMessage(file *spec.File) error {
	for _, m := range g.messages {
		if m.Package() != file.Package() {
			continue
		}
		m.Depend(spec.DependTypeMessage, file.Package())
		m.Depend(spec.DependTypeInput, file.Package())
		if err := g.analyzeFields(file.Package(), m, m.Fields(), false, false); err != nil {
			return err
		}
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeEnum(file *spec.File) {
	for _, e := range g.enums {
		if e.Package() != file.Package() {
			continue
		}
		e.Depend(spec.DependTypeEnum, file.Package())
	}
}

func (g *Generator) analyzeServices() (map[string][]*spec.Service, error) {
	services := make(map[string][]*spec.Service)

	for _, f := range g.files {
		services[f.Package()] = []*spec.Service{}

		for _, s := range f.Services() {
			if err := g.analyzeService(f, s); err != nil {
				return nil, err
			}
			if len(s.Queries) > 0 || len(s.Mutations) > 0 {
				services[f.Package()] = append(services[f.Package()], s)
			}
		}
	}
	return services, nil
}

func (g *Generator) analyzeService(f *spec.File, s *spec.Service) error {
	for _, m := range s.Methods() {
		if m.Schema == nil {
			continue
		}
		var input, output *spec.Message

		if input = g.getMessage(m.Input()); input == nil {
			return errors.New("failed to resolve input message: " + m.Input())
		}
		if output = g.getMessage(m.Output()); output == nil {
			return errors.New("failed to resolve output message: " + m.Output())
		}

		switch m.Schema.GetType() {
		case graphql.GraphqlType_QUERY, graphql.GraphqlType_RESOLVER:
			q := spec.NewQuery(m, input, output, g.args.FieldCamelCase)
			if err := g.analyzeQuery(f, q); err != nil {
				return err
			}
			s.Queries = append(s.Queries, q)
		case graphql.GraphqlType_MUTATION:
			mu := spec.NewMutation(m, input, output, g.args.FieldCamelCase)
			if err := g.analyzeMutation(f, mu); err != nil {
				return err
			}
			s.Mutations = append(s.Mutations, mu)
		}
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeQuery(f *spec.File, q *spec.Query) error {
	g.logger.Write("package %s depends on query request %s", f.Package(), q.Input.FullPath())
	q.Input.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Input, q.PluckRequest(), false, false); err != nil {
		return err
	}

	q.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Output, q.PluckResponse(), false, false); err != nil {
		return err
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeMutation(f *spec.File, m *spec.Mutation) error {
	g.logger.Write("package %s depends on mutation request %s", f.Package(), m.Input.FullPath())
	m.Input.Depend(spec.DependTypeInput, f.Package())
	if err := g.analyzeFields(f.Package(), m.Input, m.PluckRequest(), true, false); err != nil {
		return err
	}
	m.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), m.Output, m.PluckResponse(), false, false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeFields(
	rootPkg string,
	orig *spec.Message,
	fields []*spec.Field,
	asInput,
	recursive bool,
) error {

	for _, f := range fields {
		switch f.Type() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m := g.getMessage(f.TypeName())
			if m == nil {
				return errors.New("failed to resolve field message type: " + f.TypeName())
			}
			f.DependType = m
			if asInput {
				g.logger.Write("package %s depends on input %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeInput, rootPkg)
			} else {
				g.logger.Write("package %s depends on message %s", rootPkg, m.FullPath())
				switch {
				case m == orig:
					g.logger.Write("%s has cyclic dependencies of field %s\n", m.Name(), f.Name())
					f.IsCyclic = true
					m.Depend(spec.DependTypeInterface, rootPkg)
				case !recursive:
					m.Depend(spec.DependTypeMessage, rootPkg)
				default:
					return nil
				}
			}

			// Guard from recursive with infinite loop
			if m != orig {
				if err := g.analyzeFields(rootPkg, m, m.Fields(), asInput, true); err != nil {
					return err
				}
			}
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			e := g.getEnum(f.TypeName())
			if e == nil {
				return errors.New("failed to resolve field enum name: " + f.TypeName())
			}
			f.DependType = e
			g.logger.Write("package %s depends on enum %s", rootPkg, e.FullPath())
			e.Depend(spec.DependTypeEnum, rootPkg)
		}
	}
	return nil
}
