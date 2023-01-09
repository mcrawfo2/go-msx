# Traced Operations

The Tracing module has some convenience methods which allow you to create traced operations,
along with the ability to execute them in the foreground or background.

## Factories

The Tracing module defines two Operation factories which can create named and traced Operations.

### NewOperation

To create an Operation which has Distributed Tracing enabled, use the `NewOperation`
factory.  

For example:

```go
op := trace.NewOperation("multipleTwoNumbers", multiplyTwoNumbers)
```

This will create an Operation with two decorators:
- **SpanDecorator**: Records an operation and its outcome in the distributed trace.  If there 
  is a trace in progress, a new child span will be created inside the current span.
- **RecoverLogDecorator**: If your Action panics, this will stop propagation and log the details.

You can add further Middleware to the returned Operation, or consume it as-is.

### NewIsolatedOperation

The create an operation which has Distributed Tracing enabled, but is _not_ part of the 
current trace, use the `NewIsolatedOperation` factory.

For example:

```go
op := trace.NewIsolatedOperation("multipleTwoNumbers", multiplyTwoNumbers)
```

This uses `NewOperation` above, and then applies the following decorator:
- **UntracedContextDecorator**: Removes reference to the current span before a new span is created.
  This has the effect of starting a new trace, completely independent of the calling context.

## Execution

The Tracing module defines two Operation executors which can create and then execute
Traced Operations.

### ForegroundOperation

ForegroundOperation executes the action inside a new, isolated trace:

```go
err := trace.ForegroundOperation(ctx, "simple", mySimpleAction)
```

### BackgroundOperation

BackgroundOperation executes an action inside a background goroutine, using
a new, isolated trace:

```go
trace.BackgroundOperation(ctx, "simple", mySimpleAction)
```

This call does not offer persistence, cancellation, or restartability, so should not be used
for job execution or management.