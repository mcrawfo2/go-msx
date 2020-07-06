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
