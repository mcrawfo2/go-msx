# Errors

Go has a built-in `error` interface to be implemented by error
models.  

go-msx has chosen to use the `github.com/pkg/errors`
module to implement errors.  This custom error module enables collecting
stack traces, critical for logging and debugging of errors.

When instantiating or wrapping an error,
use this package instead of the standard library `errors` package.

```go

import (
	"context"
	"github.com/pkg/errors"
)

// Create a globally visible error
var MyStaticError = errors.New("Static error occurred")
var MyOtherError = errors.New("Other error occurred")
	
// Return the global error
func mine(ctx context.Context) error {
  return MyStaticError
}

// Wrap the error into your own domain
func yours(ctx context.Context) error {
  return errors.Wrap(mine(ctx), "Something bad happened")
}

func callYours(ctx context.Context) error {
  err := yours(ctx)
  if errors.Is(err, MyStaticError) {
    // Special handling for this error type		
  } else {
    // General hanlding for any other error types
    return err
  }
  return nil
}
```

The above example shows how to create a global error, and how to re-contextualize (wrap)
inside the parent.

## Composition

Composite errors implement the `CompositeError` interface:

```go
type CompositeError interface {
  Errors() interface{}
}
```

go-msx provides two composite error models: `ErrorMap` and `ErrorList`.
Each of these represents a set of errors.

* `ErrorMap` : Represents a set of key-error pairs, intended to map to sub-parts
of a structured parent component.

    ```go
    return types.ErrorMap{
        "element1": validation.Validate(&element1, validation.Required)
        "element2": validation.Validate(&element2, validation.MinLength(1))
    }
    ```
    
* `ErrorList` : Represents a series of error instances (or nils), intended to map
to elements in a parent sequence.

    ```go
    return types.ErrorList{
        validation.Validate(&parent[0], validation.Required)
        validation.Validate(&parent[1], validation.MinLength(1))
    }
    ```

The above error models also implement `Filterable`:

```go
type Filterable interface {
  Filter() error
}
```

This allows the composite error to collect non-error (nil) values, which will
be removed from the return value of `Filter()`.  This feature is used by the 
validate package during DTO validation.

## Log Customization

To enable attaching custom log fields from your error, the logging
subsystem checks if your error implements the LogFielder interface:

```go
type LogFielder interface {
  LogFields() map[string]any
}
```

Any fields returned by the `LogFields()` function will be added
as log fields if the error is output to the log via `WithError()`.