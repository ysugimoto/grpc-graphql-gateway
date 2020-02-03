package main

import (
	"net/http"

	"github.com/ysugimoto/grpc-graphql-gateway/middleware"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	_ "google.golang.org/grpc"
)

func main() {
	// conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()
	mux := runtime.NewServeMux(middleware.Cors())

	RegisterGraphqlHandler(mux, nil, "/graphql")
	http.ListenAndServe(":8888", mux)
}
