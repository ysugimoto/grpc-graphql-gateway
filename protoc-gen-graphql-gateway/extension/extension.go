package extension

import (
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

func GraphqlQueryOption(m *descriptor.MethodDescriptorProto) *graphql.GraphqlQuery {
	opts := m.GetOptions()
	if opts == nil {
		return nil
	}
	ext, err := proto.GetExtension(opts, graphql.E_Query)
	if err != nil {
		return nil
	}
	if v, ok := ext.(*graphql.GraphqlQuery); ok {
		return v
	}
	return nil
}

func GraphqlMutationOption(m *descriptor.MethodDescriptorProto) *graphql.GraphqlMutation {
	opts := m.GetOptions()
	if opts == nil {
		return nil
	}
	ext, err := proto.GetExtension(opts, graphql.E_Mutation)
	if err != nil {
		return nil
	}
	if v, ok := ext.(*graphql.GraphqlMutation); ok {
		return v
	}
	return nil
}

func GraphqlFieldExtension(f *descriptor.FieldDescriptorProto) *graphql.GraphqlField {
	if opts := f.GetOptions(); opts == nil {
		return nil
	} else if ext, err := proto.GetExtension(opts, graphql.E_Field); err != nil {
		return nil
	} else if field, ok := ext.(*graphql.GraphqlField); !ok {
		return nil
	} else {
		return field
	}
}
