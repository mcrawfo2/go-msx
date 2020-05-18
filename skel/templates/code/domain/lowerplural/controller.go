package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	//#if TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	//#endif TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/domain/api"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
)

const (
	pathRoot                            = "/api/v1/lowerplural"
	pathSuffixUpperCamelSingularName    = "/{lowerCamelSingularName}"
	pathParamNameUpperCamelSingularName = "lowerCamelSingularName"
)

var (
	viewPermissionFilter        = webservice.PermissionsFilter("VIEW_SCREAMING_SNAKE_PLURAL")
	managePermissionFilter      = webservice.PermissionsFilter("MANAGE_SCREAMING_SNAKE_PLURAL")
	paramUpperCamelSingularName = restful.PathParameter(pathParamNameUpperCamelSingularName, "Title Singular Name")
)

type lowerCamelSingularController struct {
	lowerCamelSingularService   lowerCamelSingularServiceApi
	lowerCamelSingularConverter lowerCamelSingularConverter
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
	return svc.GET("").
		Operation("listUpperCamelPlural").
		Doc("List all the Title Plural").
		Do(webservice.StandardList).
		Do(webservice.ResponseRawPayload([]api.UpperCamelSingularResponse{})).
		Filter(viewPermissionFilter).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				lowerCamelPlural, err := c.lowerCamelSingularService.ListUpperCamelPlural(req.Request.Context())
				if err != nil {
					return nil, err
				}

				return c.lowerCamelSingularConverter.ToUpperCamelSingularListResponse(lowerCamelPlural), nil
			}))
}

func (c *lowerCamelSingularController) getUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularName string `req:"path"`
	}

	return svc.GET(pathSuffixUpperCamelSingularName).
		Operation("getUpperCamelSingular").
		Doc("Retrieve the specified Title Singular").
		Do(webservice.StandardRetrieve).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Param(paramUpperCamelSingularName).
		Do(webservice.PopulateParams(params{})).
		Filter(viewPermissionFilter).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				var params = webservice.Params(req).(*params)

				lowerCamelSingular, err := c.lowerCamelSingularService.GetUpperCamelSingular(
					req.Request.Context(),
					params.UpperCamelSingularName)
				if err == lowerCamelSingularErrNotFound {
					return nil, webservice.NewNotFoundError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return c.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular), nil
			}))
}

func (c *lowerCamelSingularController) createUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		Request api.UpperCamelSingularCreateRequest `req:"body"`
	}

	return svc.POST("").
		Operation("createUpperCamelSingular").
		Doc("Create a new Title Singular").
		Do(webservice.StandardCreate).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Reads(api.UpperCamelSingularCreateRequest{}).
		Filter(managePermissionFilter).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				params := webservice.Params(req).(*params)

				lowerCamelSingular, err := c.lowerCamelSingularService.CreateUpperCamelSingular(
					req.Request.Context(),
					params.Request)

				if err == lowerCamelSingularErrAlreadyExists {
					return nil, webservice.NewConflictError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return c.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular), nil
			}))
}

func (c *lowerCamelSingularController) updateUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularName string                              `req:"path"`
		Request                api.UpperCamelSingularUpdateRequest `req:"body"`
	}

	return svc.PUT(pathSuffixUpperCamelSingularName).
		Operation("updateUpperCamelSingular").
		Doc("Update the specified Title Singular").
		Do(webservice.StandardUpdate).
		Do(webservice.ResponseRawPayload(api.UpperCamelSingularResponse{})).
		Param(paramUpperCamelSingularName).
		Reads(api.UpperCamelSingularUpdateRequest{}).
		Filter(managePermissionFilter).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				params := webservice.Params(req).(*params)

				lowerCamelSingular, err := c.lowerCamelSingularService.UpdateUpperCamelSingular(
					req.Request.Context(),
					params.UpperCamelSingularName,
					params.Request)

				if err == lowerCamelSingularErrNotFound {
					return nil, webservice.NewNotFoundError(err)
					//#if TENANT_DOMAIN
				} else if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
					//#endif TENANT_DOMAIN
				} else if err != nil {
					return nil, err
				}

				return c.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular), nil
			}))
}

func (c *lowerCamelSingularController) deleteUpperCamelSingular(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		UpperCamelSingularName string `req:"path"`
	}

	return svc.DELETE(pathSuffixUpperCamelSingularName).
		Operation("deleteUpperCamelSingular").
		Doc("Delete the specified Title Singular").
		Do(webservice.StandardDelete).
		Do(webservice.ResponseRawPayload(struct{}{})).
		Param(paramUpperCamelSingularName).
		Filter(managePermissionFilter).
		Do(webservice.PopulateParams(params{})).
		To(webservice.RawController(
			func(req *restful.Request) (body interface{}, err error) {
				params := webservice.Params(req).(*params)

				err = c.lowerCamelSingularService.DeleteUpperCamelSingular(
					req.Request.Context(),
					params.UpperCamelSingularName)
				//#if TENANT_DOMAIN
				if err == rbac.ErrUserDoesNotHaveTenantAccess {
					return nil, webservice.NewForbiddenError(err)
				}
				//#endif TENANT_DOMAIN
				if err != nil {
					return nil, err
				}

				return nil, nil
			}))
}

func init() {
	app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseBefore, func(ctx context.Context) error {
			controller := &lowerCamelSingularController{
				lowerCamelSingularService:   newUpperCamelSingularService(ctx),
				lowerCamelSingularConverter: lowerCamelSingularConverter{},
			}

			return webservice.
				WebServerFromContext(ctx).
				RegisterRestController(pathRoot, controller)
		})
		return nil
	})
}
