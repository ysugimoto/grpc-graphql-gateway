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
			--plugin=../..//dist/protoc-gen-graphql \
			--graphql_out=${SPECDIR} \
			--graphql_opt=paths=source_relative \
			--go_out=:${SPECDIR} \
			--go_opt=paths=source_relative \
			--go-grpc_out=:${SPECDIR}  \
			--go-grpc_opt=paths=source_relative \
			$$file; \
	done

init:
	if [ -L "./proto/google" ]; then rm ./proto/google; fi
	if [ -L "./proto/graphql" ]; then rm ./proto/graphql; fi
	$(LINKCMD)
	ln -s $(shell pwd)/../..//include/graphql ./proto/graphql
	cd ../../ && make

start: build
	cd app/graphql && GOOS=linux go build -o app main.go
	cd app/grpc && GOOS=linux go build -o app main.go
	docker-compose up -d --build

stop:
	docker-compose down

start_graphql:
	go run app/graphql/main.go

start_grpc:
	go run app/grpc/main.go