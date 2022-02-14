package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	ErrDisabled           = errors.New("Kafka disabled")
	ErrTopicAlreadyExists = &retry.PermanentError{Cause: errors.New("Topic already exists")}
	logger                = log.NewLogger("msx.kafka")
)

type Connection struct {
	cfg    *ConnectionConfig
	client sarama.Client
}

func (c *Connection) Client() sarama.Client {
	return c.client
}

func (c *Connection) ClusterAdmin() (sarama.ClusterAdmin, error) {
	return sarama.NewClusterAdminFromClient(c.client)
}

func (c *Connection) SyncProducer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(c.client)
}

func (c *Connection) Close() {
	if err := c.client.Close(); err != nil {
		logger.WithError(err).Error("Failed to close sarama client")
	}
}

// NewConnection creates a new Connection using the supplied configuration
func NewConnection(ctx context.Context, cfg *ConnectionConfig) (*Connection, error) {
	saramaConfig, err := cfg.SaramaConfig(ctx)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		cfg: cfg,
	}

	saramaClient, err := sarama.NewClient(cfg.BrokerAddresses(), saramaConfig)
	if err != nil {
		return nil, err
	} else {
		conn.client = saramaClient
	}

	return conn, nil
}
