package health

import "cto-github.cisco.com/NFV-BU/go-msx/stats"

const (
	statsSubsystemHealth           = "health"
	statsCounterHealthCheckReports = "check_reports"
	statsGaugeHealthChecks         = "checks"

	statsHistogramHealthCheckStatus = "check_status"
	statsGaugeHealthChecksUp        = "checks_up"
	statsGaugeHealthChecksDown      = "checks_down"
	statsGaugeHealthChecksUnknown   = "checks_unknown"
)

var (
	counterHealthCheckReports = stats.NewCounter(statsSubsystemHealth, statsCounterHealthCheckReports)
	gaugeHealthChecks         = stats.NewGauge(statsSubsystemHealth, statsGaugeHealthChecks)
	gaugeHealthChecksUp       = stats.NewGauge(statsSubsystemHealth, statsGaugeHealthChecksUp)
	gaugeHealthChecksDown     = stats.NewGauge(statsSubsystemHealth, statsGaugeHealthChecksDown)
	gaugeHealthChecksUnknown  = stats.NewGauge(statsSubsystemHealth, statsGaugeHealthChecksUnknown)
	histHealthCheckStatus     = stats.NewHistogram(statsSubsystemHealth, statsHistogramHealthCheckStatus, nil)
)

func sendReportStats(report *Report) {
	counterHealthCheckReports.Inc()
	gaugeHealthChecks.Set(float64(len(report.Details)))
	gaugeHealthChecksUp.Set(float64(report.Up()))
	gaugeHealthChecksDown.Set(float64(report.Down()))
	gaugeHealthChecksUnknown.Set(float64(report.Unknown()))
	if len(report.Details) > 0 {
		histHealthCheckStatus.Observe(float64(report.Up()) / float64(len(report.Details)))
	}
}
