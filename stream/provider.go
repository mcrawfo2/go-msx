package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type Provider interface {
	NewPublisher(cfg *config.Config, name string, configuration *BindingConfiguration) (Publisher, error)
	NewSubscriber(cfg *config.Config, name string, configuration *BindingConfiguration) (Subscriber, error)
}

var providers = make(map[string]Provider)

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

func NewPublisher(cfg *config.Config, name string) (Publisher, error) {
	bindingConfig, err := NewBindingConfigurationFromConfig(cfg, name)
	if err != nil {
		return nil, err
	}

	provider, ok := providers[bindingConfig.Binder]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Binder not found: %s", bindingConfig.Binder))
	}

	publisher, err := provider.NewPublisher(cfg, name, bindingConfig)
	if err != nil {
		return nil, err
	}

	return NewStatsPublisher(publisher, bindingConfig), nil
}

func NewSubscriber(cfg *config.Config, name string) (message.Subscriber, error) {
	bindingConfig, err := NewBindingConfigurationFromConfig(cfg, name)
	if err != nil {
		return nil, err
	}

	provider, ok := providers[bindingConfig.Binder]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Binder not found: %s", bindingConfig.Binder))
	}

	subscriber, err := provider.NewSubscriber(cfg, name, bindingConfig)
	if err != nil {
		return nil, err
	}

	return NewStatsSubscriber(subscriber, bindingConfig), nil
}
