package main

import (
	"log"
	"os"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/generator"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

func main() {
	var genError error

	resp := &plugin.CodeGeneratorResponse{}
	defer func() {
		// If some error has been occurred in generate process,
		// add error message to plugin response
		if genError != nil {
			message := genError.Error()
			resp.Error = &message
		}
		buf, err := proto.Marshal(resp)
		if err != nil {
			log.Fatalln(err)
		}
		os.Stdout.Write(buf)
	}()

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		genError = err
		return
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(buf, &req); err != nil {
		genError = err
		return
	}

	var args *spec.Params
	if req.Parameter != nil {
		args, err = spec.NewParams(req.GetParameter())
		if err != nil {
			genError = err
			return
		}
	}

	// We're dealing with each descriptors to out wrapper struct
	// in order to access easily plugin options, pakcage name, comment, etc...
	var files []*spec.File
	for _, f := range req.GetProtoFile() {
		files = append(files, spec.NewFile(f))
	}

	g := generator.New()
	genFiles, err := g.Generate(files, args)
	if err != nil {
		genError = err
		return
	}
	resp.File = append(resp.File, genFiles...)
}
