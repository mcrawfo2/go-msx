package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/prometheusprovider"
)

func init() {
	OnEvent(EventStart, PhaseBefore, registerPrometheusWebService)
	OnEvent(EventStart, PhaseBefore, startStatsPusher)
	OnEvent(EventStop, PhaseAfter, stopStatsPusher)
}

func registerPrometheusWebService(ctx context.Context) error {
	logger.Info("Registering prometheus metrics endpoint")
	return prometheusprovider.RegisterProvider(ctx)
}

func startStatsPusher(ctx context.Context) error {
	logger.Info("Configuring stats pusher")
	if err := stats.Configure(ctx); err != nil && err != stats.ErrDisabled {
		return err
	} else if err == stats.ErrDisabled {
		logger.WithContext(ctx).Info("Stats pusher disabled.")
		return nil
	}

	logger.Info("Starting stats pusher")
	return stats.Start(ctx)
}

func stopStatsPusher(ctx context.Context) error {
	logger.Info("Stopping stats pusher")
	if err := stats.Stop(ctx); err != nil && err != stats.ErrDisabled {
		return err
	}
	return nil
}
