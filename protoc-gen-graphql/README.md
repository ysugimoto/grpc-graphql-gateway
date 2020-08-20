# protoc-gen-graphql

`protoc` plugin binary.

## Usage

Download binary which corresponds to runtime version via [Release Page](https://github.com/ysugimoto/grpc-graphql-gateway/releases) and place in `$PATH` directory.

Compile with this plugin example:

```shell
protoc \
    -I. \
    --graphql_out=[arguments]:[output_dir] \
    --go_out=plugins=grpc:[output_dir] \
    example.proto
```

## Compilation Arguments

This plugin accepts some compile arguments:

- `--graphql_out=verbose`: verbose debug output
- `--graphql_out=exclude=[regex]`: exclude generation package with regexp
- `--graphql_out=field_camel`: all graphql field name transform to lower-camel-case

All arguments can be provide by splitting comma.

## Binary Option

Typically `protoc-gen-graphql` will process via stdin but some option can accept for development:

- `protoc-gen-graphql -v`: Print plugin version
