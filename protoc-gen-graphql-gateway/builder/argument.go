package builder

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql-gateway/types"
)

type Argument struct {
	a *types.ArgumentSpec
}

func NewArgument(a *types.ArgumentSpec) *Argument {
	return &Argument{
		a: a,
	}
}

func (b *Argument) BuildQuery() string {
	return ""
}

func (b *Argument) BuildProgram() string {
	fields := b.a.Message.Descriptor.GetField()
	args := make([]string, len(fields))

	for i, f := range fields {
		fieldType := ext.ConvertGoType(f)
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			fieldType = "graphql.NewList(" + fieldType + ")"
		}

		var optional bool
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			optional = opt.GetOptional()
		}
		if !optional {
			fieldType = "graphql.NewNonNull(" + fieldType + ")"
		}
		args[i] = fmt.Sprintf(`
			"%s": &graphql.ArgumentConfig{
				Type: %s,
			},`,
			f.GetName(),
			fieldType,
		)
	}

	return strings.TrimSpace(strings.Join(args, "\n"))
}
