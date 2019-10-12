package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/support/stats"
)

func init() {
	OnEvent(EventStart, PhaseBefore, createStatsCollector)
	OnEvent(EventStop, PhaseAfter, closeStatsCollector)
}

func createStatsCollector(ctx context.Context) error {
	stats.Configure(Config())
	return nil
}

func closeStatsCollector(ctx context.Context) error {
	stats.Close()
	return nil
}