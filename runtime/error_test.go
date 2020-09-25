package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorHandler(t *testing.T) {
	rpcError := []GraphqlError{
		{
			Message: "rpc error: code = NotFound desc = Example Message",
		},
	}
	defaultGraphqlErrorHandler(rpcError)
	assert.Len(t, rpcError, 1)
	err := rpcError[0]
	assert.Equal(t, "Example Message", err.Message)
	assert.NotNil(t, err.Extensions)
	ext, ok := err.Extensions["code"]
	assert.True(t, ok)
	assert.Equal(t, "NOTFOUND", ext)
}
