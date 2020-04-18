.PHONY: build

build:
	protoc \
	  -I. \
		--plugin=../../dist/protoc-gen-graphql \
	  --go_out=plugins=grpc:./greeter \
		--graphql_out=verbose:./greeter \
	  greeter.proto
