# Development Targets

### `generate`

The `generate` target will execute any custom (or default) generate commands defined in
the `generate.*` entries.

If no command is specified for an entry, it will default to running `go generate` on that folder.

Generate commands can also be specified using `go:generate` [comments](https://blog.golang.org/generate).
Generation will be executed when `generate` executes on the directory containing files with these comments.

### `go-fmt`

The `go-fmt` target executes `go fmt` on directories which contain `*.go` files (excluding the `vendor` directory).

### `license`

