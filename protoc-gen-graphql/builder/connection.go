package builder

import (
	"fmt"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Connection builder is used for Go program generation.
// This builder geneartes gRPC connection function for each services.
type Connection struct {
	s *spec.Service
}

func NewConnection(s *spec.Service) *Connection {
	return &Connection{
		s: s,
	}
}

func (c *Connection) BuildSchema() (string, error) {
	return "", nil
}

func (c *Connection) BuildProgram() (string, error) {
	host := c.s.Host()
	var option string
	if c.s.Insecure() {
		option = ", grpc.WithInsecure()"
	}
	return fmt.Sprintf(`
// Create gRPC connection to host which specified via Service directive.
// If you registered handler via ReegisterXXXGraphqlHandler with your *grpc.ClientConn,
// this function won't be called.
func create%sConnection() (*grpc.ClientConn, error) {
	return grpc.Dial("%s"%s)
}`,
		c.s.Name(),
		host,
		option,
	), nil
}
