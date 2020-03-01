package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"github.com/pkg/errors"
)

const (
	PermissionIsApiAdmin            = "IS_API_ADMIN"
	PermissionViewServices          = "VIEW_SERVICES"
	PermissionManageServices        = "MANAGE_SERVICES"
	PermissionViewContact           = "VIEW_CONTACT"
	PermissionManageContact         = "MANAGE_CONTACT"
	PermissionViewLocaleString      = "VIEW_LOCALE_STRING"
	PermissionManageLocaleString    = "MANAGE_LOCALE_STRING"
	PermissionViewIntegration       = "VIEW_INTEGRATION"
	PermissionManageIntegration     = "MANAGE_INTEGRATION"
	PermissionViewMaintenanceInfo   = "VIEW_MAINTENANCE_INFO"
	PermissionManageMaintenanceInfo = "MANAGE_MAINTENANCE_INFO"
)

var ErrUserDosNotHavePermission = errors.New("User does not have any of the required permissions")

func HasPermission(ctx context.Context, required []string) error {
	usermanagementIntegration, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return err
	}

	s, err := usermanagementIntegration.GetMyCapabilities()
	if err != nil {
		return err
	}

	payload, ok := s.Payload.(*usermanagement.UserCapabilityListResponse)
	if !ok {
		return errors.New("Failed to convert response payload")
	}

	for _, p := range required {
		for _, c := range payload.Capabilities {
			if p == c.Name {
				return nil
			}
		}
	}

	return ErrUserDosNotHavePermission
}
