# Retry

Retry enables attempting an operation multiple times, stopping
on success (no error returned) or permanent operation failure.

## Retrying

To retry an action, create a new instance of `Retry` via `NewRetry` and
call its `Retry` method:

```go

tooEarly := time.Parse("2020-01-01T00:00:00Z")
tooLate := time.Parse("2020-12-31T23:59:59.999999999Z")

// Retry once an hour
r := NewRetry(ctx, RetryConfig{Attempts:10, Delay: 60 * 60 * 1000})

err := r.Retry(func() error {
  now := time.Now()
  if now.Before(tooEarly) {
    return retry.TransientError{
      Cause: errors.New("Will succeed in the future")
    } 		
  } else if now.After(tooLate) {
    return retry.PermanentError{
      Cause: errors.New("Will never succeed again")			
    }
  }
  return nil
})
```

The above retries its action once per hour, with up to 10 attempts.
If the time is before `tooEarly`, it will continue retrying, since it
returns a `TransientError`.  If the time is after `tooLate`, it will
stop retrying, since it returns `PermanentError`.  If the time is after `tooEarly`
but before `tooLate`, it will succeed and cease further attempts.

Retry distinguishes between Transient and Permanent errors by inspecting
the returned `error` instance.  If it implements the failure interface,
it can be queried for transience/permanence:

```go
type failure interface {
	IsPermanent() bool
}
```

Permanent errors should return `true` from `IsPermanent()`, transient
errors should return `false`.  As above, this can be handled by wrapping
the error in either `PermanentError` or `TransientError`.

## Configuration Examples

- Retries without delays
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 2,
        Delay:    0,
        BackOff:  0.0,
        Linear:   true,
    })
    ```

- Retries with fixed delays (1 second)
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 2,
        Delay:    1000,
        BackOff:  1.0,
        Linear:   true,
    })
    ```

- Retries with linear delays (1, 2, 3, 4)
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 5,
        Delay:    1000,
        BackOff:  1.0,
        Linear:   true,
    })
    ```

- Retries with exponential delays (1, 2, 4, 8)
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 5,
        Delay: 1000,
        BackOff: 2.0,
        Linear: false,
    })
    ```

- Retries with linear delay and Jitter (low random) (1, 2.452, 3.571, 4.357) <br />
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 5,
        Delay:    1000,
        BackOff:  1.0,
        Linear:   true,
        Jitter:   1000,
    })
    ```

- Retries with linear delay and Jitter (extreme random) (1, 7.8, 20.3, 8.45) <br />
  With higher Jitter value you could expect greater randomness.
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 5,
        Delay:    1000,
        BackOff:  1.0,
        Linear:   true,
        Jitter:   20000,
    })
    ```

- Retries with exponential delay and Jitter (1, 2, 4, 8) (note: jitter is negligible so this is just like exponential backoff with no jitter)
    ```go
    r := NewRetry(ctx, RetryConfig{
        Attempts: 5,
        Delay: 1000,
        BackOff: 2.0,
        Linear: false,
        Jitter: 1,
    })
    ```

- Using retry with decorator
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
