package spec

import (
	// nolint: staticcheck
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Service spec wraps ServiceDescriptorProto with GraphqlService option.
type Service struct {
	descriptor *descriptor.ServiceDescriptorProto
	Option     *graphql.GraphqlService
	*File
	paths   []int
	methods []*Method

	Queries   []*Query
	Mutations []*Mutation
}

func NewService(
	d *descriptor.ServiceDescriptorProto,
	f *File,
	paths ...int,
) *Service {

	var o *graphql.GraphqlService
	if opts := d.GetOptions(); opts != nil {
		if ext, err := proto.GetExtension(opts, graphql.E_Service); err == nil {
			if service, ok := ext.(*graphql.GraphqlService); ok {
				o = service
			}
		}
	}

	s := &Service{
		descriptor: d,
		Option:     o,
		File:       f,
		paths:      paths,
		methods:    make([]*Method, 0),
		Queries:    make([]*Query, 0),
		Mutations:  make([]*Mutation, 0),
	}

	for i, m := range d.GetMethod() {
		ps := make([]int, len(paths))
		copy(ps, paths)
		s.methods = append(s.methods, NewMethod(m, s, append(ps, 4, i)...))
	}
	return s
}

func (s *Service) Comment() string {
	return s.File.getComment(s.paths)
}

func (s *Service) Name() string {
	return s.descriptor.GetName()
}

func (s *Service) Methods() []*Method {
	return s.methods
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
