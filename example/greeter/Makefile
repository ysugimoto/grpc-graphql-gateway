.PHONY: build

build:
	protoc \
	  -I. \
		-I../../include/graphql \
		--plugin=../../dist/protoc-gen-graphql \
	  	--go_out=:./greeter \
		--go-grpc_out=./greeter \
		--graphql_out=verbose:./greeter \
	  greeter.proto

start_grpc:
	go run service/main.go

start_graphql:
	go run gateway/main.go