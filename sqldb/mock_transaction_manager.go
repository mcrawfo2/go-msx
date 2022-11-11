// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	context "context"

	types "cto-github.cisco.com/NFV-BU/go-msx/types"
)

type MockTransactionManager struct{}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, action types.ActionFunc) error {
	return action(ctx)
}

func InjectMockTransactionManager(ctx context.Context) context.Context {
	return ContextTransactionManager().Set(ctx, new(MockTransactionManager))
}
