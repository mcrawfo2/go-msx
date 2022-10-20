// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"github.com/pkg/errors"
	"testing"
)

func Test_TransactionManager_Commit(t *testing.T) {
	ctx := context.Background()
	NewTransactionManager().
		WithTransaction(ctx, func(ctx context.Context) error {
			// do all your db processes in here (preferably prepared)

			// then return nil to commit or an error to rollback
			// return errors.New("some error") // to rollback
			return nil // to commit
		})
}

func Test_TransactionManager_Rollback(t *testing.T) {
	ctx := context.Background()
	NewTransactionManager().
		WithTransaction(ctx, func(ctx context.Context) error {
			// do all your db processes in here (preferably prepared)

			// then return nil to commit or an error to rollback
			return errors.New("some error") // to rollback
		})
}
