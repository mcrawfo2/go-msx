# Development Targets

### `download-generate-deps`

The `download-generate-deps` target installs cross-project generation dependencies, including:

- github.com/rust-lang/mdBook
- github.com/badboy/mdbook-mermaid
- github.com/vektra/mockery/v2

### `generate`

The `generate` target will execute any custom (or default) generate commands defined in
the `generate.*` [entries](config.md#generate).

If no command is specified for an entry, it will default to running `go generate` on that folder.

Generate commands can also be specified using `go:generate` [comments](https://blog.golang.org/generate).
Generation will be executed when `generate` executes on the directory containing files with these comments.

### `go-fmt`

The `go-fmt` target executes `go fmt` on directories which contain `*.go` files (excluding the `vendor` directory).

### `license`

The `license` target verifies that all go source code files contain the appropriate Cisco license header. 

### `update-openapi-producer-spec`

The `update-openapi-producer-spec` target will obtain the latest version of the
microservice _producer_ OpenApi contract specification and overwrite the stored version.

Producer specification file is configured via the `openapi.spec` build [setting](config.md#openapi).

### `update-openapi-consumer-spec`

The `update-openapi-consumer-spec` target will obtain the latest version of the
_consumer_ OpenApi contract specification and overwrite the stored version.

Consumer local and remote specification are configured via the `openapi.contracts[*]` build [setting](config.md#openapi). 
