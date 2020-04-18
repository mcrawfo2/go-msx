package kafka

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	ErrDisabled = errors.New("Kafka disabled")
	logger      = log.NewLogger("msx.kafka")
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

func NewConnection(cfg *ConnectionConfig) (*Connection, error) {
	saramaConfig, err := cfg.SaramaConfig()
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		cfg: cfg,
	}

	if saramaClient, err := sarama.NewClient(cfg.BrokerAddresses(), saramaConfig); err != nil {
		return nil, err
	} else {
		conn.client = saramaClient
	}

	return conn, nil
}

func NewConnectionFromConfig(cfg *config.Config) (*Connection, error) {
	connectionConfig, err := NewConnectionConfig(cfg)
	if err != nil {
		return nil, err
	}

	if !connectionConfig.Enabled {
		return nil, ErrDisabled
	}

	return NewConnection(connectionConfig)
}
