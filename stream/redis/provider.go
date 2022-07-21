// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"github.com/ThreeDotsLabs/watermill/message"
	redisstream "github.com/minghsu0107/watermill-redistream/pkg/redis"
	"github.com/pkg/errors"
)

const (
	providerNameRedis = "redis"
)

var ErrDisabled = redis.ErrDisabled
var loggerWatermillRedis = log.NewLogger("watermill.redis")
var loggerAdapter = stream.NewWatermillLoggerAdapter(loggerWatermillRedis)

type Provider struct{}

func (p *Provider) newPublisher(ctx context.Context) (message.Publisher, error) {
	client := redis.Pool().Connection().Client(ctx)

	publisherConfig := redisstream.PublisherConfig{}

	publisher, err := redisstream.NewPublisher(ctx,
		publisherConfig,
		client,
		new(redisstream.DefaultMarshaller),
		loggerAdapter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create SQL publisher")
	}

	return publisher, nil
}

func (p *Provider) newSubscriber(ctx context.Context, name string, streamBinding *stream.BindingConfiguration) (message.Subscriber, error) {
	client := redis.Pool().Connection().Client(ctx)

	bindingConfiguration, err := NewBindingConfiguration(ctx, name, streamBinding)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create binding configuration")
	}

	subscriberConfig := redisstream.SubscriberConfig{
		Consumer:        bindingConfiguration.Consumer.ClientId + "-" + bindingConfiguration.Consumer.ClientIdSuffix,
		ConsumerGroup:   bindingConfiguration.StreamBindingConfig.Group,
		NackResendSleep: bindingConfiguration.Consumer.NackResendSleep,
		MaxIdleTime:     bindingConfiguration.Consumer.MaxIdleTime,
	}

	subscriber, err := redisstream.NewSubscriber(ctx, subscriberConfig, client, new(redisstream.DefaultMarshaller), loggerAdapter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create redis subscriber")
	}

	return subscriber, nil
}

func (p *Provider) NewPublisher(ctx context.Context, _ string, streamBinding *stream.BindingConfiguration) (stream.Publisher, error) {
	sqlPublisher, err := p.newPublisher(ctx)
	if err != nil {
		return nil, err
	}

	publisher := stream.NewTopicPublisher(sqlPublisher, streamBinding)
	return publisher, nil
}

func (p *Provider) NewSubscriber(ctx context.Context, name string, streamBinding *stream.BindingConfiguration) (stream.Subscriber, error) {
	sqlSubscriber, err := p.newSubscriber(ctx, name, streamBinding)
	if err != nil {
		return nil, err
	}

	return sqlSubscriber, nil
}

func RegisterProvider(cfg *config.Config) error {
	sqlConfig, err := redis.NewConnectionConfigFromConfig(cfg)
	if err != nil {
		return err
	}

	if !sqlConfig.Enable {
		return ErrDisabled
	}

	stream.RegisterProvider(providerNameRedis, &Provider{})
	return nil
}
