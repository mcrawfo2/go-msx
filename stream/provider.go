// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type Provider interface {
	NewPublisher(ctx context.Context, name string, configuration *BindingConfiguration) (Publisher, error)
	NewSubscriber(ctx context.Context, name string, configuration *BindingConfiguration) (Subscriber, error)
}

var (
	providers             = make(map[string]Provider)
	ErrBinderNotEnabled   = errors.New("Binder not enabled")
	ErrConsumerNotEnabled = errors.New("Consumer not enabled")
	ErrDisconnected       = errors.New("Disconnected mode")
)

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

// NewPublisher creates a Publisher instance based on the specified binding name
func NewPublisher(ctx context.Context, name string) (Publisher, error) {
	bindingConfig, err := NewBindingConfiguration(ctx, name)
	if err != nil {
		return nil, err
	}

	provider, ok := providers[bindingConfig.Binder]
	if !ok {
		return nil, ErrBinderNotEnabled
	}

	publisher, err := provider.NewPublisher(ctx, name, bindingConfig)
	if err != nil {
		return nil, err
	}

	return NewStatsPublisher(publisher, bindingConfig), nil
}

// NewSubscriber creates a Subscriber instance based on the specified binding name
func NewSubscriber(ctx context.Context, name string) (message.Subscriber, error) {
	bindingConfig, err := NewBindingConfiguration(ctx, name)
	if err != nil {
		return nil, err
	}
	if bindingConfig.Disconnected {
		return nil, ErrDisconnected
	}

	if bindingConfig.Consumer.AutoStartup == false {
		return nil, ErrConsumerNotEnabled
	}

	provider, ok := providers[bindingConfig.Binder]
	if !ok {
		return nil, ErrBinderNotEnabled
	}

	subscriber, err := provider.NewSubscriber(ctx, name, bindingConfig)
	if err != nil {
		return nil, err
	}

	return NewStatsSubscriber(subscriber, bindingConfig), nil
}
