// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/Shopify/sarama"
	"github.com/jackc/puddle"
	"github.com/pkg/errors"
	"sync"
	"time"
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
	IdleCleanup        time.Duration `config:"default=15s"`
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
	pool       *puddle.Pool
	poolConfig *ConnectionPoolConfig
	connConfig *ConnectionConfig
}

func (p *ConnectionPool) IdleCleanup(ctx context.Context) {
	logger.WithContext(ctx).Info("Starting kafka connection pool idle cleanup.")
	timer := time.NewTimer(p.poolConfig.IdleCleanup)

	for {
		select {
		case <-ctx.Done():
			logger.WithContext(ctx).WithError(ctx.Err()).Info("Stopping kafka connection pool idle cleanup.")
			return

		case <-timer.C:
			logger.WithContext(ctx).Debug("Cleaning kafka idle connections")
			for _, resource := range p.pool.AcquireAllIdle() {
				if resource.IdleDuration() > p.poolConfig.IdleTimeout {
					resource.Destroy()
				} else {
					resource.ReleaseUnused()
				}
			}
			timer.Reset(p.poolConfig.IdleCleanup)
		}
	}
}

func (p *ConnectionPool) WithConnection(ctx context.Context, action func(*Connection) error) error {
	connResource, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	conn, ok := connResource.Value().(*Connection)
	if !ok {
		return errors.New("Failed to retrieve connection")
	}

	defer func() {
		if !conn.Client().Closed() {
			connResource.Release()
		} else {
			connResource.Destroy()
		}
	}()

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

func (p *ConnectionPool) construct(ctx context.Context) (interface{}, error) {
	return NewConnection(ctx, p.connConfig)
}

func (p *ConnectionPool) destruct(res interface{}) {
	conn := res.(*Connection)
	conn.Close()
}

func (p *ConnectionPool) ConnectionConfig() *ConnectionConfig {
	return p.connConfig
}

func Pool() *ConnectionPool {
	return pool
}

func ConfigurePool(ctx context.Context) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()

	if pool != nil {
		return nil
	}

	cfg := config.MustFromContext(ctx)
	connConfig, err := NewConnectionConfig(cfg)
	if err != nil {
		return err
	}

	if !connConfig.Enabled {
		pool = &ConnectionPool{
			connConfig: connConfig,
		}
		return ErrDisabled
	}

	poolConfig, err := NewConnectionPoolConfig(cfg)
	if err != nil {
		return err
	}

	if !poolConfig.Enabled {
		pool = &ConnectionPool{
			connConfig: connConfig,
		}
		return ErrDisabled
	}

	p := ConnectionPool{
		connConfig: connConfig,
		poolConfig: poolConfig,
	}

	p.pool = puddle.NewPool(
		p.construct,
		p.destruct,
		int32(poolConfig.MaxCapacity))

	pool = &p

	go pool.IdleCleanup(ctx)

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
