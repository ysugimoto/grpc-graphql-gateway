# CHANGELOG

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
