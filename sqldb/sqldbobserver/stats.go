package sqldbobserver

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"strings"
	"time"
)

const (
	statsSubsystemSql          = "sql"
	statsCounterSqlQueries     = "queries"
	statsGaugeActiveSqlQueries = "active_queries"
	statsCounterConnections    = "connections"
	statsCounterSqlQueryErrors = "query_errors"
	statsTimerSqlQueryTime     = "query_time"
)

var (
	countVecSqlQueries       = stats.NewCounterVec(statsSubsystemSql, statsCounterSqlQueries, "action")
	gaugeVecActiveSqlQueries = stats.NewGaugeVec(statsSubsystemSql, statsGaugeActiveSqlQueries, "action")
	countConnections         = stats.NewGauge(statsSubsystemSql, statsCounterConnections)
	countSqlQueryErrors      = stats.NewCounter(statsSubsystemSql, statsCounterSqlQueryErrors)
	histVecSqlQueryTime      = stats.NewHistogramVec(statsSubsystemSql, statsTimerSqlQueryTime, nil, "action")
)

type errorFunc func() error

func queryLabel(stmt string) string {
	parts := strings.Fields(stmt)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func observeStats(query string, action errorFunc) (err error) {
	start := time.Now()
	label := queryLabel(query)
	countVecSqlQueries.WithLabelValues(label).Inc()
	gaugeVecActiveSqlQueries.WithLabelValues(label).Inc()

	defer func() {
		gaugeVecActiveSqlQueries.WithLabelValues(label).Dec()
		histVecSqlQueryTime.WithLabelValues(label).Observe(float64(time.Since(start)) / float64(time.Millisecond))
		if err != nil {
			countSqlQueryErrors.Inc()
		}
	}()

	err = action()
	return err
}

func ObserveConnection(action errorFunc) (err error) {
	countConnections.Inc()
	return action()
}
