package sqldb

import (
	"context"
	"github.com/pkg/errors"
)

type contextKey int

const (
	contextKeySqlPool contextKey = iota
)

var ErrDisabled = errors.New("Sql connection disabled")

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeySqlPool, pool)
}

func PoolFromContext(ctx context.Context) (*ConnectionPool, error) {
	connectionPoolInterface := ctx.Value(contextKeySqlPool)
	if connectionPoolInterface == nil {
		return nil, ErrDisabled
	}
	return connectionPoolInterface.(*ConnectionPool), nil
}
