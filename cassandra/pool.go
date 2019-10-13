package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/gocql/gocql"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPool struct {
	cluster *Cluster
}

func (p *ConnectionPool) WithSession(action func(*gocql.Session) error) error {
	session, err := p.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return action(session)
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

	if conn, err := NewClusterFromConfig(cfg); err != nil {
		return err
	} else {
		pool = &ConnectionPool{
			cluster: conn,
		}
	}
	return nil
}

type cassandraContextKey int
const contextKeyCassandraPool cassandraContextKey = iota

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyCassandraPool, pool)
}

func PoolFromContext(ctx context.Context) *ConnectionPool {
	connectionPoolInterface := ctx.Value(contextKeyCassandraPool)
	if connectionPoolInterface == nil {
		return nil
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		logger.Warn("Context cassandra connection pool value is the wrong type")
		return nil
	} else {
		return connectionPool
	}
}
