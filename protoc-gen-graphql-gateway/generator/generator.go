package generator

import (
	"errors"
	"log"
	"strings"

	"path/filepath"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/format"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/graphql"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Generator struct {
	req *plugin.CodeGeneratorRequest

	messages  map[string]*types.Message
	queries   format.Queries
	mutations format.Mutations
	types     format.Types
}

func New(req *plugin.CodeGeneratorRequest) *Generator {
	messages := make(map[string]*types.Message)

	for _, f := range req.GetProtoFile() {
		if strings.HasPrefix(f.GetPackage(), "google.protobuf") {
			continue
		}
		for _, m := range f.GetMessageType() {
			key := f.GetPackage() + "." + m.GetName()
			messages[key] = types.NewMessage(m, f)
		}
	}
	return &Generator{
		req:      req,
		messages: messages,

		queries:   format.Queries{},
		mutations: format.Mutations{},
		types:     format.Types{},
	}
}

func (g *Generator) FindMessage(names ...string) *types.Message {
	for _, n := range names {
		if m, ok := g.messages[n]; ok {
			return m
		}
	}
	return nil
}

func (g *Generator) Generate() *plugin.CodeGeneratorResponse {
	resp := &plugin.CodeGeneratorResponse{}

	var genError error
	defer func() {
		if genError != nil {
			msg := genError.Error()
			resp.Error = &msg
		}
	}()

	args := &types.Params{}
	if g.req.Parameter != nil {
		args, genError = types.NewParams(g.req.GetParameter())
		if genError != nil {
			return resp
		}
	}

	for _, f := range g.req.GetProtoFile() {
		pkg := f.GetPackage()
		for _, s := range f.GetService() {
			for _, m := range s.GetMethod() {
				if opt := ext.GraphqlQueryOption(m); opt != nil {
					qs, err := g.AnalyzeQuery(m, opt)
					if err != nil {
						genError = err
						return resp
					}
					g.queries.Add(pkg, qs)
				}
				if opt := ext.GraphqlMutationOption(m); opt != nil {
					ms, err := g.AnalyzeMutation(m, opt)
					if err != nil {
						genError = err
						return resp
					}
					g.mutations.Add(pkg, ms)
				}
			}
		}
	}

	if len(g.queries) == 0 {
		genError = errors.New("nothing to generate queries")
		return resp
	}

	queries := g.queries.Concat()
	mutations := g.mutations.Concat()
	stack := make(map[string]struct{})
	for _, q := range queries {
		g.types = append(g.types, g.resolveMessages(q.Output, stack)...)
	}
	for _, m := range mutations {
		g.types = append(g.types, g.resolveMessages(m.Output, stack)...)
	}

	if args.QueryOut != "" {
		schema := format.NewSchema(queries, mutations, g.types.Sort())
		file, err := schema.Format(filepath.Join(args.QueryOut, "./query.graphql"))
		if err != nil {
			genError = err
			return resp
		}
		resp.File = append(resp.File, file)
	}

	// program := format.NewProgram(g.queries, g.mutations)
	// file, err = program.Format()
	// if err != nil {
	// 	genError = err
	// 	return resp
	// }
	// resp.File = append(resp.File, file)
	return resp
}

func (g *Generator) AnalyzeQuery(
	m *descriptor.MethodDescriptorProto,
	opt *graphql.GraphqlQuery,
) (*types.QuerySpec, error) {
	var req, resp *types.Message
	if req = g.FindMessage(
		m.GetInputType(),
		strings.TrimPrefix(m.GetInputType(), "."),
		"."+m.GetInputType(),
	); req == nil {
		return nil, errors.New("InputType: " + m.GetInputType() + " not exists")
	}
	if resp = g.FindMessage(
		m.GetOutputType(),
		strings.TrimPrefix(m.GetOutputType(), "."),
		"."+m.GetOutputType(),
	); resp == nil {
		return nil, errors.New("OutputType: " + m.GetOutputType() + " not exists")
	}

	return &types.QuerySpec{
		Input:  req,
		Output: resp,
		Option: opt,
	}, nil
}

func (g *Generator) AnalyzeMutation(
	m *descriptor.MethodDescriptorProto,
	opt *graphql.GraphqlMutation,
) (*types.MutationSpec, error) {
	var req, resp *types.Message
	var ok bool
	if req, ok = g.messages[m.GetInputType()]; !ok {
		return nil, errors.New("InputType: " + m.GetInputType() + " not exists")
	}
	if resp, ok = g.messages[m.GetOutputType()]; !ok {
		return nil, errors.New("OutputType: " + m.GetOutputType() + " not exists")
	}

	return &types.MutationSpec{
		Input:  req,
		Output: resp,
		Option: opt,
	}, nil
}

func (g *Generator) resolveMessages(m *types.Message, stack map[string]struct{}) []*types.Message {
	ret := []*types.Message{}
	if _, ok := stack[m.Descriptor.GetName()]; !ok {
		ret = append(ret, m)
		stack[m.Descriptor.GetName()] = struct{}{}
	}

	for _, f := range m.Descriptor.GetField() {
		if f.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}
		mm := g.FindMessage(
			f.GetTypeName(),
			strings.TrimPrefix(f.GetTypeName(), "."),
			"."+f.GetTypeName(),
		)
		if mm == nil {
			log.Println("resolveMessages: undefined: " + f.GetTypeName())
			continue
		}
		ret = append(ret, g.resolveMessages(mm, stack)...)
	}
	return ret
}
