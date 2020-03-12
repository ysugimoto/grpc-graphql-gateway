# StarWars API example

A bit actual usage, define as Go package, using `google.protobuf.XXX` package, run with docker container.

## How to try

This example uses google's protobuf package, so you need to be able to import them:

```shell
# mac
brew install protobuf
# linux
follow specific package manager...
```

Build protoc to Go code, and start docker-compose:

```shell
make start
```

Then service and gateway containers are running:

- service - `localhost:50051`
- gateway - `localhost:8888`

Try to send GraphQL request:

```shell
curl -g "http://localhost:8888/graphql" -d '
{
  humans {
    id
    name
    friends {
      name
    }
  }
  hero(episode: JEDI) {
    name
  }
}'
```
