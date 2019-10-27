package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"sync"
	"time"
	"vitess.io/vitess/go/pools"
)

const (
	configRootKafkaPool = "spring.cloud.stream.kafka.binder.pool"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPoolConfig struct {
	Enabled            bool          `config:"default=true"`
	Capacity           int           `config:"default=1"`
	MaxCapacity        int           `config:"default=24"`
	IdleTimeout        time.Duration `config:"default=60s"`
	PrefillParallelism int           `config:"default=0"`
}

func NewConnectionPoolConfig(cfg *config.Config) (*ConnectionPoolConfig, error) {
	var poolConfig ConnectionPoolConfig
	if err := cfg.Populate(&poolConfig, configRootKafkaPool); err != nil {
		return nil, errors.Wrap(err, "Failed to populate kafka pool config")
	}

	return &poolConfig, nil
}

type ConnectionPool struct {
	pool *pools.ResourcePool
	cfg  *ConnectionConfig
}

func (p *ConnectionPool) WithConnection(ctx context.Context, action func(*Connection) error) error {
	connResource, err := p.pool.Get(ctx)
	if err != nil {
		return err
	}
	defer p.pool.Put(connResource)

	conn, ok := connResource.(*Connection)
	if !ok {
		return errors.New("Failed to retrieve connection")
	}

	return action(conn)
}

func (p *ConnectionPool) WithAsyncProducer(ctx context.Context, action func(sarama.AsyncProducer) error) error {
	return p.WithConnection(ctx, func(conn *Connection) error {
		asyncProducer, err := sarama.NewAsyncProducerFromClient(conn.Client())
		if err != nil {
			return errors.Wrap(err, "Failed to create async producer")
		}

		return action(asyncProducer)
	})
}

func (p *ConnectionPool) WithSyncProducer(ctx context.Context, action func(sarama.SyncProducer) error) error {
	return p.WithConnection(ctx, func(conn *Connection) error {
		syncProducer, err := sarama.NewSyncProducerFromClient(conn.Client())
		if err != nil {
			return errors.Wrap(err, "Failed to create sync producer")
		}

		return action(syncProducer)
	})
}

func (p *ConnectionPool) WithConsumer(ctx context.Context, action func(sarama.Consumer) error) error {
	return p.WithConnection(ctx, func(conn *Connection) error {
		consumer, err := sarama.NewConsumerFromClient(conn.Client())
		if err != nil {
			return errors.Wrap(err, "Failed to create consumer")
		}

		return action(consumer)
	})
}

func (p *ConnectionPool) WithClusterAdmin(ctx context.Context, action func(sarama.ClusterAdmin) error) error {
	return p.WithConnection(ctx, func(conn *Connection) error {
		clusterAdmin, err := sarama.NewClusterAdminFromClient(conn.Client())
		if err != nil {
			return errors.Wrap(err, "Failed to create consumer")
		}

		return action(clusterAdmin)
	})
}

func (p *ConnectionPool) connection() (pools.Resource, error) {
	return NewConnection(p.cfg)
}

func (p *ConnectionPool) ConnectionConfig() *ConnectionConfig {
	return p.cfg
}

func Pool() *ConnectionPool {
	return pool
}

func ConfigurePool(cfg *config.Config) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()

	if pool != nil {
		return nil
	}

	connConfig, err := NewConnectionConfig(cfg)
	if err != nil {
		return err
	}

	if !connConfig.Enabled {
		return ErrDisabled
	}

	poolConfig, err := NewConnectionPoolConfig(cfg)
	if err != nil {
		return err
	}

	if !poolConfig.Enabled {
		return ErrDisabled
	}

	p := ConnectionPool{
		cfg: connConfig,
	}

	p.pool = pools.NewResourcePool(
		p.connection,
		poolConfig.Capacity,
		poolConfig.MaxCapacity,
		poolConfig.IdleTimeout,
		poolConfig.PrefillParallelism)

	pool = &p

	return nil
}

type kafkaContextKey int

const contextKeyKafkaPool kafkaContextKey = iota

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyKafkaPool, pool)
}

func PoolFromContext(ctx context.Context) *ConnectionPool {
	connectionPoolInterface := ctx.Value(contextKeyKafkaPool)
	if connectionPoolInterface == nil {
		return nil
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		logger.Warn("Context kafka connection pool value is the wrong type")
		return nil
	} else {
		return connectionPool
	}
}
