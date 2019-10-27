package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/kafka"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(registerKafkaStreamProvider))
	OnEvent(EventStart, PhaseAfter, stream.StartRouter)
	OnEvent(EventStop, PhaseBefore, stream.StopRouter)
}

func registerKafkaStreamProvider(cfg *config.Config) error {
	if err := kafka.RegisterProvider(cfg); err != nil && err != kafka.ErrDisabled {
		return err
	}
	return nil
}