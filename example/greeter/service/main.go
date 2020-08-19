package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ysugimoto/grpc-graphql-gateway/example/greeter/greeter"
	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	return &greeter.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *Server) SayGoodbye(ctx context.Context, req *greeter.GoodbyeRequest) (*greeter.GoodbyeReply, error) {
	return &greeter.GoodbyeReply{
		Message: fmt.Sprintf("Good-bye, %s!", req.GetName()),
	}, nil
}

func main() {
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	server := grpc.NewServer()
	greeter.RegisterGreeterServer(server, &Server{})
	if err := server.Serve(conn); err != nil {
		log.Println(err)
	}
}
