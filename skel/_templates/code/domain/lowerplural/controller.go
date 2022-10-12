package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	//#if TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	//#endif TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/domain/api"
	//#if TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	//#endif TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
)

const (
	pathRoot                          = "/api"
	pathPrefixUpperCamelPlural        = "/v1/lowerplural"
	pathSuffixUpperCamelSingularId    = "/{lowerCamelSingularId}"
	pathParamNameUpperCamelSingularId = "lowerCamelSingularId"
	//#if TENANT_DOMAIN
	queryParamNameTenantId = "tenantId"
	//#endif TENANT_DOMAIN
)

var (
	viewPermission            = webservice.Permissions("VIEW_SCREAMING_SNAKE_PLURAL")
	managePermission          = webservice.Permissions("MANAGE_SCREAMING_SNAKE_PLURAL")
	paramUpperCamelSingularId = restful.PathParameter(pathParamNameUpperCamelSingularId, "Title Singular Id")
)

type lowerCamelSingularController struct {
	lowerCamelSingularService lowerCamelSingularServiceApi
}

func (c *lowerCamelSingularController) Routes(svc *restful.WebService) {
	svc.ApiVersion("v2")
	tag := webservice.TagDefinition("Title Singular", "Title Singular Controller")
	webservice.Routes(svc, tag,
		c.listUpperCamelPlural,
		c.getUpperCamelSingular,
		c.createUpperCamelSingular,
		c.updateUpperCamelSingular,
		c.deleteUpperCamelSingular)
}

func (c *lowerCamelSingularController) listUpperCamelPlural(svc *restful.WebService) *restful.RouteBuilder {
	//#if TENANT_DOMAIN
	var paramTenantId = restful.QueryParameter(queryParamNameTenantId, "Tenant Id").Required(true)

	type params struct {
		TenantId types.UUID `req:"query"`
	}
	//#endif TENANT_DOMAIN

	return svc.GET(pathPrefixUpperCamelPlural).
		Operation("listUpperCamelPlural").
		Doc("List all the Title Plural").
		Do(webservice.StandardList).
		Do(webservice.ResponseRawPayload([]api.UpperCamelSingularResponse{})).
		//#if TENANT_DOMAIN
		Param(paramTenantId).
		Do(webservice.PopulateParams(params{})).
		//#endif TENANT_DOMAIN
		Do(viewPermission).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				//#if TENANT_DOMAIN
				var args = webservice.Params(req).(*params)
				//#endif TENANT_DOMAIN

				body, err = c.lowerCamelSingularService.ListUpperCamelPlural(
					req.Request.Context(),
					//#if TENANT_DOMAIN
					args.TenantId,
					//#endif TENANT_DOMAIN
				)

				//#if TENANT_DOMAIN
				if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
				} else if err != nil {
					return nil, err
				}
				//#else TENANT_DOMAIN
				if err != nil {
					return nil, err
				}
				//#endif TENANT_DOMAIN

				return
			}))
}

func (c *lowerCamelSingularController) getUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularId types.UUID `req:"path"`
	}

	return svc.GET(pathPrefixUpperCamelPlural + pathSuffixUpperCamelSingularId).
		Operation("getUpperCamelSingular").
		Doc("Retrieve the specified Title Singular").
		Do(webservice.StandardRetrieve).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Param(paramUpperCamelSingularId).
		Do(webservice.PopulateParams(params{})).
		Do(viewPermission).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				var args = webservice.Params(req).(*params)

				body, err = c.lowerCamelSingularService.GetUpperCamelSingular(
					req.Request.Context(),
					args.UpperCamelSingularId)
				if err == lowerCamelSingularErrNotFound {
					return nil, webservice.NewNotFoundError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return
			}))
}

func (c *lowerCamelSingularController) createUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		Request api.UpperCamelSingularCreateRequest `req:"body,san"`
	}

	return svc.POST(pathPrefixUpperCamelPlural).
		Operation("createUpperCamelSingular").
		Doc("Create a new Title Singular").
		Do(webservice.StandardCreate).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Reads(api.UpperCamelSingularCreateRequest{}).
		Do(managePermission).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				args := webservice.Params(req).(*params)

				body, err = c.lowerCamelSingularService.CreateUpperCamelSingular(
					req.Request.Context(),
					args.Request)

				if err == lowerCamelSingularErrAlreadyExists {
					return nil, webservice.NewConflictError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return
			}))
}

func (c *lowerCamelSingularController) updateUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularId types.UUID                          `req:"path"`
		Request              api.UpperCamelSingularUpdateRequest `req:"body,san"`
	}

	return svc.PUT(pathPrefixUpperCamelPlural + pathSuffixUpperCamelSingularId).
		Operation("updateUpperCamelSingular").
		Doc("Update the specified Title Singular").
		Do(webservice.StandardUpdate).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Param(paramUpperCamelSingularId).
		Reads(api.UpperCamelSingularUpdateRequest{}).
		Do(managePermission).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				args := webservice.Params(req).(*params)

				body, err = c.lowerCamelSingularService.UpdateUpperCamelSingular(
					req.Request.Context(),
					args.UpperCamelSingularId,
					args.Request)

				if err == lowerCamelSingularErrNotFound {
					return nil, webservice.NewNotFoundError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return
			}))
}

func (c *lowerCamelSingularController) deleteUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularId types.UUID `req:"path"`
	}

	return svc.DELETE(pathPrefixUpperCamelPlural + pathSuffixUpperCamelSingularId).
		Operation("deleteUpperCamelSingular").
		Doc("Delete the specified Title Singular").
		Do(webservice.StandardDelete).
		Do(webservice.ResponseRawPayload(types.Empty{})).
		Param(paramUpperCamelSingularId).
		Do(managePermission).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				args := webservice.Params(req).(*params)

				err = c.lowerCamelSingularService.DeleteUpperCamelSingular(
					req.Request.Context(),
					args.UpperCamelSingularId)
				//#if TENANT_DOMAIN
				if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
				}
				//#endif TENANT_DOMAIN
				if err != nil {
					return nil, err
				}

				return
			}))
}

func newUpperCamelSingularController(ctx context.Context) (webservice.RestController, error) {
	controller := lowerCamelSingularControllerFromContext(ctx)
	if controller == nil {
		lowerCamelSingularService, err := newUpperCamelSingularService(ctx)
		if err != nil {
			return nil, err
		}

		controller = &lowerCamelSingularController{
			lowerCamelSingularService: lowerCamelSingularService,
		}
	}
	return controller, nil
}

func init() {
	app.OnRootEvent(app.EventStart, app.PhaseBefore, func(ctx context.Context) error {
		controller, err := newUpperCamelSingularController(ctx)
		if err != nil {
			return err
		}

		return webservice.
			WebServerFromContext(ctx).
			RegisterRestController(pathRoot, controller)
	})
}
