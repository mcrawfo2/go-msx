package app

import "cto-github.cisco.com/NFV-BU/go-msx/support/stats"

func init() {
	OnEvent(EventStart, PhaseBefore, createStatsCollector)
	OnEvent(EventStop, PhaseAfter, closeStatsCollector)
}

func createStatsCollector() {
	stats.Configure(Config())
}

func closeStatsCollector() {
	stats.Close()
}