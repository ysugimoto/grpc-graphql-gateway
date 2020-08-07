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
