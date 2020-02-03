.PHONY: clean graphql

graphql: build
	find ./proto -name "*.proto" | xargs protoc \
		-Iproto \
		--plugin=../../dist/protoc-gen-graphql \
		--graphql_out=query=./,go=./app:./

build: clean init
	cd ../../ && make
	for file in `find ./proto -name '*.proto'`; do \
		protoc \
			-Iproto \
			--go_out=plugins=grpc:./app \
			$$file; \
	done
	mv app/github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/* ./app
	rm -rf app/github.com

clean:
	rm -rf ./app/book
	rm -rf ./app/author

init:
	if [ -L "./proto/google" ]; then rm ./proto/google; fi
	if [ -L "./proto/graphql" ]; then rm ./proto/graphql; fi
	ln -s $(shell brew --prefix protobuf)/include/google ./proto/google
	ln -s $(shell pwd)/../../include/graphql ./proto/graphql

