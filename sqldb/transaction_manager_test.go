// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TransactionManager_Commit(t *testing.T) {
	ctx := context.Background()

	tm, err := NewTransactionManager(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tm)

	err = tm.WithTransaction(ctx, func(ctx context.Context) error {
		// do all your db processes in here (preferably prepared)

		// then return nil to commit or an error to rollback
		// return errors.New("some error") // to rollback
		return nil // to commit
	})
	assert.Error(t, err)
}

func Test_TransactionManager_Rollback(t *testing.T) {
	ctx := context.Background()

	tm, err := NewTransactionManager(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tm)

	err = tm.WithTransaction(ctx, func(ctx context.Context) error {
		// do all your db processes in here (preferably prepared)

		// then return nil to commit or an error to rollback
		return errors.New("some error") // to rollback
	})
	assert.Error(t, err)
}
