package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
)

func init() {
	OnEvent(EventStart, PhaseBefore, createStatsCollector)
	OnEvent(EventStop, PhaseAfter, closeStatsCollector)
}

func createStatsCollector(ctx context.Context) error {
	statsContext := config.ContextWithConfig(ctx, applicationConfig)
	return stats.Configure(statsContext)
}

func closeStatsCollector(ctx context.Context) error {
	stats.Close()
	return nil
}