# MSX Scheduled Module

MSX Scheduled manages the periodic execution of tasks within microservices.

## Tasks

The work to be performed on a periodic basis must be surrounded in an Action (a function signature matching`types.ActionFunc`):

```go
func doWork(ctx context.Context) error {
  // TODO: perform the desired work.
} 
```

Actions can be anonymous functions, struct methods (as above), or static methods, and can also be derived from Operations (`types.Operation`).

## Scheduling

Scheduling a task requires two steps: Configuration and Registration.

### Configuration

To configure the periodic execution, your task will need a simple name to identify its configuration.  For example, the `do-work` task can be configured as:

```yaml
scheduled.tasks:
    do-work:
        fixed-interval: 10m
        # fixed-delay: 5m
        # initial-delay: 15m
        # cron-expression: "0 0 0 * *"
```

This example configuration will execute the `do-work` task (once registered) every 10 minutes.

To ensure a fixed period _between_ executions, use the `fixed-delay` configuration instead.

To specify an initial delay before first execution that is different from `fixed-delay` or `fixed-interval`, specify the `initial-delay`.

To use a CRON expression to specify the execution schedule, use the `cron-expression` configuration.  For an overview of CRON expressions, see [here](https://en.wikipedia.org/wiki/Cron).

### Registration

To register your task at runtime, call the `scheduled.ScheduleTask` function during the application Start:

```go
const taskNameDoWork = "do-work"

func init() {
  app.OnRootEvent(app.EventStart, app.PhaseAfter, func(ctx context.Context) error {
        return scheduled.ScheduleTask(ctx, taskNameDoWork, doWork)
  })
}
```

This will load the configuration using the supplied task name, and schedule the task according to the configuration.

