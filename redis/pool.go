// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPool struct {
	conn *Connection
	cfg  *ConnectionConfig
}

func (p *ConnectionPool) WithConnection(action func(*Connection) error) error {
	return action(p.conn)
}

func (p *ConnectionPool) Connection() *Connection {
	return p.conn
}

func (p *ConnectionPool) ConnectionConfig() *ConnectionConfig {
	return p.cfg
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

	if conn, err := NewConnection(ctx); err != nil && err != ErrDisabled {
		return err
	} else if err == ErrDisabled {
		connectionConfig, _ := NewConnectionConfigFromConfig(config.FromContext(ctx))
		pool = &ConnectionPool{
			cfg: connectionConfig,
		}
	} else {
		pool = &ConnectionPool{
			conn: conn,
			cfg:  conn.config,
		}
	}
	return nil
}

type redisContextKey int

const contextKeyRedisPool redisContextKey = iota

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyRedisPool, pool)
}

func PoolFromContext(ctx context.Context) *ConnectionPool {
	connectionPoolInterface := ctx.Value(contextKeyRedisPool)
	if connectionPoolInterface == nil {
		return nil
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		logger.Warn("Context redis connection pool value is the wrong type")
		return nil
	} else {
		return connectionPool
	}
}
