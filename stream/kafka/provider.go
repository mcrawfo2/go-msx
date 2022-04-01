// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	msxKafka "cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/pkg/errors"
)

const (
	providerNameKafka = "kafka"
)

var ErrDisabled = msxKafka.ErrDisabled
var ErrTopicAlreadyExists = msxKafka.ErrTopicAlreadyExists
var loggerWatermillKafka = log.NewLogger("watermill.kafka")
var loggerAdapter = stream.NewWatermillLoggerAdapter(loggerWatermillKafka)

type Provider struct{}

func (p *Provider) NewPublisher(ctx context.Context, name string, streamBinding *stream.BindingConfiguration) (stream.Publisher, error) {
	connectionConfig := msxKafka.Pool().ConnectionConfig()

	saramaConfig, err := connectionConfig.SaramaConfig(ctx)
	if err != nil {
		return nil, err
	}

	bindingConfig, err := NewBindingConfiguration(ctx, name, streamBinding)
	if err != nil {
		return nil, err
	}

	if bindingConfig.Producer.Sync == false {
		saramaConfig.Producer.Return.Successes = false
	}

	if connectionConfig.AutoCreateTopics {
		err = retry.NewRetry(ctx, retry.RetryConfig{
			Attempts: 3,
			Delay:    500,
			BackOff:  0.0,
			Linear:   true,
		}).Retry(func() error {
			return msxKafka.Pool().WithConnection(ctx, func(connection *msxKafka.Connection) error {
				return msxKafka.CreateTopics(ctx, connection, streamBinding.Destination)
			})
		})
		if errors.Is(err, ErrTopicAlreadyExists) {
			err = nil
		} else if err != nil {
			return nil, errors.Wrap(err, "Failed to create topic")
		}
	}

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               connectionConfig.BrokerAddresses(),
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaConfig,
		},
		loggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	return stream.NewTopicPublisher(publisher, bindingConfig.StreamBindingConfig), nil
}

func (p *Provider) NewSubscriber(ctx context.Context, name string, streamBinding *stream.BindingConfiguration) (stream.Subscriber, error) {
	connectionConfig := msxKafka.Pool().ConnectionConfig()
	saramaConfig, err := connectionConfig.SaramaConfig(ctx)
	if err != nil {
		return nil, err
	}
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               connectionConfig.BrokerAddresses(),
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaConfig,
			ConsumerGroup:         streamBinding.Group,
			InitializeTopicDetails: &sarama.TopicDetail{
				NumPartitions:     int32(connectionConfig.DefaultPartitions),
				ReplicationFactor: int16(connectionConfig.ReplicationFactor),
			},
		},
		loggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func RegisterProvider(cfg *config.Config) error {
	kafkaConfig, err := msxKafka.NewConnectionConfig(cfg)
	if err != nil {
		return err
	}

	if !kafkaConfig.Enabled {
		return ErrDisabled
	}

	stream.RegisterProvider(providerNameKafka, &Provider{})
	return nil
}
