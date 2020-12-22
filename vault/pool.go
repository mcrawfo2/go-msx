package vault

import (
	"context"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

// Deprecated
type ConnectionPool struct {
	conn *Connection
}

func (p *ConnectionPool) WithConnection(action func(*Connection) error) error {
	return action(p.conn)
}

func (p *ConnectionPool) Connection() *Connection {
	return p.conn
}

// Deprecated.  Use ConnectionFromContext instead.
func Pool() *ConnectionPool {
	return pool
}

// Deprecated
func ConfigurePool(ctx context.Context) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()

	if pool != nil {
		return nil
	}

	if conn, err := NewConnection(ctx); err != nil {
		return err
	} else {
		pool = &ConnectionPool{
			conn: conn,
		}
	}
	return nil
}
