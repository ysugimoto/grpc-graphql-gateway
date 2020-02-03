package builder

import (
	"fmt"
	"strings"

	"path/filepath"

	"github.com/iancoleman/strcase"
	ext "github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/extension"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/types"
)

type Resolver struct {
	q *types.QuerySpec
}

func NewResolver(q *types.QuerySpec) *Resolver {
	return &Resolver{
		q: q,
	}
}

func (b *Resolver) BuildSchema() string {
	return ""
}

func (b *Resolver) BuildProgram() string {
	pkgName := filepath.Base(b.q.Input.GoPackageName())
	if pkgName == "main" {
		pkgName = ""
	} else {
		pkgName += "."
	}

	var response string
	if expose, _ := b.q.GetExposeField(); expose != nil {
		response = ".Get" + strcase.ToCamel(expose.GetName()) + "()"
	}

	input := []string{
		fmt.Sprintf("req := &%s{}", b.q.Input.StructName(false)),
	}
	for _, f := range b.q.Input.Descriptor.GetField() {
		var optional bool
		if opt := ext.GraphqlFieldExtension(f); opt != nil {
			optional = opt.GetOptional()
		}
		t := ext.ConvertGoPrimitiveType(f)
		if !optional {
			if t == "int64" {
				// In graphql, argument comes with int, not in64
				// so we need to cast from int to int64 internaly
				input = append(input, fmt.Sprintf(
					`req.%s = int64(p.Args["%s"].(int))`,
					strcase.ToCamel(f.GetName()),
					f.GetName(),
				))
			} else {
				input = append(input, fmt.Sprintf(
					`req.%s = p.Args["%s"].(%s)`,
					strcase.ToCamel(f.GetName()),
					f.GetName(),
					ext.ConvertGoPrimitiveType(f),
				))
			}
		} else {
			if t == "int64" {
				input = append(input, strings.TrimSpace(fmt.Sprintf(`
				  if v, ok := p.Args["%s"]; !ok {
						return nil, errors.New("%s is not found in parameter")
				  } else if arg, ok := v.(int); !ok {
						return nil, errors.New("%s is not found in parameter")
				  } else {
						req.%s = int64(arg)
				}`,
					f.GetName(),
					f.GetName(),
					f.GetName(),
					strcase.ToCamel(f.GetName()),
				)))
			} else {
				input = append(input, strings.TrimSpace(fmt.Sprintf(`
				  if v, ok := p.Args["%s"]; !ok {
						return nil, errors.New("%s is not found in parameter")
				  } else if arg, ok := v.(%s); !ok {
						return nil, errors.New("%s is not found in parameter")
				  } else {
						req.%s = arg
				}`,
					f.GetName(),
					f.GetName(),
					ext.ConvertGoPrimitiveType(f),
					f.GetName(),
					strcase.ToCamel(f.GetName()),
				)))
			}
		}

	}

	var connLogic string
	service := ext.GraphqlServiceExtension(b.q.Service)
	if service == nil || service.GetHost() == "" {
		connLogic = `
		conn := c.Default
		if conn == nil {
			return nil, errors.New("failed to find default grpc connection")
		}`
	} else {
		var option string
		if service.GetInsecure() {
			option = ", grpc.WithInsecure()"
		}
		connLogic = strings.TrimSpace(fmt.Sprintf(`
			var keep bool
			conn := c.Find("%s")
			if conn == nil {
				if c, err := grpc.Dial("%s"%s); err != nil {
					return nil, errors.New("failed to find grpc connection for '%s'")
				} else {
					conn = c
				}
			} else {
				keep = true
			}
			if !keep {
				defer conn.Close()
			}`,
			service.GetHost(),
			service.GetHost(),
			option,
			service.GetHost(),
		))
	}

	// TODO: implement
	return fmt.Sprintf(`func(p graphql.ResolveParams) (interface{}, error) {
		%s
		client := %sNew%sClient(conn)
		%s
		resp,err := client.%s(p.Context, req)
		if err != nil {
			return nil, err
		}
		return resp%s, nil
	}`,
		connLogic,
		pkgName,
		b.q.Service.GetName(),
		strings.Join(input, "\n"),
		b.q.Method.GetName(),
		response,
	)
}
