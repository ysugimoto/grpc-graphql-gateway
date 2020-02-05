package builder

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

var handlerTemplate = `

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you worried about open/close performance for each handling graphql request,
// then you can call Register%sHandler with *grpc.ClientConn manually.
func Register%sGraphql(mux *runtime.ServeMux) {
	Register%sGraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your client connection, so that we reuse that and never close connection inside.
// You need to close it maunally when appication will terminate.
func Register%sGraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) {
	mux.AddQueryField(getQueryFields(conn))
	mux.AddMutationField(getMutationFields(conn))
}`

// Handler Builder is used only for Go program generation.
// This builder generates export functions which register service
type Handler struct {
	pkgName string
}

func NewHandler(p string) *Handler {
	return &Handler{
		pkgName: p,
	}
}

func (b *Handler) BuildQuery() (string, error) {
	return "", nil
}

func (b *Handler) BuildProgram() (string, error) {
	n := strcase.ToCamel(b.pkgName)

	return fmt.Sprintf(handlerTemplate, n, n, n, n), nil
}
