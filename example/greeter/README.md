# Greeter example

Most simplest way to try this project example.

## How to try

Build protoc to Go code, and start serivce and gateway:

```shell
make
go run service/main.go
go run gateway/main.go
```

Then services are listening:

- service - `localhost:50051`
- gateway - `localhost:8888`

Try to send GraphQL request:

```shell
curl -g "http://localhost:8888/graphql" -d '
{
  hello(name: "GraphQL Gateway") {
    message
  }
}'
```
