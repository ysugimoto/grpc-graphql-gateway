package runtime

import (
	"regexp"
	"strings"

	"github.com/graphql-go/graphql/gqlerrors"
)

var (
	// gRPC common status error describes following error message.
	// Find submatch of status text and pluck it.
	// See: https://github.com/grpc/grpc-go/blob/master/internal/status/status.go#L146
	grpcBackendErrorMatcher = regexp.MustCompile(`rpc error: code = ([^\s]+).*desc = (.+)`)
)

// Type alias for gqlerrors in order to avoid to import gqlerrors in user's application
type GraphqlError = gqlerrors.FormattedError

type (
	// Custom error handler which is called on graphql result has an error
	GraphqlErrorHandler func(errs []GraphqlError)
)

// Default error handler for addigng error extension code from gRPC error message
func defaultGraphqlErrorHandler(errs []GraphqlError) {
	for i := 0; i < len(errs); i++ {
		e := errs[i]
		m := grpcBackendErrorMatcher.FindStringSubmatch(e.Message)
		if m == nil {
			continue
		}
		e.Message = m[2]
		if e.Extensions == nil {
			e.Extensions = make(map[string]interface{})
		}
		e.Extensions["code"] = strings.ToUpper(m[1])
		errs[i] = e
	}
}
