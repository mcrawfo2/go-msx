package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
)

const (
	FmtUserDoesNotHavePermissions = "User does not have permissions to use this api: %v"
)

func PermissionsFilter(anyOf ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var ctx = req.Request.Context()
		if err := HasPermission(ctx, anyOf); err != nil {
			if err = WriteErrorEnvelope(req, resp, http.StatusUnauthorized, err); err != nil {
				logger.WithError(err).Error("Failed to write error envelope")
			}
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

func HasPermission(ctx context.Context, required []string) error {
	usermanagementIntegration, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return err
	}

	s, err := usermanagementIntegration.GetMyCapabilities()
	if err != nil {
		return err
	}

	payload, ok := s.Payload.(*usermanagement.UserCapabilityListDTO)
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

	return errors.Errorf(FmtUserDoesNotHavePermissions, required)
}
