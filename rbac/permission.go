package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"github.com/pkg/errors"
)

const (
	PermissionIsApiAdmin = "IS_API_ADMIN"
	PermissionViewServices = "VIEW_SERVICES"
	PermissionManageServices = "MANAGE_SERVICES"
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
