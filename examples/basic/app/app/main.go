package main

import (
	"net/http"

	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	mux := runtime.NewServeMux()

	RegisterGraphqlHandler(mux, conn, "/graphql")
	http.ListenAndServe(":8888", mux)
}
