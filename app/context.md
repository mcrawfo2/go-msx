# MSX Dependencies

In large applications, inter-object dependency management becomes more challenging.  Within the go standard library, the [`Context`](https://blog.golang.org/context) object is provided to share dependencies and cancellation.  This simplifies writing unit tests, since dependencies can be injected via the context.

MSX Application provides a `Context` object to event `Observer`s so they may inject new dependencies for their subsystems.  The context object also carries Trace information for logging and trace publishing.

By default, the following dependencies are added to the MSX Application context:
- Configuration
- Cockroach client pool
- Consul client pool
- Vault client pool
- Redis client pool
- Kafka client pool
- Http client factory

During `migrate` execution, the Migration Manifest is also available from the context.

## Accessing Dependencies

Each substitutable component in go-msx requires a context accessor to allow injecting and
inspecting overrides:

```go
type contextKeyNamed string

func ContextDomainService() types.ContextAccessor[DomainService] {
	return types.NewContextAccessor[DomainService](contextKeyNamed("DomainService"))
}
```

Key to type safety is the external invisibility of the context key.  This is guaranteed
by defined a module-local type (`contextKey` or `contextKeyNamed`) and using an instance of it
to index the context inspection/injection.

To inject your custom dependency to the current context:

```go
ctx = domain.ContextDomainService().Set(ctx, domainService)
```

To retrieve a dependency from the current context:
```go
domainServiceApi := domain.ContextDomainService().Get(ctx)
```
or:
```go
domainServiceApi, ok := domain.ContextDomainService().TryGet(ctx)
```

## Logging and Tracing

To apply logging and tracing fields from the current context:

```go
myLogger.WithContext(ctx).Info("My log message")
```
