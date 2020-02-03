package runtime

import (
	"google.golang.org/grpc"
)

const (
	optNameErrorHandler = "errorhandler"
	optNameAllowCORS    = "allowcors"
	optNameGrpcOption   = "grpcoption"
)

type Option struct {
	name  string
	value interface{}
}

func WithErrorHandler(eh GraphqlErrorHandler) Option {
	return Option{
		name:  optNameErrorHandler,
		value: eh,
	}
}

func WithCORS() Option {
	return Option{
		name:  optNameAllowCORS,
		value: true,
	}
}

func WithGrpcOption(opts ...grpc.DialOption) Option {
	return Option{
		name:  optNameGrpcOption,
		value: opts,
	}
}
