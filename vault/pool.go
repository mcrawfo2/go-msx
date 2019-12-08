package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPool struct {
	conn *Connection
}

func (p *ConnectionPool) WithConnection(action func(*Connection) error) error {
	return action(p.conn)
}

func (p *ConnectionPool) Connection() *Connection {
	return p.conn
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

	if conn, err := NewConnectionFromConfig(cfg); err != nil {
		return err
	} else {
		pool = &ConnectionPool{
			conn: conn,
		}
	}
	return nil
}

type vaultContextKey int

const contextKeyVaultPool vaultContextKey = iota

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyVaultPool, pool)
}

func PoolFromContext(ctx context.Context) *ConnectionPool {
	connectionPoolInterface := ctx.Value(contextKeyVaultPool)
	if connectionPoolInterface == nil {
		return nil
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		logger.Warn("Context vault connection pool value is the wrong type")
		return nil
	} else {
		return connectionPool
	}
}
