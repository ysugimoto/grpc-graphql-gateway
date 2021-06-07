# grpc-graphql-gateway

[![CircleCI](https://circleci.com/gh/ysugimoto/grpc-graphql-gateway/tree/master.svg?style=svg)](https://circleci.com/gh/ysugimoto/grpc-graphql-gateway/tree/master)

`grpc-graphql-gateway` is a protoc plugin that generates graphql execution code from Protocol Buffers.

![image](https://raw.githubusercontent.com/ysugimoto/grpc-graphql-gateway/master/misc/grpc-graphql-gateway.png)

## Motivation

On API development, frequently we choose some IDL, in order to manage API definitions from a file.
Considering two of IDL -- GraphQL and Protocol Buffers (for gRPC) -- these have positive point respectively:

- **GraphQL** -- Can put together multiple resources getting into one HTTP request, appropriate for BFF
- **gRPC** -- Easy syntax in Protocol Buffers, and easy to implement API server using HTTP/2

But sometimes it's hard to maintain both GraphQL and Protocol Buffers, so we created this plugin in order to generate GraphQL Schema from Protocol Buffers.

This project much refers to [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) how to generate a file, provide a plugin. many thanks!

## Installation

### Get plugin binary

Get `protoc-gen-graphql` binary from [releases](https://github.com/ysugimoto/grpc-graphql-gateway/releases) page and set $PATH to be executable.

Or, simply get the latest one:

```shell
go get github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
```

Then the binary will place in `$GOBIN`.

### Get protobuf file

Get `include/graphql.proto` from this repository and put it into your project under the protobuf files.

```shell
git submodule add https://github.com/ysugimoto/grpc-graphql-gateway.git grpc-graphql-gateway
## Or another way...
```

## Usage

Please replace  `[your/project]` section to your appropriate project.

### Write Protocol Buffers

Declare gRPC service with protobuf using `grpc-graphql-gateway` options.
This example has two RPCs that names `SayHello` and `SayGoodbye`:

```protobuf
// greeter.proto
syntax = "proto3";

import "graphql.proto";

service Greeter {
  // gRPC service information
  option (graphql.service) = {
    host: "localhost:50051"
    insecure: true
  };

  rpc SayHello (HelloRequest) returns (HelloReply) {
    // Here is plugin definition
    option (graphql.schema) = {
      type: QUERY   // declare as Query
      name: "hello" // query name
    };
  }

  rpc SayGoodbye (GoodbyeRequest) returns (GoodbyeReply) {
    // Here is plugin definition
    option (graphql.schema) = {
      type: QUERY     // declare as Query
      name: "goodbye" // query name
    };
  }
}

message HelloRequest {
  // Below line means the "name" field is required in GraphQL argument
  string name = 1 [(graphql.field) = {required: true}];
}

message HelloReply {
  string message = 1;
}

message GoodbyeRequest {
  // Below line means the "name" field is required in GraphQL argument
  string name = 1 [(graphql.field) = {required: true}];
}

message GoodbyeReply {
  string message = 1;
}
```

### Compile to Go code

Compile protobuf file with the plugin:

```shell
protoc \
  -I. \
  --go_out=plugins=grpc:./greeter \
  --graphql_out=./greeter \
  greeter.proto
```

Then you can see `greeter/greeter.pb.go` and `greeter/greeter.graphql.go`.

### Implement service

For example, gRPC service will be:

```go
// service/main.go
package main

import (
    "context"
    "fmt"
    "net"
    "log"

    "github.com/[your/project]/greeter"
    "google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	return &greeter.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *Server) SayGoodbye(ctx context.Context, req *greeter.GoodbyeRequest) (*greeter.GoodbyeReply, error) {
	return &greeter.GoodbyeReply{
		Message: fmt.Sprintf("Good-bye, %s!", req.GetName()),
	}, nil
}

func main() {
	conn, err := net.Listen("tcp", ":500５1")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := grpc.NewServer()
	greeter.RegisterGreeterServer(server, &Server{})
	server.Serve(conn)
}
```

Then let's start service:

```shell
go run service/main.go
```

The gRPC service will start on `localhost:50051`.

Next, GraphQL gateway service should be:

```go
// gateway/main.go
package main

import (
    "log"

    "net/http"

    "github.com/[your/project]/greeter"
    "github.com/ysugimoto/grpc-graphql-gateway/runtime"
)

func main() {
    mux := runtime.NewServeMux()

    if err := greeter.RegisterGreeterGraphql(mux); err != nil {
        log.Fatalln(err)
    }
    http.Handle("/graphql", mux)
    log.Fatalln(http.ListenAndServe(":8888", nil))
}
```

Then let's start gateway:

```shell
go run gateway/main.go
```

The GraphQL gateway will start on `localhost:8888`

### Send request via the gateway

Now you can access gRPC service via GraphQL gateway!

```shell
curl -g "http://localhost:8888/graphql" -d '
{
  hello(name: "GraphQL Gateway") {
    message
  }
}'
#=> {"data":{"hello":{"message":"Hello, GraphQL Gateway!"}}}
```

You can also send request via `POST` method with operation name like:

```shell
curl -XPOST "http://localhost:8888/graphql" -d '
query greeting($name: String = "GraphQL Gateway") {
  hello(name: $name) {
    message
  }
  goodbye(name: $name) {
    message
  }
}'
#=> {"data":{"goodbye":{"message":"Good-bye, GraphQL Gateway!"},"hello":{"message":"Hello, GraphQL Gateway!"}}}
```

This is the most simplest way :-) 

## Resources

To learn more, please see the following resources:

- `graphql.proto` Plugin option definition. See a comment section for custom usage (e.g mutation).
- [example/greeter](https://github.com/ysugimoto/grpc-graphql-gateway/tree/master/example/greeter)  Files of above usage.
- [example/starwars](https://github.com/ysugimoto/grpc-graphql-gateway/tree/master/example/starwars) Common implementation for GraphQL explanation, the StarWars API example

This plugin generates graphql execution code using [graphql-go/graphql](https://github.com/graphql-go/graphql), see that repository in detail.

## Limitations

This plugin just aims to generate a simple gateway of gRPC.

Some of things could be solved and could not be solved.
The most of limitations come from the IDL's power of expression -- some kind of GraphQL schema feature cannot implement by Protocol Buffers X(

Currently we don't support some Protobuf types:

- Builtin `oneof` type

## Contribute

- Fork this repository
- Customize / Fix problem
- Send PR :-)
- Or feel free to create issue for us. We'll look into it

## Author

Yoshiaki Sugimoto

## License

MIT
