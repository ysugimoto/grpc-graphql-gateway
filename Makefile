.PHONY: command

CMD=protoc-gen-graphql-gateway

command:
	protoc -I $(shell brew --prefix protobuf)/include/google \
		-I include/graphql \
		--go_out=./${CMD}/graphql \
		include/graphql/graphql.proto
	cd ${CMD} && go build -o ../dist/${CMD}
