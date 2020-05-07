package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type Provider interface {
	NewPublisher(cfg *config.Config, name string, configuration *BindingConfiguration) (Publisher, error)
	NewSubscriber(cfg *config.Config, name string, configuration *BindingConfiguration) (Subscriber, error)
}

var (
	providers             = make(map[string]Provider)
	ErrBinderNotEnabled   = errors.New("Binder not enabled")
	ErrConsumerNotEnabled = errors.New("Consumer not enabled")
)

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
		return nil, ErrBinderNotEnabled
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

	if bindingConfig.Consumer.AutoStartup == false {
		return nil, ErrConsumerNotEnabled
	}

	provider, ok := providers[bindingConfig.Binder]
	if !ok {
		return nil, ErrBinderNotEnabled
	}

	subscriber, err := provider.NewSubscriber(cfg, name, bindingConfig)
	if err != nil {
		return nil, err
	}

	return NewStatsSubscriber(subscriber, bindingConfig), nil
}
