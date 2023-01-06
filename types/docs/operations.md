# Operations

In the following sections, we will explore mechanisms for reusing and augmenting
code using Functions, Actions, Operations, Decorators, Filters, and Handlers.

## Functions

In go, functions are the lowest level of reusable execution. They accept a specific
set of arguments and return a set of values.

```go
func add(number1, number2 int) int {
  return number1 + number2
}
```

## Function Types

Go functions are first-class, meaning they can be passed around as values, including as
parameters to other functions, or return values from functions. This enables powerful code
composition and reuse. This is also a key feature of the functional style of programming:

> In functional programming, functions are treated as first-class citizens,
> meaning that they can be bound to names (including local identifiers), passed
> as arguments, and returned from other functions, just as any other data type can.
> This allows programs to be written in a declarative and composable style, where
> small functions are combined in a modular manner.
>
> -- [Functional Programming](https://en.wikipedia.org/wiki/Functional_programming), [Wikipedia](https://en.wikipedia.org)

In go, to accept a function as an argument, you can declare a function type, and use it to declare the receiving
parameter:

```go
type unaryOperator func (int) int
type binaryOperator func(int, int) int

func evaluateUnaryExpression(operand int, operator unaryOperator) int {
  return unaryOperator(operand)
}

func evaluateBinaryExpression(leftOperand, rightOperand int, operator binaryOperator) int {
  return binaryOperator(leftOperatnd, rightOperand)
}
```

## Actions

go-msx defines an ActionFunc type to describe an executable function (Action) signature:

```go
type ActionFunc func(ctx context.Context) error
```

An `ActionFunc` accepts a single `Context` argument (to allow access to dependencies and operation-scoped data),
and returns a single `error` value indicating success (`nil`) or failure (non-`nil`).  As described
above, this enables you to pass around these functions and abstractly re-use them:

```go
// Send a message to the ANSWER_TOPIC channel
func deepThought(ctx context.Context) error {
  return stream.PublishObject('ANSWER_TOPIC', map[string]any{
    "answer": 42, 		
  })
}

// Call the deepThought function when the application is running
func init() {
  app.OnEvent(app.EventRun, app.PhaseDuring, deepThought)
}
```

In this example, we register an application event observer Action to be executed
when the application has finished startup.  The Action sends a simple message to
a stream.

## Operations

To simplify reusing code to work with Actions, go-msx has an Operation type:

```go
type Operation struct {...}
func (o Operation) Run(ctx context.Context) error {...}
```

Operations provide a `Run` method to execute the operation, along with other methods
to create derived Operations using Filters and Decorators.  These will be discussed in the 
next section.
