package main

import (
	"log"
	"os"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/generator"
)

func main() {
	resp := &plugin.CodeGeneratorResponse{}

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(buf, &req); err != nil {
		log.Fatalln(err)
	}

	g := generator.New(&req)
	g.Generate(resp)

	buf, err = proto.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}
	os.Stdout.Write(buf)
}
