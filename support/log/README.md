# MSX Logging Module

MSX logging is an extension of the popular `logrus` logging library, to include:
- Log names
- Level-specific loggers
- Improved context handling

## Usage

After importing the MSX log package, you can use the default named logger `msx` simply:

```
import "cto-github.cisco.com/NFV-BU/go-msx/support/log"

var logger = log.StandardLogger()

func main() {
    var action = "started"
    logger.Infof("Something happened: %s", action) 
}
```

To use a logger with a custom name:

```
var logger = log.NewLogger("alert/api")
```

To create a levelled logger, which outputs Print at the defined log level:

```
debugLogger := logger.Level(log.DebugLevel)
debugLogger.Printf("Some template: %s", "inserted")
```

To add one-time custom diagnostic fields:

```
var logger = log.NewLogger("tenant")

func HandleGetTenantRequest(tenantId string) {
    logger.
        WithExtendedField("tenantId", tenantId).
        Debugf("Tenant retrieval requested")
}
```

To create a sub-logger with custom diagnostic fields:

```
var logger = log.NewLogger("services/tenant")

func HandleGetTenantRequest(tenantId string) {
    requestLogger := logger.WithExtendedField("tenantId", tenantId)
    requestLogger.Debugf("some message")
}
```

To use the log context that was embedded in a Context object:

```
func HandleRequest(ctx context.Context) {
    requestLogger := logger.WithContext(ctx)
    ...
}
```

## Configuration

Output configuration is done directly with logrus:

```
import "github.com/sirupsen/logrus"

// Log as JSON instead of the default ASCII formatter.
logrus.SetFormatter(&logrus.JSONFormatter{})
```

By default, all output is sent to standard output, with high-resolution
timestamps. See [init.go](init.go) for specifics.
