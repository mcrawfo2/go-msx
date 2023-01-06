# Filters

When reusing multiple Middleware types for a given Operation, it may be important that some 
middleware instances are consistently applied before others.

For example, our TokenFilter is applied to our HTTP pipeline to extract the token from Request
and inject it into the Context.  This must occur _before_ our AuthenticationFilter checks the
token in the Context and verifies the user is properly authenticated.   Since we can swap in 
CertificateFilter for TokenFilter (when using Certificate Authentication), it important that 
the Middleware not be coupled, and always be applied in the correct order.

Decorators do not directly allow application in middleware-specified order.  This presents many
problems when using factories to generate higher level Operation abstractions such as Endpoints
and Message Subscribers.  These factories do not need to know about the variety or application 
order of these Decorators, especially when they are mixed with framework-specified Middleware.

To enable this scenario, go-msx offers Filters.  Filters allow you to ensure the correct
ordering of Middleware when passed along from other components without requiring tight
coupling or specialization.

## Usage

To apply a filter to an Operation, use the `WithFilter` method:
```go
err := types.NewOperation(myAction).
	WithFilter(NewLoggingFilter(logger)).
	Run(ctx)
```

### Order

Filter Order can be though of as priority: a higher number means it will be applied earlier
to the target Action.

For example, if Filter A has an order of `0` and Filter B has an order of `100`, then
Filter B will be applied first (executed second inbound, first outbound), and Filter A will
be applied second (executed first inbound, second outbound).

Note that when combining Operation instances using the `Operation.Run` method, filters are
only ordered relative to the other Filters and Decorators on the Operation to which it was directly
applied.

## Implementation

A Filter is envisioned as a simple wrapper around a Decorator which also
provides a method to inspect the order that it should be applied:

```go
// ActionFilter is an ordered Decorator
type ActionFilter interface {
  Order() int
  Decorator() ActionFuncDecorator
}
```

When authoring a Filter, any type that implements `ActionFilter` can be used.
You can even transform a Decorator into a Filter using a simple factory:

```go
recoveryDecorator := NewRecoveryDecorator()
filter := types.NewOrderedDecorator(100, recoveryDecorator)
```
