# CHANGELOG

## v0.16.0

## New graphql.proto option

- Add new graphql option of `resolver` [#24](https://github.com/ysugimoto/grpc-graphql-gateway/pull/24)

```protobuf
message GraphqlField {
    bool required = 1;
    string name = 2;
    string default = 3;
    bool omit = 4;
    string resolver = 5;
}
```

A new field of `resolve` which resolves as a nested query.

## v0.15.0

Add `omit` option in graphql.field option.

## v0.14.6

Bugfixes

- enum: use protobuf enum type for value [#18](https://github.com/ysugimoto/grpc-graphql-gateway/pull/18) (@bmkessler)

## v0.14.5

Bugfixes

- Implement request transformation from CamelCase to SnakeCase.

## v0.14.3, v0.14.4

Bugfixes

- Fix unexpected panic caused by reflect package in marshaling JSON response with camel-case

## v0.14.1, v0.14.2

Bugfixes

- Fix tiny syntax error
- Fix field camelcase generation
- Fix required sign in repeated array and input object

## v0.14.0

### Partially support google's type

Implement specific types and provide from this repository:

- google.protobuf.Timestamp
- google.protobuf.Wrappers
- google.protobuf.Empty

Note that we could only support the above types hence if a user imports other types e.g. google.protobuf.Any will raise an error.

And for GraphQL spec, google.protobuf.Empty defined empty field like `_: Boolean` because GraphQL raises error when type fields are empty.

### Fix map type in proto3

In proto3, map type can define as `map<string, string>` but PB parse as xxxEntry message. This PR also can use as graphql type with required key and value.

### Version printing in plugin binary

`protoc-gen-graphql` accepts `-v` option for print build version.

## v0.13.0

### Add infix typename

Add infix typename to GraphQL typename in order to avoid conflicting name between `Type` and `Input`.
After this version, GraphQL typename is modified that `[PackageName]_[Type]_[MessageName]` for example:

```
package user

message User {
  int64 user_id = 1;
}
```

Then typename is `User_Type_User`.

### Convert field name to CamelCase option

In protocol buffers, all message field name should define as *lower_snake_case* referred by [Style Guide](https://developers.google.com/protocol-buffers/docs/style#message_and_field_names). But in GraphQL Schema, typically each field name should define as *lowerCamelCase*, so we add more option in `protoc-gen-graphql`:

```shell
protoc -I.
    --graphql_out=field_camel=true:.
    --go_out=plugins=grpc:.
    example.proto
```

The new option, `field_camel=true` converts all message field name to camel case like:

```
// protobuf
message User {
    int64 user_id = 1 [(graphql.field).required = true];
    string user_name = 2 [(graphql.field).required = true];
}

// Graphql Schema
type User_Type_User {
    userId Int!
    userName String!
}
```

To keep backward compatibility, compile as *lower_snake_case* as default. If you want to define any graphql field as lowerCamelCase, please supply this option.

## v0.12.0

### Define MiddlewareError and respond with error code

We defined `MiddleWareError` which can return in Middleware function. If middleware function returns that pointer instead of common error,
The runtime responds error with that field data.

The definition is:

```go
type MiddlewareError struct {
  Code string
  Message string
}

// generate error
return runtime.NewMiddlewareError("CODE", "MESSAGE")
```

If middleware returns common error, the runtime error will be:

```
{"data": null, "errors": [{"message": "error message of err.Error()", "extensions": {"code": "MIDDLEWARE_ERROR"}]}
```

If middleware returns `MiddlewareError`, the runtime error will be:

```
{"data": null, "errors": [{"message": "Message field value of MiddlewareError", "extensions": {"code": "Code field value of MiddlewareError"}]}
```


## v0.9.1

### Changed middleware fucntion type

On MiddlewareFunc, you need to return `context.Context` as first return value. this is because we need to make custom metadata to gRPC on middleware process.
If you are already using your onw middleware, plase change its interface. see https://github.com/ysugimoto/grpc-graphql-gateway/pull/10
