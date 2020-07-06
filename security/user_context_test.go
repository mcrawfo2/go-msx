package security

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
	"testing"
	"time"
)

func TestUserContext_Clone(t *testing.T) {
	tenantId, _ := types.NewUUID()
	expiry := time.Now().Nanosecond()
	issuedAt := time.Now().Nanosecond()

	c := &UserContext{
		UserName:    "userName",
		Roles:       []string{"role1", "role2"},
		TenantId:    tenantId,
		Scopes:      []string{"scope1", "scope2"},
		Authorities: []string{"authority1", "authority2"},
		FirstName:   "firstName",
		LastName:    "lastName",
		Issuer:      "issuer",
		Subject:     "subject",
		Exp:         expiry,
		IssuedAt:    issuedAt,
		Jti:         "jti",
		Email:       "email",
		Token:       "token",
		ClientId:    "clientId",
	}

	if got := c.Clone(); !reflect.DeepEqual(got, c) {
		t.Errorf("Clone() = %v, want %v", got, c)
	}
}
