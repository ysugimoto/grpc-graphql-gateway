.PHONY: command clean

GRAPHQL_CMD=protoc-gen-graphql
SCHEMA_CMD=protoc-gen-graphql-schema

command: plugin clean
	cd ${GRAPHQL_CMD} && go build -o ../dist/${GRAPHQL_CMD}
	cd ${SCHEMA_CMD} && go build -o ../dist/${SCHEMA_CMD}

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
	cd ${SCHEMA_CMD} && GOOS=darwin GOARCH=amd64 go build -o ../dist/${SCHEMA_CMD}.darwin
	cd ${SCHEMA_CMD} && GOOS=linux GOARCH=amd64 go build -o ../dist/${SCHEMA_CMD}.linux
