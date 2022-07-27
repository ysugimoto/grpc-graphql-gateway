package main

import (
	"log"

	"net/http"

	"github.com/alehechka/grpc-graphql-gateway/example/greeter/greeter"
	"github.com/alehechka/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := greeter.RegisterGreeterGraphql(mux, "localhost:50051", opts...); err != nil {
		log.Fatalln(err)
	}
	http.Handle("/graphql", mux)
	log.Fatalln(http.ListenAndServe(":8888", nil))
}
