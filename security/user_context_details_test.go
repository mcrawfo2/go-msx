// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package security

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserContextDetails_HasTenantId(t *testing.T) {
	tenantId, _ := types.NewUUID()
	notTenantId := types.EmptyUUID()
	userContextDetails := UserContextDetails{
		Tenants: []types.UUID{tenantId},
	}

	assert.True(t, userContextDetails.HasTenantId(tenantId))
	assert.False(t, userContextDetails.HasTenantId(notTenantId))
}
