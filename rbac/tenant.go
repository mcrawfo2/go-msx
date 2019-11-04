package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"github.com/pkg/errors"
)

const FmtUserDoesNotHaveTenantAccess = "User does not have access to the tenant: %v"

func HasTenant(ctx context.Context, tenantId string) error {
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

	return errors.Errorf(FmtUserDoesNotHavePermissions, tenantId)
}
