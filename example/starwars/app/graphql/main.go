package main

import (
	"log"
	"net/http"

	"github.com/ysugimoto/grpc-graphql-gateway/example/starwars/spec/starwars"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	mux := runtime.NewServeMux(runtime.Cors())
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := starwars.RegisterStartwarsServiceGraphqlHandler(mux, nil, "localhost:50051", opts...); err != nil {
		log.Fatalln(err)
	}
	http.Handle("/graphql", mux)
	log.Fatalln(http.ListenAndServe(":8888", nil))
}
