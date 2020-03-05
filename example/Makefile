.PHONY: clean all

UNAME := $(shell uname)
SPECDIR := "./spec"

ifeq ($(UNAME), Darwin)
	LINKCMD := ln -s $(shell brew --prefix protobuf)/include/google ./proto/google
endif
ifeq ($(UNAME), Linux)
	LINKCMD := ln -s /usr/local/include/protobuf/google ./proto/google
endif

all: build

build: init
	rm -rf ./spec/starwars
	for file in `find ./proto -name '*.proto'`; do \
		echo $$file; \
		protoc \
			-Iproto \
			--plugin=../dist/protoc-gen-graphql \
			--graphql_out=${SPECDIR} \
			--go_out=plugins=grpc:${SPECDIR} \
			$$file; \
	done
	mv ${SPECDIR}/github.com/ysugimoto/grpc-graphql-gateway/example/spec/* ${SPECDIR}
	rm -rf ${SPECDIR}/github.com

init:
	if [ -L "./proto/google" ]; then rm ./proto/google; fi
	if [ -L "./proto/graphql" ]; then rm ./proto/graphql; fi
	$(LINKCMD)
	ln -s $(shell pwd)/../include/graphql ./proto/graphql
	cd ../ && make

