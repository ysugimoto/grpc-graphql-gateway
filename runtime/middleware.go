package runtime

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

type (
	// MiddlewareFunc type definition
	MiddlewareFunc func(w http.ResponseWriter, r *http.Request) error

	// Custom error handler which is called on graphql result has an error
	GraphqlErrorHandler func(errs gqlerrors.FormattedErrors)
)
