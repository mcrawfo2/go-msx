package health

import "cto-github.cisco.com/NFV-BU/go-msx/stats"

const (
	statsGaugeHealthCheckReports  = "health.checkReports"
	statsGaugeHealthChecks        = "health.checks"
	statsGaugeHealthChecksUp      = "health.checksUp"
	statsGaugeHealthChecksDown    = "health.checksDown"
	statsGaugeHealthChecksUnknown = "health.checksUnknown"
)

func sendReportStats(report *Report) {
	stats.Incr(statsGaugeHealthCheckReports, 1)
	stats.Gauge(statsGaugeHealthChecks, int64(len(report.Details)))
	stats.Gauge(statsGaugeHealthChecksUp, int64(report.Up()))
	stats.Gauge(statsGaugeHealthChecksDown, int64(report.Down()))
	stats.Gauge(statsGaugeHealthChecksUnknown, int64(report.Unknown()))
}
