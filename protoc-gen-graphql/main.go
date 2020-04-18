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
	if err = proto.Unmarshal(buf, &req); err != nil {
		genError = err
		return
	}

	var parameter string
	if req.Parameter != nil {
		parameter = req.GetParameter()
	}
	args, err := spec.NewParams(parameter)
	if err != nil {
		genError = err
		return
	}

	// We're dealing with each descriptors to out wrapper struct
	// in order to access easily plugin options, package name, comment, etc...
	var files []*spec.File
	for _, f := range req.GetProtoFile() {
		files = append(files, spec.NewFile(f))
	}

	g := generator.New(files, args)
	var ftg []string
	for _, f := range req.GetFileToGenerate() {
		if !args.IsExclude(f) {
			ftg = append(ftg, f)
		}
	}
	if len(ftg) > 0 {
		genFiles, err := g.Generate(goTemplate, ftg)
		if err != nil {
			genError = err
			return
		}
		resp.File = append(resp.File, genFiles...)
	}
}
