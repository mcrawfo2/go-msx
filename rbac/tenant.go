package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

var logger = log.NewLogger("msx.rbac")
var ErrUserDoesNotHaveTenantAccess = errors.New("User does not have access to the tenant")

func HasTenant(ctx context.Context, tenantId string) error {
	tenantUuid, err := types.ParseUUID(tenantId)
	if err != nil {
		return err
	}

	userContext := security.UserContextFromContext(ctx)
	if userContext != nil {
		if userContext.TenantId == tenantId {
			userTenantUuid, err := types.ParseUUID(tenantId)
			if err != nil {
				return err
			}
			if userTenantUuid.Equals(tenantUuid) {
				return nil
			}
		}
	}

	usermanagementIntegration, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return err
	}

	s, err := usermanagementIntegration.GetMyTenantIds()
	if err != nil {
		return err
	}

	payload, ok := s.Payload.(*usermanagement.TenantIdList)
	if !ok {
		return errors.New("Failed to convert response payload")
	}

	for _, t := range *payload {
		if tenantId == t {
			return nil
		}
	}

	return ErrUserDoesNotHaveTenantAccess
}
