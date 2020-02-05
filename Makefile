.PHONY: command

CMD=protoc-gen-graphql

command: plugin
	cd ${CMD} && go build -o ../dist/${CMD}

plugin:
	protoc -I $(shell brew --prefix protobuf)/include/google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/ysugimoto/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

all: plugin
	cd ${CMD} && GOOS=darwin GOARCH=amd64 go build -o ../dist/${CMD}.darwin
	cd ${CMD} && GOOS=linux GOARCH=amd64 go build -o ../dist/${CMD}.linux

