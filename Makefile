.PHONY: command

CMD=protoc-gen-graphql

command:
	protoc -I $(shell brew --prefix protobuf)/include/google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/ysugimoto/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com
	cd ${CMD} && go build -o ../dist/${CMD}
