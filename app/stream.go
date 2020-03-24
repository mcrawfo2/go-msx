package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/gochannel"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/kafka"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerKafkaStreamProvider)
	OnEvent(EventConfigure, PhaseAfter, registerGoChannelStreamProvider)
	OnEvent(EventStart, PhaseAfter, stream.StartRouter)
	OnEvent(EventStop, PhaseBefore, stream.StopRouter)
}

func registerKafkaStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if err := kafka.RegisterProvider(cfg); err != nil && err != kafka.ErrDisabled {
		return err
	} else if err == kafka.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("Kafka disabled.  Not registering stream provider.")
	}
	return nil
}

func registerGoChannelStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if err := gochannel.RegisterProvider(cfg); err != nil {
		return err
	}
	return nil
}
