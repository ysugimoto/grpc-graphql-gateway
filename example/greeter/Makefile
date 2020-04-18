.PHONY: build

build:
	cd ../../ && make
	protoc \
	  -I. \
		-I./grpc-gateway \
		--plugin=../../dist/protoc-gen-graphql \
	  --go_out=plugins=grpc:./gateway \
		--graphql_out=verbose:./gateway \
	  ./grpc-gateway/google/api/annotations.proto
	#protoc \
	#  -I. \
	#	--plugin=../../dist/protoc-gen-graphql \
	#  --go_out=plugins=grpc:./greeter \
	#	--graphql_out=verbose:./greeter \
	#  greeter.proto
