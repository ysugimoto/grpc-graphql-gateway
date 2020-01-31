package resolver

import (
	"errors"
	"log"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/builder"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Resolver struct {
	// factries stacks
	messages map[string]*types.Message
	enums    map[string]*types.Enum
}

func New(req *plugin.CodeGeneratorRequest) *Resolver {
	messages := make(map[string]*types.Message)
	enums := make(map[string]*types.Enum)

	for _, f := range req.GetProtoFile() {
		pkgName := f.GetPackage()
		if strings.HasPrefix(pkgName, "google.protobuf") {
			continue
		}
		for _, m := range f.GetMessageType() {
			key := pkgName + "." + m.GetName()
			messages[key] = types.NewMessage(m, f)
		}
		for _, e := range f.GetEnumType() {
			key := pkgName + "." + e.GetName()
			enums[key] = &types.Enum{e, f}
		}
	}

	return &Resolver{
		messages: messages,
		enums:    enums,
	}
}

func (r *Resolver) FindMessage(names ...string) *types.Message {
	for _, n := range names {
		if m, ok := r.messages[n]; ok {
			return m
		}
	}
	return nil
}

func (r *Resolver) FindEnum(names ...string) *types.Enum {
	for _, n := range names {
		if m, ok := r.enums[n]; ok {
			return m
		}
	}
	return nil
}

func (r *Resolver) ResolveTypes(
	queries []*types.QuerySpec,
	mutations []*types.MutationSpec,
) ([]builder.Builder, error) {

	var builders []builder.Builder
	stack := make(map[string]struct{})

	for _, qs := range queries {
		bs, err := r.resolveMessage(qs.Output, stack, qs.GetExposeField())
		if err != nil {
			return nil, err
		}
		builders = append(builders, bs...)
	}
	for _, mu := range mutations {
		bs, err := r.resolveMessage(mu.Output, stack, mu.GetExposeField())
		if err != nil {
			return nil, err
		}
		builders = append(builders, bs...)
	}
	return builders, nil
}

func (r *Resolver) resolveMessage(
	m *types.Message,
	stack map[string]struct{},
	filterFields ...string,
) ([]builder.Builder, error) {
	var ret []builder.Builder

	if _, ok := stack["m_"+m.Descriptor.GetName()]; !ok {
		ret = append(ret, builder.NewMessage(m))
		stack["m_"+m.Descriptor.GetName()] = struct{}{}
	}
	for _, f := range m.Descriptor.GetField() {
		var ignored bool
		if len(filterFields) > 0 {
			ignored = true
			for _, i := range filterFields {
				if f.GetName() == i {
					ignored = false
					break
				}
			}
		}
		log.Printf("type: %s, field: %s, filters: %v, ignored %t", m.Descriptor.GetName(), f.GetName(), filterFields, ignored)
		if ignored {
			continue
		}
		switch f.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			mm := r.FindMessage(
				f.GetTypeName(),
				strings.TrimPrefix(f.GetTypeName(), "."),
				"."+f.GetTypeName(),
			)
			if mm == nil {
				return nil, errors.New("resolveMessages: undefined message: " + f.GetTypeName())
			}
			nested, err := r.resolveMessage(mm, stack)
			if err != nil {
				return nil, err
			}
			ret = append(ret, nested...)
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			en := r.FindEnum(
				f.GetTypeName(),
				strings.TrimPrefix(f.GetTypeName(), "."),
				"."+f.GetTypeName(),
			)
			if en == nil {
				return nil, errors.New("resolveMessages: undefined enum: " + f.GetTypeName())
			}
			if _, ok := stack["e_"+en.Descriptor.GetName()]; !ok {
				ret = append(ret, builder.NewEnum(en))
				stack["e_"+en.Descriptor.GetName()] = struct{}{}
			}
		default:
			continue
		}
	}
	return ret, nil
}
