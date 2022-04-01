// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cassandra

import (
	"context"
	"github.com/pkg/errors"
)

type contextKey int

const (
	contextKeyCassandraPool contextKey = iota
	contextKeyCrudRepositoryFactory
)

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyCassandraPool, pool)
}

func PoolFromContext(ctx context.Context) (*ConnectionPool, error) {
	connectionPoolInterface := ctx.Value(contextKeyCassandraPool)
	if connectionPoolInterface == nil {
		return nil, ErrDisabled
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		return nil, errors.New("Context cassandra connection pool value is the wrong type")
	} else {
		return connectionPool, nil
	}
}

func ContextWithCrudRepositoryFactory(ctx context.Context, factory CrudRepositoryFactoryApi) context.Context {
	return context.WithValue(ctx, contextKeyCrudRepositoryFactory, factory)
}

func CrudRepositoryFactoryFromContext(ctx context.Context) CrudRepositoryFactoryApi {
	api, _ := ctx.Value(contextKeyCrudRepositoryFactory).(CrudRepositoryFactoryApi)
	return api
}
