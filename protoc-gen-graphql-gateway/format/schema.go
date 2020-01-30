package format

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Schema struct {
	q   []*types.QuerySpec
	m   []*types.MutationSpec
	t   []*types.Message
	out *bytes.Buffer
}

func NewSchema(q []*types.QuerySpec, m []*types.MutationSpec, t Types) *Schema {
	return &Schema{
		q:   q,
		m:   m,
		t:   t,
		out: new(bytes.Buffer),
	}
}

func (s Schema) write(line string) {
	io.WriteString(s.out, line+"\n")
}

func (s *Schema) Format(file string) (*plugin.CodeGeneratorResponse_File, error) {
	s.write(`type Query {`)
	for _, q := range s.q {
		s.write("  " + q.BuildQuery())
	}
	s.write("}")

	for _, m := range s.m {
		s.write("")
		s.write(m.BuildQuery())
	}

	for _, t := range s.t {
		s.write("")
		s.write(t.BuildQuery())
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(file),
		Content: proto.String(s.out.String()),
	}, nil
}
