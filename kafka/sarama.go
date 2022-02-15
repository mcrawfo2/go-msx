package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
)

var (
	saramaLogger = log.NewLogger("sarama")
)

func init() {
	metrics.UseNilMetrics = true
	sarama.Logger = saramaLogger
	saramaLogger.SetLevel(log.WarnLevel)
}

func NewSaramaClient(config *ConnectionConfig) (sarama.Client, error) {
	saramaConfig, err := config.SaramaConfig(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to configure kafka client")
	}

	brokerAddresses := config.BrokerAddresses()
	if len(brokerAddresses) == 0 {
		return nil, errors.New("No brokers defined")
	}

	return sarama.NewClient(brokerAddresses, saramaConfig)
}

func NewSyncProducer(config *ConnectionConfig) (sarama.SyncProducer, error) {
	saramaClient, err := NewSaramaClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create kafka client")
	}

	return sarama.NewSyncProducerFromClient(saramaClient)
}

func NewClusterAdmin(config *ConnectionConfig) (sarama.ClusterAdmin, error) {
	saramaConfig, err := config.SaramaConfig(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to configure kafka client")
	}

	brokerAddresses := config.BrokerAddresses()
	if len(brokerAddresses) == 0 {
		return nil, errors.New("No brokers defined")
	}

	return sarama.NewClusterAdmin(brokerAddresses, saramaConfig)
}
