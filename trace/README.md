# MSX Distributed Tracing

MSX Distributed Tracing allows the collection of an operational flow graph.  Based on [OpenTracing](https://opentracing.io/docs/overview/), tracing helps pinpoint where failures occur and what causes poor performance.

## Model

- **Span**

  A span is a named, timed operation representing a piece of the operational flow.  Spans can have parents and children.

- **Trace**

  A trace is the complete tree of spans from an entire operational flow.  A new trace (with a new root span) is created by input from an external system, such as a REST API client.  Traces extend across synchronous and asynchronous message flows (interal RPC and events).

## Usage

The most common usage of tracing is to create a new child span within the current span, and execute an operation inside it.  To facilitate this, you can use the `trace.Operation()` function:

```go
err := trace.Operation(ctx, "myChildOperation", func(ctx context.Context) error { 
    myLogger.WithContext(ctx).Info("Inside myChildOperation...")
    return nil
})
```

To create a new child span and attach data to it, you can use the `trace.NewSpan()` function:

```go
// Create the new span
ctx, span := trace.NewSpan(ctx, spanName)
defer span.Finish()

// Tag the operation name
span.SetTag(trace.FieldOperation, operationName)

// Execute the operation and record the result
if err := myOperation(); err != nil {
    span.LogFields(trace.Status("ERROR"), trace.Error(err))
} else {
    span.LogFields(trace.Status("OK"))
}
```

Common trace log tags include:
- `trace.FieldOperation`: Generic operation name
- `trace.FieldStatus`: Terminal status of the operation
- `trace.FieldHttpCode`: Response status code
- `trace.FieldHttpUrl`: Request url
- `trace.FieldHttpMethod`: Request method

Other tags can be defined as needed using simple period-separated strings (e.g. `grpc.response.code`).

## Advanced Usage

When writing a new driver for external input (such as a new RPC transport listener), you can retrieve the untraced context:

```go
ctx = trace.UntracedContextFromContext(ctx)
```

This context object should be passed to the input handlers, who will be responsible for starting a new (root) span:

```go
err := trace.Operation(ctx, "myInputReceiver", myInputHandler)
```

## Configuration

By default, MSX tracing will send trace data to a Jaeger listener at `udp://localhost:6831`.

The following configuration settings can be specified to override the default behaviour:

| Key                   | Description | Default |
|-----------------------|-------------|---------|
| `trace.enabled`       | collect and forward distributed tracing data | `true` |
| `trace.service-name`  | name of service to supply with the trace | `${info.app.name}` |
| `trace.reporter.name` | which reporter to use, `jaeger` or `zipkin` | `jaeger` |
| `trace.reporter.host` | jaeger host | `localhost` |
| `trace.reporter.port` | jaeger udp port | `6831` |
| `trace.reporter.url`  | zipkin url | `http://localhost:9411/api/v1/spans` |

