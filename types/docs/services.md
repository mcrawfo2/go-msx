# Services

A service is a reusable part of an application.  Services are used to isolate responsibility
and are composed to provide functionality. Example services include:
- REST Controller
- Application Service
- Stream Message Subscriber
- Database Repository
- API Integration

Each service may have a number of supporting components and functions:
- Interface Definition
- Mock
- Structure
- Dependencies
- Implementation
- Abstract Constructor
- Lifecycle Registration (for root components)

Let's look at each of these using an example Application Service, `HelloWorldService`. 

## Interface Definition

To enable test duplicate (mock) substitutions, you should define an interface declaring the
public methods of your component.  In our example, we have one method, `SayHello`:

```go
type HelloWorldService interface {
	SayHello(context.Context) (string, error)
}
```

This interface type will be used later by our Abstract Constructor to ensure we provide
any substituted dependency instead of returning a live object when requested during testing.

The interface should be externally visible (capitalized) so that other modules can re-use it.
This opposes the standard go convention of "interface definition by consumer" enabled by
duck typing, however it allows you to pre-generate mocks for your consumers' testing needs.

## Mock

Each service (other than root components) will be re-used by one ore more other services,
and therefore should provide a mock.  Using mockery, for example:

```go
//go:generate mockery --name=HelloWorldService --case=snake --with-expecter
```

This mock will be generated automatically when you run `go generate`, such as when using
the `make generate` target.

## Structure

Each service is defined using a simple go structure:

```go
type helloWorldService struct {} 
```

In most situations, you do not want to make the implementation visible outside the current module,
and therefore the structure name should start with a lowercase letter.  Consumers of your structure 
will receive a reference via the Interface, which will be visible externally.

## Dependencies

A service often depends on other services (provided by the go-msx framework, your application,
or third parties).  These dependencies are declared in the service structure.  For example,
our `HelloWorldService` can depend on a repository:

```go
type helloWorldService struct {
	helloWorldRepository HelloWorldRepository
}
```

Dependencies should be declared in the structure by referring to their abstract (interface)
type, so that during coverage testing, you can use Mocks to test all code paths.

Module dependencies should not be declared in the structure, but rather in the local `pkg.go`.
This includes loggers.

## Implementation

Each service will have a series of public functions matching the Interface Definition:

```go
func (r *helloWorldService) SayHello(ctx context.Context) (string, error) {
	return "Hello, World", nil
}
```

Service methods should use a pointer receiver, as they will be passed around on the heap
inside an interface reference.

## Context Accessor

Each substitutable component in go-msx requires a context accessor to allow injecting and
inspecting overrides:

```go
func ContextHelloWorldService() types.ContextAccessor[HelloWorldService] {
	return types.NewContextAccessor[HelloWorldService](contextKeyNamed("HelloWorldService"))
}
```

Key to type safety is the external invisibility of the context key.  This is guaranteed
by defined a module-local type (`contextKey` or `contextKeyNamed`) and using an instance of it
to index the context inspection/injection.

## Abstract Constructor

To manage the injection of dependencies, go-msx applications use abstract constructors in
the style of go factories.  In particular, they:
- return a reference to an interface type instead of the concrete implementation;
- check the passed-in Context for overrides for the component, and if found, return it;
- fail on error in constructing any subcomponents
- check configuration to select from alternative dependencies

Our service has a single, simple dependency:

```go
func NewHelloWorldService(ctx context.Context) (result HelloWorldService, err error) {
	var ok bool
	if result, ok = ContextHelloWorldService().TryGet(ctx); !ok {
        helloWorldRepository, err := NewHelloWorldRepository(ctx)
		if err != nil { 
			return nil, err
		}
		
		result = &helloWorldService{
            helloWorldRepository: helloWorldRepository,
        }
    }
	
	return
}
```

## Lifecycle Registration

For root components (those not instantiated by other components), you must instantiate them
during application startup.  Components that should be created for all commands should use
the `OnEvent` registration:

```go
func init() {
	var svc helloWorldService
	app.OnEvent(
		app.EventStart,
		app.EventDuring,
		func (ctx context.Context) (err error) {
			svc, err = NewHelloWorldService(ctx)
			return err
        })
}
```
