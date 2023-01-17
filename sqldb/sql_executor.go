// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type SqlExecutor interface {
	DriverName() string
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func ContextSqlExecutor() types.ContextKeyAccessor[SqlExecutor] {
	return types.NewContextKeyAccessor[SqlExecutor](contextKeySqlExecutor)
}

type SqlExecutorAction func(ctx context.Context, sqlExecutor SqlExecutor) error

func WithSqlExecutor(ctx context.Context, action SqlExecutorAction) error {
	sqlExecutor := ContextSqlExecutor().Get(ctx)

	if sqlExecutor == nil {
		pool, err := PoolFromContext(ctx)
		if err != nil {
			return err
		}

		return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
			return action(ctx, conn)
		})

	} else {
		// connection should have been established
		return action(ctx, sqlExecutor)
	}
}

func SqlDriverName(ctx context.Context) (string, error) {
	driverName, err := config.FromContext(ctx).String(
		config.PrefixWithName(configRootSpringDatasourceConfig, "driver"))
	if err != nil {
		return "", err
	}
	return driverName, nil
}
