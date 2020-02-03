package runtime

import (
	"google.golang.org/grpc"
)

type Connection struct {
	Default     *grpc.ClientConn
	distributed map[string]*grpc.ClientConn
}

func NewConnection(defaultConn *grpc.ClientConn) *Connection {
	return &Connection{
		Default:     defaultConn,
		distributed: map[string]*grpc.ClientConn{},
	}
}

func (c *Connection) Close() {
	if c.Default != nil {
		c.Default.Close()
	}
	for _, v := range c.distributed {
		if v != nil {
			v.Close()
		}
	}
}

func (c *Connection) Distribute(host string, conn *grpc.ClientConn) *Connection {
	c.distributed[host] = conn
	return c
}

func (c *Connection) Find(host string) *grpc.ClientConn {
	if conn, ok := c.distributed[host]; !ok {
		return nil
	} else {
		return conn
	}
}
