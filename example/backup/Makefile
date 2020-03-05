.PHONY: clean all

UNAME := $(shell uname)

ifeq ($(UNAME), Darwin)
	LINKCMD := ln -s $(shell brew --prefix protobuf)/include/google ./proto/google
endif
ifeq ($(UNAME), Linux)
	LINKCMD := ln -s /usr/local/include/protobuf/google ./proto/google
endif

all: build

build: init
	rm -rf ./app/book
	rm -rf ./app/author
	for file in `find ./proto -name '*.proto'`; do \
		echo $$file; \
		protoc \
			-Iproto \
			--plugin=../../dist/protoc-gen-graphql \
			--graphql_out=./app \
			--go_out=plugins=grpc:./app \
			$$file; \
	done
	mv app/github.com/ysugimoto/grpc-graphql-gateway/example/backup/app/* ./app
	rm -rf app/github.com

init:
	if [ -L "./proto/google" ]; then rm ./proto/google; fi
	if [ -L "./proto/graphql" ]; then rm ./proto/graphql; fi
	$(LINKCMD)
	ln -s $(shell pwd)/../../include/graphql ./proto/graphql
	cd ../../ && make

