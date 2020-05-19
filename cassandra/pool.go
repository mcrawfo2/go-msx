package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"github.com/gocql/gocql"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPool struct {
	cfg        *ClusterConfig
	cluster    *Cluster
	session    *gocql.Session
	sessionMtx sync.Mutex
}

func (p *ConnectionPool) withFreshSession(action func(*gocql.Session) error) (err error) {
	session, err := p.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return action(session)
}

func (p *ConnectionPool) withPersistentSession(action func(*gocql.Session) error) (err error) {
	p.sessionMtx.Lock()
	defer p.sessionMtx.Unlock()

	if p.session == nil || p.session.Closed() {
		p.session, err = p.cluster.CreateSession()
		if err != nil {
			return err
		}
	}

	return action(p.session)
}

func (p *ConnectionPool) WithSession(action func(*gocql.Session) error) (err error) {
	if p.cfg.PersistentSessions {
		return p.withPersistentSession(action)
	} else {
		return p.withFreshSession(action)
	}
}

func (p *ConnectionPool) WithSessionRetry(ctx context.Context, action func(*gocql.Session) error) error {
	return p.WithSession(func(session *gocql.Session) error {
		r, err := retry.NewRetryFromContext(ctx)
		if err != nil {
			return err
		}

		return r.Retry(func() error { return action(session) })
	})
}

func (p *ConnectionPool) ClusterConfig() ClusterConfig {
	return *p.cfg
}

func CreateKeyspaceForPool(ctx context.Context) error {
	if pool == nil {
		// Cassandra is disabled
		return nil
	}

	systemClusterConfig := pool.ClusterConfig()
	targetKeyspaceName := systemClusterConfig.KeyspaceName
	systemClusterConfig.KeyspaceName = keyspaceSystem

	logger.WithContext(ctx).Infof("Ensuring keyspace %s exists", targetKeyspaceName)

	cluster, err := NewCluster(&systemClusterConfig)
	if err != nil {
		return err
	}

	return cluster.createKeyspace(ctx, targetKeyspaceName, systemClusterConfig.KeyspaceOptions)
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

	clusterConfig, err := NewClusterConfigFromConfig(cfg)
	if err != nil {
		return err
	}

	cluster, err := NewClusterFromConfig(cfg)
	if err != nil {
		return err
	}

	pool = &ConnectionPool{
		cfg:     clusterConfig,
		cluster: cluster,
	}

	return nil
}
