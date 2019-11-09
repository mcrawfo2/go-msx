package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, activateTracing)
	OnEvent(EventFinal, PhaseAfter, deactivateTracing)
}

func activateTracing(ctx context.Context) error {
	logger.Info("Activating tracing")
	return trace.ConfigureTracer(ctx)
}

func deactivateTracing(ctx context.Context) error {
	logger.Info("Deactivating tracing")
	return trace.ShutdownTracer(ctx)
}
