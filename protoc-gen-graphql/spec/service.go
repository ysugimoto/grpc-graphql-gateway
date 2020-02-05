package spec

import (
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Service spec wraps ServiceDescriptorProto with GraphqlService option.
type Service struct {
	descriptor *descriptor.ServiceDescriptorProto
	Option     *graphql.GraphqlService
	*File
	paths []int
}

func NewService(
	s *descriptor.ServiceDescriptorProto,
	f *File,
	paths ...int,
) *Service {
	var o *graphql.GraphqlService
	if opts := s.GetOptions(); opts != nil {
		if ext, err := proto.GetExtension(opts, graphql.E_Service); err == nil {
			if service, ok := ext.(*graphql.GraphqlService); ok {
				o = service
			}
		}
	}

	return &Service{
		descriptor: s,
		Option:     o,
		File:       f,
	}
}

func (s *Service) Comment(t CommentType) string {
	return s.File.getComment(s.paths, t)
}

func (s *Service) Name() string {
	return s.descriptor.GetName()
}

func (s *Service) Methods() []*Method {
	var methods []*Method
	for i, m := range s.descriptor.GetMethod() {
		paths := make([]int, len(s.paths))
		copy(paths, s.paths)
		methods = append(methods, NewMethod(m, s, append(paths, 4, i)...))
	}
	return methods
}

func (s *Service) Host() string {
	if s.Option == nil {
		return ""
	}
	return s.Option.GetHost()
}

func (s *Service) Insecure() bool {
	if s.Option == nil {
		return false
	}
	return s.Option.GetInsecure()
}
