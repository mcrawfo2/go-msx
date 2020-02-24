# MSX Statistics

MSX Statistics allows applications to monitor and record application metrics for display on dashboards and generation of application alarms.  We have chosen to support Prometheus and its OpenMetrics format to expose collected data.

## Statistics Types

MSX Statistics supports several base data collection types:

- **Counter**

    A counter is an ever-increasing number. For example, "Completed Requests" will continuously increase throughout the lifetime of the application.  It is initialized to zero on application startup.

- **Gauge**
    
    A gauge is a metric that represents a single numerical value that can arbitrarily go up and down. For example, "Active Requests" increases when a new request arrives, and decreases when a request is fully serviced.

- **Histogram**
    
    A histogram samples observations (usually things like request durations or response sizes) and counts them in configurable buckets. It also provides a sum of all observed values.  For example, "Query Duration" has a range of time values (from 0 seconds and up).  These can be put into buckets to see what the 99th percentile Query Duration is (using the prometheus `histogram_quantile` function in the dashboard).

If you wish to further group the data, you can use the Vector version of each of the above types.  For example, we can group a "Request Duration Histogram" by API endpoint, in order to see the distribution of request duration distributions for each endpoint separated from other endpoints.

## Usage

### Instantiation

To start collecting a statistic, you must first initialize its collector.  This can be accomplished during module initialization by assigning the collector to a module-global variable:

```go
const (
    statsSubsystemConsul               = "consul"
    statsHistogramConsulCallTime       = "call_time"
    statsGaugeConsulCalls              = "calls"
    statsCounterConsulCallErrors       = "call_errors"
    statsGaugeConsulRegisteredServices = "registrations"
)

var (
    // Collect the number of errors for each api
    countVecConsulCallErrors = stats.NewCounterVec(
        statsSubsystemConsul, 
        statsCounterConsulCallErrors, 
        "api", "param")

    // Collect the number of active requests for each api
    gaugeVecConsulCalls      = stats.NewGaugeVec(
        statsSubsystemConsul, 
        statsGaugeConsulCalls, 
        "api", "param")

    // Collect the distribution of call execution times for each api
    histVecConsulCallTime    = stats.NewHistogramVec(
        statsSubsystemConsul, 
        statsHistogramConsulCallTime, 
        nil, 
        "api", "param")
)
```

As you can see above, each of the collector constructors start with two required arguments:
- **Subsystem**
    
    Identifies the application subsystem being monitored.  In this case, `consul`.

- **Metric Name**
    
    Identifies the individual metric dimension.  By convention, duration histograms end with `_time`, and counters are pluralized.

The histogram (and histogram vector) constructors require an argument specifying the buckets and their upper limits.  To use the default buckets, pass `nil` for this argument.  The current default buckets are calculated by executing `prometheus.ExponentialBuckets(10, 2, 16)`: this evaluates to `[10, 20, 40, ..., 655360]`.  For more information about histograms, you can visit the [Prometheus documentation](https://prometheus.io/docs/practices/histograms/).

Vector constructors, as shown above, accept a final series of dimensions to be applied to each of the measurements.  In the example above, each of our vectors accepts the `api` and `param` groupings.  In the consul stats collector:
- `api` identifies which Consul API endpoint is being called (by path)
- `param` identifies eg. the servicename for discovery

### Collection

After initializing your collectors, you can start to measure your application as the relevant events occur.

A common pattern is define a wrapper function whose only purpose is to collect statistics.  In the Consul package, we can see an example of this:

```go
func observeConsulCall(api, param string, fn func() error) (err error) {
    // Collect the start time of the call
    start := time.Now()
    // Increase the number of active calls
    gaugeVecConsulCalls.WithLabelValues(api, param).Inc()

    // Execute this code before returning, even in case of panic()
    defer func() {
        // Reduce the number of active calls
        gaugeVecConsulCalls.WithLabelValues(api, param).Dec()
        
        // Bucket the call duration in the histogram
        histVecConsulCallTime.WithLabelValues(api, param).Observe(
            float64(time.Since(start)) / float64(time.Millisecond))

        if err != nil {
            // Increase the error count if an error was returned from fn
            countVecConsulCallErrors.WithLabelValues(api, param).Inc()
        }
    }()

    // Call the wrapped function and intercept it's error return value
    err = fn()
    
    // Return the wrapped function's value, after the defer block
    return err
}
```

There are a few things to note here not covered in the inline comments:
1. We directly pass `api` and `param` group values to each of the vectors from the wrapper using `.WithLabelValues()`.  These must be passed in the same order as in the constructor.
2. Time periods should be calculated as `float64` milliseconds.
3. Counters and Gauges can be incremented by `1.0` using the `.Inc()` method.
4. Gauges can be decremented by `1.0` using the `.Dec()` method.
5. Histograms can record an observation using the `.Observe()` method.

## Push Gateway

By default, the MSX Statistics package expects the statistics to be polled by an external application.  If such a poller is not available, MSX Statistics can be configured to push
to an external Prometheus push gateway.

### Configuration

The following configuration settings can be specified to configure the stats pusher:

| Key                   | Description | Default |
|-----------------------|-------------|---------|
| `stats.push.enabled`  | enable the stats pusher | `false` |
| `stats.push.url`      | url to push stats too | |
| `stats.push.job-name` | prometheus job name to send | `go_msx` |
| `stats.push.frequency` | duration between pushes | `15s` |
