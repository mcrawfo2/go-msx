# MSX Logging Module

MSX logging is an extension of the popular `logrus` logging library, to include:
- Log names
- Level-specific loggers
- Improved context handling

## Usage

After importing the MSX log package, you can use the default named logger `msx` simply:

```go
import "cto-github.cisco.com/NFV-BU/go-msx/log"

var logger = log.StandardLogger()

func main() {
    var action = "started"
    logger.Infof("Something happened: %s", action) 
}
```

To use a logger with a custom name:

```go
var logger = log.NewLogger("alert.api")
```

To create a levelled logger, which outputs Print at the defined log level:

```go
debugLogger := logger.Level(log.DebugLevel)
debugLogger.Printf("Some template: %s", "inserted")
```

To record a golang `error` object:

```go
func DeployResource(data []byte) {
    var body ResourceDeployment
    if err := json.Unmarshal(data, &body); err != nil {
        logger.
            WithError(err).
            Error("Failed to parse Resource Deployment request")
    }
}
```

To use the log context that was embedded in a Context object:

```go
func HandleRequest(ctx context.Context) {
    requestLogger := logger.WithContext(ctx)
    ...
}
```

To add one-time custom diagnostic fields:

```go
var logger = log.NewLogger("tenant")

func HandleGetTenantRequest(tenantId string) {
    logger.
        WithExtendedField("tenantId", tenantId).
        Debug("Tenant retrieval requested")
}
```

To create a sub-logger with custom diagnostic fields:

```go
var logger = log.NewLogger("services.tenant")

func HandleGetTenantRequest(tenantId string) {
    requestLogger := logger.WithExtendedField("tenantId", tenantId)
    requestLogger.Debugf("some message")
}
```

## Configuration

Output configuration is done directly with logrus:

```go
import "github.com/sirupsen/logrus"

// Log as JSON instead of the default ASCII formatter.
logrus.SetFormatter(&logrus.JSONFormatter{})
```

By default, all output is sent to standard output, with high-resolution
timestamps. See [init.go](init.go) for specifics.
