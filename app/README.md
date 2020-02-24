# MSX Application Module

MSX Application is a simple state machine for managing application lifecycle.  It installs observers to configure and instantiate the standard components for use with MSX applications.  This includes listeners for external events to advance the state machine (e.g. POSIX signals, configuration changes).

## Lifecycle Events

MSX Application defines various lifecycle events:
- `app.EventCommand` - mode selection based on CLI sub-commands
- `app.EventInit` - pre-configure application 
- `app.EventConfigure` - application and component configuration 
- `app.EventStart` - start services for consumers
- `app.EventReady` - application fully initialized and ready to service requests
- `app.EventRefresh` - update configuration after change
- `app.EventStop` - stop services for consumers
- `app.EventFinalize` - pre-termination cleanup

Each event (except `app.EventCommand`) proceeds with three phases:
- `app.PhaseBefore` - early
- `app.PhaseDuring` - normal
- `app.PhaseAfter` - late

The `app.EventCommand` event will be execute with the `phase` containing the command being executed.  The following commands are pre-defined:
- `app.CommandRoot` - Root (default)
- `app.CommandMigrate` - Migrate
- `app.CommandPopulate` - Populate

### Event Observers

When a lifecycle event phase is occurring, the MSX Application will call each of the `Observer`s registered for the event phase.  These callbacks can be registered during any previous lifecycle event callback or during your module `init()`.  

For example, to call the `addWebService` observer during `start.before` for all commands:

```go
func init() {
    app.OnEvent(app.EventStart, app.PhaseBefore, addWebService)
}
```

To see an example showing command-specific event observers, see [Commanding](#Commanding), below.

### Short-Circuiting

Sometimes an application is not able to correctly execute a lifecycle phase, or receives an external interruption.  This will result in a short-circuit of the lifecycle.  If an error is returned from one of the observers in the following phases, the lifecycle will move to the specified phase:

- `app.EventInit` => `app.EventFinalize`
- `app.EventConfigure` => `app.EventFinalize`
- `app.EventStart` => `app.EventStop`
- `app.EventReady` => `app.EventStop`

## Application Observers

### Command

The `app.EventCommand` events are the first events fired during startup.  They provide the opportunity to execute custom logic and register event observers specific to the command.

As above, the `app.EventCommand` event will be executed with the phase containing the command being executed.  For example the phase could be one of the default commands:
- Root (`app.CommandRoot`)
- Migrate (`app.CommandMigrate`)
- Populate (`app.CommandPopulate`)

To add a new command:

```go
func main() {
    if _, err := app.AddCommand("token", "Create OAuth2 token", renew, app.Noop); err != nil {
        cli.Fatal(err)
    }
}
```

To configure event observers in response to a specific command being executed:

```go
func init() {
    app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
        app.OnEvent(app.EventStart, app.PhaseBefore, addWebService)
        return nil
    })
}
```

### Init

The `app.EventInit` events are fired second, after the `app.EventCommand` events.

Observers attached to the `app.EventInit` events should be restricted to modifying the application environment.  This includes registering custom config providers or custom context injectors.

### Configure

The `app.EventConfigure` events are fired third during startup, after the `app.EventInit` events.

By default, the application is configured:
- `app.PhaseBefore`
    - Register remote config providers
- `app.PhaseDuring`
    - [Load configuration](#configuration-loading)
- `app.PhaseAfter`
    - HTTP Client
    - Consul connection pool
    - Vault connection pool
    - Cassandra connection pool
    - Redis connection pool
    - Kafka connection pool
    - Web server
    - Create Cassandra Keyspace

Typically, user applications will not register new event handlers for the `app.EventConfigure` events.

### Start

The `app.EventStart` events are fired fourth during startup, after the `app.EventConfigure` events.

By default, application infrastructure is connected:
- `app.PhaseBefore`:
  - Authentication Providers
  - Spring Actuators
  - Swagger
  - Prometheus Actuator
  - Stats Pusher
- `app.PhaseAfter`:
  - Health logging
  - Stream Router
  - Web Server
  - Config Watcher

Custom application startup code is expected to run inside the `app.PhaseDuring` phase.  This would include starting any long-running services or scheduling background tasks.

### Ready

The `app.EventReady` events are fired fifth during startup, after the `app.EventStart` events.

By default, application ready observers are executed:
- `app.PhaseBefore`:
  - Service Registration (consul)
- `app.PhaseAfter`:
  - Command Execution (sub-commands)

### Refresh

*TBD*

### Stop

The `app.EventStop` events are fired first during shutdown.

By default, application services are stopped and infrastructure and disconnected:

- `app.PhaseBefore`:
  - Service De-Registration (consul)
  - Health logging
  - Stream router
  - Web Server
  - Stats Pusher

Any custom application code running in the background should be shutdown during `app.PhaseDuring`.

### Finalize

The `app.EventFinal` events are fired last during shutdown.

By default, tracing is stopped during `app.PhaseAfter` to allow trace collection to include `app.EventStop`.

## Configuration Loading

In response to the `app.EventConfigure` event, MSX Application combines all registered sources of configuration.  This occurs in three phases:
- **Phase 1** - In-Memory
  - Application Static Defaults
  - Environment Variables
  - Application Runtime Overrides
  - Command Line
- **Phase 2** - Filesystem
  - Defaults Files
  - Bootstrap Files
  - Application Files
  - Profile Files 
  - Build Files
- **Phase 3** - Remote
  - Consul
  - Vault

Note that this loading order is not the same as the order of precendence for calculating values:
  - Application Static Defaults
  - Defaults Files
  - Bootstrap Files
  - Application Files
  - Build Files
  - Consul
  - Vault
  - Profile Files
  - Environment Variables
  - Command Line
  - Application Runtime Overrides
