package types

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/graphql"
)

type FieldSpec struct {
	Descriptor *descriptor.FieldDescriptorProto
	Option     *graphql.GraphqlField
}
