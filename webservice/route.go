package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
)

const (
	HeaderNameAuthorization = "Authorization"
)

var (
	HeaderAuthorization *restful.Parameter
	logger              = log.NewLogger("msx.webservice")
)

func init() {
	HeaderAuthorization = restful.
		HeaderParameter(HeaderNameAuthorization, "Authentication token in form 'Bearer {token}'").
		Required(false)
}

func StandardRoute(b *restful.RouteBuilder) {
	StandardAuthenticationRequired(b)
	StandardReturns(b)
}

func StandardAuthenticationRequired(b *restful.RouteBuilder) {
	b.Filter(RequireAuthenticatedFilter).
		Notes("This endpoint is secured").
		Param(HeaderAuthorization)
}

func StandardReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 400, 401, 403))
}

func RequireAuthenticatedFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	securityProvider := SecurityProviderFromContext(req.Request.Context())
	if securityProvider != nil {
		err := securityProvider.Authentication(req)
		if err != nil {
			WriteErrorEnvelope(req, resp, http.StatusUnauthorized, err)
			return
		}
	}

	chain.ProcessFilter(req, resp)
}

func PermissionsFilter(anyOf ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var ctx = req.Request.Context()
		if err := rbac.HasPermission(ctx, anyOf); err != nil {
			logger.WithError(err).WithField("perms", anyOf).Error("Permission denied")
			WriteErrorEnvelope(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

func getParameter(parameter *restful.Parameter, req *restful.Request) (string, error) {
	switch parameter.Kind() {
	case restful.PathParameterKind:
		return req.PathParameter(parameter.Data().Name), nil

	case restful.BodyParameterKind:
		return req.BodyParameter(parameter.Data().Name)

	case restful.QueryParameterKind:
		return req.QueryParameter(parameter.Data().Name), nil

	case restful.HeaderParameterKind:
		return req.HeaderParameter(parameter.Data().Name), nil

	default:
		return "", errors.Errorf("Unsupported parameter type: %v", parameter.Kind())
	}
}

func TenantFilter(parameter *restful.Parameter) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		tenantId, err := getParameter(parameter, req)
		if err != nil {
			WriteErrorEnvelope(req, resp, http.StatusBadRequest, err)
			return
		}

		ctx := req.Request.Context()
		if err := rbac.HasTenant(ctx, tenantId); err != nil {
			logger.WithError(err).WithField("tenant", tenantId).Error("Permission denied")
			WriteErrorEnvelope(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

var DefaultSuccessEnvelope = integration.MsxEnvelope{}

type RouteBuilderFunc func(*restful.RouteBuilder)

func Returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", DefaultSuccessEnvelope)
}
func Returns201(b *restful.RouteBuilder) {
	b.Returns(http.StatusCreated, "Created", DefaultSuccessEnvelope)
}
func Returns204(b *restful.RouteBuilder) {
	b.Returns(http.StatusNoContent, "No Content", DefaultSuccessEnvelope)
}
func Returns400(b *restful.RouteBuilder) {
	b.Returns(http.StatusBadRequest, "Bad Request", nil)
}
func Returns401(b *restful.RouteBuilder) {
	b.Returns(http.StatusUnauthorized, "Not Authorized", nil)
}
func Returns403(b *restful.RouteBuilder) {
	b.Returns(http.StatusForbidden, "Forbidden", nil)
}
func Returns404(b *restful.RouteBuilder) {
	b.Returns(http.StatusNotFound, "Not Found", nil)
}
func Returns409(b *restful.RouteBuilder) {
	b.Returns(http.StatusConflict, "Conflict", nil)
}
func Returns424(b *restful.RouteBuilder) {
	b.Returns(http.StatusFailedDependency, "Failed Dependency", nil)
}
func Returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Internal Server Error", nil)
}
func Returns503(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bad Gateway", nil)
}

func Returns(statuses ...int) RouteBuilderFunc {
	var statusFuncs []RouteBuilderFunc
	for _, status := range statuses {
		switch status {
		case 200:
			statusFuncs = append(statusFuncs, Returns200)
		case 201:
			statusFuncs = append(statusFuncs, Returns201)
		case 204:
			statusFuncs = append(statusFuncs, Returns204)
		case 400:
			statusFuncs = append(statusFuncs, Returns400)
		case 401:
			statusFuncs = append(statusFuncs, Returns401)
		case 404:
			statusFuncs = append(statusFuncs, Returns404)
		case 409:
			statusFuncs = append(statusFuncs, Returns409)
		case 424:
			statusFuncs = append(statusFuncs, Returns424)
		case 500:
			statusFuncs = append(statusFuncs, Returns500)
		case 503:
			statusFuncs = append(statusFuncs, Returns503)
		}
	}

	return func(b *restful.RouteBuilder) {
		for _, statusFunc := range statusFuncs {
			statusFunc(b)
		}
	}
}
