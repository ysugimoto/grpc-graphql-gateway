package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ysugimoto/grpc-graphql-gateway/example/greeter/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	greeter.GreeterServer
}

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

// StreamGreetings is a server-side streaming RPC method that sends multiple greetings
// to the client in response to a single request.
func (s *Server) StreamGreetings(req *greeter.HelloRequest, stream greeter.Greeter_StreamGreetingsServer) error {
	// Get the name from the request
	name := req.GetName()
	if name == "" {
		return status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	// Define a list of greeting messages to send
	greetings := []string{
		fmt.Sprintf("Hello, %s!", name),
		fmt.Sprintf("Greetings, %s!", name),
		fmt.Sprintf("Good day, %s!", name),
		fmt.Sprintf("Welcome, %s!", name),
		fmt.Sprintf("Hi there, %s!", name),
	}

	// Stream each greeting to the client
	for _, greeting := range greetings {
		// Create the response message
		reply := &greeter.HelloReply{
			Message: greeting,
		}

		// Send the message to the client
		if err := stream.Send(reply); err != nil {
			return status.Errorf(codes.Internal, "failed to send greeting: %v", err)
		}

		// Add a small delay between messages (optional)
		time.Sleep(time.Millisecond * 200) // nolint:mnd
	}

	return nil
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
