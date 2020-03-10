package main

import (
	"log"
	"net/http"

	"github.com/ysugimoto/grpc-graphql-gateway/example/starwars/spec/starwars"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
)

func main() {
	mux := runtime.NewServeMux(runtime.Cors())

	starwars.RegisterStartwarsServiceGraphqlHandler(mux, nil)
	http.Handle("/graphql", mux)
	log.Fatalln(http.ListenAndServe(":8888", nil))
}
