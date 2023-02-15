// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"errors"
	"github.com/jmoiron/sqlx"
)

const (
	ConfigRootTransactionManager        = "sqldb.transaction-manager"
	ConfigTransactionManagerMaxRetries  = ConfigRootTransactionManager + ".max-retries"
	TransactionManagerMaxRetriesDefault = 5
)

func TransactionDecorator(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		pool, err := PoolFromContext(ctx)
		if err != nil {
			return err
		}

		err = pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
			tx, err := conn.Beginx()
			if err != nil {
				return err
			}

			ctx = ContextSqlExecutor().Set(ctx, tx)

			cfg := config.MustFromContext(ctx)
			driver, err := cfg.String("spring.datasource.driver")
			if err != nil {
				return err
			}
			maxRetries, err := cfg.IntOr(ConfigTransactionManagerMaxRetries, TransactionManagerMaxRetriesDefault)
			if err != nil {
				return err
			}

			retryCounter := 0
			doRetry := true

			for doRetry {
				doRetry = false

				actionErr := types.RecoverErrorDecorator(func(ctx context.Context) error {
					return action(ctx)
				})(ctx)

				if actionErr == nil {
					err = tx.Commit()

					// if retryable error happens at the commit
					// pq: restart transaction: TransactionRetryWithProtoRefreshError: TransactionRetryError: retry txn (RETRY_SERIALIZABLE - failed preemptive refresh): "sql txn"
					if err != nil && retryCounter < maxRetries {
						doRetry = CheckTransactionRetryable(driver, err)
					}

				} else {
					// if retryable error happens in the action
					// pq: restart transaction: TransactionRetryWithProtoRefreshError: WriteTooOldError: write at ...
					if actionErr != nil && retryCounter < maxRetries {
						doRetry = CheckTransactionRetryable(driver, actionErr)
					}

					err = tx.Rollback()
					if err != nil {
						return err
					}
				}

				if err == nil && actionErr != nil {
					err = actionErr
				}

				if doRetry {
					retryCounter++
					logger.WithContext(ctx).Error(err)
					logger.WithContext(ctx).Infof("retrying transaction: %d", retryCounter)

					// create a new transaction and re inject in the ctx
					tx, err = conn.Beginx()
					if err != nil {
						return err
					}

					ctx = ContextSqlExecutor().Set(ctx, tx)
				}
			}

			return err
		})

		return err
	}
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, action types.ActionFunc) error
}

func ContextTransactionManager() types.ContextKeyAccessor[TransactionManager] {
	return types.NewContextKeyAccessor[TransactionManager](contextKeyTransactionManager)
}

func NewTransactionManager(ctx context.Context) (TransactionManager, error) {
	mgr := ContextTransactionManager().Get(ctx)
	if mgr == nil {
		mgr = new(SqlTransactionManager)
	}
	return mgr, nil
}

type SqlTransactionManager struct{}

func (t SqlTransactionManager) WithTransaction(ctx context.Context, action types.ActionFunc) error {
	wrappedAction := TransactionDecorator(action)
	return wrappedAction(ctx)
}

func CheckTransactionRetryable(driver string, err error) bool {
	isTransactionRetryable := false
	if driver == DriverPostgres {
		code := SqlErrCode(err)
		isTransactionRetryable = (code == "CR000" || code == "40001")
	}

	return isTransactionRetryable
}

func SqlErrCode(err error) (errString string) {
	var sqlErr ErrWithSQLState
	if errors.As(err, &sqlErr) {
		errString = sqlErr.SQLState()
	}
	return
}

type ErrWithSQLState interface {
	SQLState() string
}
