# Retry
MSX Retry enables attempting an operation multiple times, stopping
on success (no error returned) or permanent operation failure.
## Configuration Examples
Retries without delays
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 2,
    Delay:    0,
    BackOff:  0.0,
    Linear:   true,
})
```
Retries with fixed delays (1 second)
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 2,
    Delay:    1000,
    BackOff:  1.0,
    Linear:   true,
})
```
Retries with linear delays (1, 2, 3, 4)
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 5,
    Delay:    1000,
    BackOff:  1.0,
    Linear:   true,
})
```
Retries with exponential delays (1, 2, 4, 8)
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 5,
    Delay: 1000,
    BackOff: 2.0,
    Linear: false,
})
```

Retries with linear delay and Jitter (extreme random) (1, 7.8, 20.3, 8.45) <br />
> With higher Jitter value you could expect greater randomness
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 5,
    Delay:    1000,
    Jitter:   20000,
    Linear:   true,
})
```
Retries with exponential delay and Jitter (1, 2, 4, 8) (note: jitter is negligible so this is just like exponential backoff with no jitter)
```go
r := NewRetry(ctx, RetryConfig{
    Attempts: 5,
    Delay: 1000,
    BackOff: 2.0,
    Linear: false,
    Jitter: 1,
})
```


Using retry with decorator
```go
	types.
		NewOperation(func(ctx context.Context) error {
			return errors.New("a transient error")
		}).
		WithDecorator(Decorator(RetryConfig{
			Attempts: 1,
			Delay:    10,
			BackOff:  2.0,
			Linear:   false,
			Jitter:   1,
		})).
		Run(ctx)
```
