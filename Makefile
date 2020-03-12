.PHONY: command clean

GRAPHQL_CMD=protoc-gen-graphql

command: plugin clean
	cd ${GRAPHQL_CMD} && go build -o ../dist/${GRAPHQL_CMD}

lint:
	golangci-lint run

plugin:
	protoc -I $(shell brew --prefix protobuf)/include/google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/ysugimoto/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

build:
	protoc -I google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/ysugimoto/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

clean:
	rm -rf ./dist/*

all: clean build
	cd ${GRAPHQL_CMD} && GOOS=darwin GOARCH=amd64 go build -o ../dist/${GRAPHQL_CMD}.darwin
	cd ${GRAPHQL_CMD} && GOOS=linux GOARCH=amd64 go build -o ../dist/${GRAPHQL_CMD}.linux
