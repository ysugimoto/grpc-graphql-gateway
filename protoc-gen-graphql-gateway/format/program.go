package format

import (
	"errors"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type Program struct {
}

func NewProgram(q Queries, m Mutations, t Types) *Program {
	return &Program{}
}

func (p *Program) Format() (*plugin.CodeGeneratorResponse_File, error) {
	return nil, errors.New("not implemented")
}
