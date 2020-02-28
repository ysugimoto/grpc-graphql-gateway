package main

import (
	"log"
	"net/http"

	// "github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/author"
	"github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/book"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
)

func main() {
	mux := runtime.NewServeMux(runtime.Cors())

	// author.RegisterAuthorServiceGraphqlHandler(mux, nil)
	book.RegisterBookServiceGraphqlHandler(mux, nil)
	http.Handle("/graphql", mux)
	log.Println("GraphQL handler starts on 0.0.0.0:8888...")
	http.ListenAndServe(":8888", nil)
}
