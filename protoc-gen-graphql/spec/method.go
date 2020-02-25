package spec

import (
	"strings"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Method spec wraps MethodDescriptorProto with GraphqlQuery and GraphqlMutation options.
type Method struct {
	descriptor *descriptor.MethodDescriptorProto
	Service    *Service
	Query      *graphql.GraphqlQuery
	Mutation   *graphql.GraphqlMutation
	*File

	paths []int
}

func NewMethod(
	m *descriptor.MethodDescriptorProto,
	s *Service,
	paths ...int,
) *Method {
	var q *graphql.GraphqlQuery
	var mt *graphql.GraphqlMutation
	if opts := m.GetOptions(); opts != nil {
		if ext, err := proto.GetExtension(opts, graphql.E_Query); err == nil {
			if query, ok := ext.(*graphql.GraphqlQuery); ok {
				q = query
			}
		}
		if ext, err := proto.GetExtension(opts, graphql.E_Mutation); err == nil {
			if mutation, ok := ext.(*graphql.GraphqlMutation); ok {
				mt = mutation
			}
		}
	}

	return &Method{
		descriptor: m,
		Service:    s,
		Query:      q,
		Mutation:   mt,
		File:       s.File,
		paths:      paths,
	}
}

// -- common functions

func (m *Method) Comment() string {
	return m.File.getComment(m.paths)
}

func (m *Method) ServiceName() string {
	return m.Service.Name()
}

func (m *Method) Name() string {
	return m.descriptor.GetName()
}

func (m *Method) Input() string {
	return strings.TrimPrefix(m.descriptor.GetInputType(), ".")
}

func (m *Method) Output() string {
	return strings.TrimPrefix(m.descriptor.GetOutputType(), ".")
}
