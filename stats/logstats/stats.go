package logstats

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/sirupsen/logrus"
)

const (
	statsSubsystemLogging = "log"
)

var (
	loggingStatsCounterOutputVec = stats.NewCounterVec(statsSubsystemLogging, "count", "logger", "level")
)

type LoggingStatsCollector struct {
}

func (l LoggingStatsCollector) Inc(loggerName string, level logrus.Level) {
	loggingStatsCounterOutputVec.WithLabelValues(loggerName, log.LoggerLevel(level).Name()).Inc()
}

func init() {
	log.RegisterStatsCollector(new(LoggingStatsCollector))
}
