package format

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/builder"
)

type Program struct {
	items []builder.Builder
	out   *bytes.Buffer
}

func NewProgram(bs []builder.Builder) *Program {
	return &Program{
		items: bs,
		out:   new(bytes.Buffer),
	}
}

func (p *Program) write(line string) {
	io.WriteString(p.out, line+"\n")
}

func (p *Program) Format(file string) *plugin.CodeGeneratorResponse_File {
	for _, item := range p.items {
		if line := item.BuildProgram(); line != "" {
			p.write(line)
		}
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(file),
		Content: proto.String(p.out.String()),
	}
}
