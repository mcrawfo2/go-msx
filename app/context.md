# MSX Dependencies

In large applications, inter-object dependency management becomes more challenging.  Within the go standard library, the [`Context`](https://blog.golang.org/context) object is provided to share dependencies and cancellation.  This simplifies writing unit tests, since dependencies can be injected via the context.

MSX Application provides a `Context` object to event `Observer`s so they may inject new dependencies for their subsystems.  The context object also carries Trace information for logging and trace publishing.

By default, the following dependencies are added to the MSX Application context:
- Configuration
- Cassandra client pool
- Cockroach client pool
- Consul client pool
- Vault client pool
- Redis client pool
- Kafka client pool
- Http client factory

During `migrate` execution, the Migration Manifest is also available from the context.

## Accessing Dependencies

To support adding a custom dependency to any context, define the standard context chaining pattern:

```go
type DomainServiceApi interface {
}

type domainContextKey int

const contextKeyDomainService domainContextKey = iota

func ContextWithDomainService(ctx context.Context, domainService DomainServiceApi) context.Context {
	return context.WithValue(ctx, contextKeyDomainService, domainService)
}

func DomainServiceFromContext(ctx context.Context) DomainServiceApi {
    return ctx.Value(contextKeyDomainService).(DomainServiceApi)
}
```

To inject your custom dependency to the current context:

```go
ctx := domain.ContextWithDomainService(ctx, domainService)
```

To retrieve a dependency from the current context:

```go
domainServiceApi := domain.DomainServiceFromContext(ctx)
```

## Logging and Tracing

To apply logging and tracing fields from the current context:

```go
myLogger.WithContext(ctx).Info("My log message")
```
