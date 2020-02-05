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
			if query, ok := ext.(*graphql.GraphqlMutation); ok {
				mt = query
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

func (m *Method) Comment(t CommentType) string {
	return m.File.getComment(m.paths, t)
}

func (m *Method) ServiceName() string {
	return m.Service.Name()
}

func (m *Method) Name() string {
	return m.descriptor.GetName()
}

func (m *Method) QueryName() string {
	if m.Query == nil {
		return ""
	}
	return m.Query.GetName()
}

func (m *Method) MutationName() string {
	if m.Mutation == nil {
		return ""
	}
	return m.Mutation.GetName()
}

func (m *Method) Input() string {
	return strings.TrimPrefix(m.descriptor.GetInputType(), ".")
}

func (m *Method) Output() string {
	return strings.TrimPrefix(m.descriptor.GetOutputType(), ".")
}

func (m *Method) QueryResponse() *graphql.GraphqlResponse {
	if m.Query == nil {
		return nil
	}
	return m.Query.GetResponse()
}

func (m *Method) MutationResponse() *graphql.GraphqlResponse {
	if m.Mutation == nil {
		return nil
	}
	return m.Mutation.GetResponse()
}

func (m *Method) MutationRequest() *graphql.GraphqlRequest {
	if m.Mutation == nil {
		return nil
	}
	return m.Mutation.GetRequest()
}

func (m *Method) ExposeQuery() string {
	if m.Query == nil {
		return ""
	} else if m.Query.GetResponse() == nil {
		return ""
	}
	return m.Query.GetResponse().GetExpose()
}

func (m *Method) ExposeMutation() string {
	if m.Mutation == nil {
		return ""
	} else if m.Mutation.GetResponse() == nil {
		return ""
	}
	return m.Mutation.GetResponse().GetExpose()
}

func (m *Method) ExposeQueryFields(msg *Message) []*Field {
	expose := m.ExposeQuery()
	if expose == "" {
		return msg.Fields()
	}
	var fields []*Field
	for _, f := range msg.Fields() {
		if expose != f.Name() {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (m *Method) ExposeMutationFields(msg *Message) []*Field {
	expose := m.ExposeMutation()
	if expose == "" {
		return msg.Fields()
	}
	var fields []*Field
	for _, f := range msg.Fields() {
		if expose != f.Name() {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}
