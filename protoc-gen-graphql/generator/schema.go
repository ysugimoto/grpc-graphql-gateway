package generator

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/builder"
)

// Schema generator is used for generating GraphQL schema.
type Schema struct {
	items []builder.Builder
	out   *bytes.Buffer
}

func NewSchema(bs []builder.Builder) *Schema {
	return &Schema{
		items: bs,
		out:   new(bytes.Buffer),
	}
}

func (s *Schema) write(line string) {
	io.WriteString(s.out, line+"\n")
}

func (s *Schema) Format(file string) (*plugin.CodeGeneratorResponse_File, error) {
	for _, item := range s.items {
		// call BuildQuery() for each builder
		if line, err := item.BuildQuery(); err != nil {
			return nil, err
		} else if line != "" {
			s.write(line)
		}
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(file),
		Content: proto.String(s.out.String()),
	}, nil
}
