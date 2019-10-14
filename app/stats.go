package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
)

func init() {
	OnEvent(EventStart, PhaseBefore, createStatsCollector)
	OnEvent(EventStop, PhaseAfter, closeStatsCollector)
}

func createStatsCollector(ctx context.Context) error {
	return stats.Configure(ctx)
}

func closeStatsCollector(ctx context.Context) error {
	stats.Close()
	return nil
}