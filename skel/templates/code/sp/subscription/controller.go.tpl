package subscription

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/sp/api"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"errors"
	"github.com/emicklei/go-restful"
)

const (
	pathRoot                = "api/v1/subscriptions"
	paramNameSubscriptionId = "subscriptionId"
)

var (
	managePermission         = webservice.Permissions(rbac.PermissionManageServices)
	pathParamSubscriptionId  = restful.PathParameter(paramNameSubscriptionId, "Subscription Id")
	pathSuffixSubscriptionId = "/{subscriptionId}"
)

type subscriptionController struct {
	subscriptionService   subscriptionServiceApi
	subscriptionConverter subscriptionConverter
}

func (c *subscriptionController) Routes(svc *restful.WebService) {
	svc.ApiVersion("v2")
	tag := webservice.TagDefinition("Subscriptions", "Subscription Controller")
	webservice.Routes(svc,
		tag,
		c.createSubscription,
		c.deleteSubscription,
	)
}

func (c *subscriptionController) createSubscription(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		Request api.SubscriptionCreateRequest `req:"body"`
	}

	return svc.POST("").Operation("createSubscription").
		Doc("API to create subscription").
		Notes("API to create subscription").
		Do(webservice.StandardCreate).
		Do(webservice.ResponsePayload(api.SubscriptionCreateResponse{})).
		Do(managePermission).
		Reads(api.SubscriptionCreateRequest{}).
		Do(webservice.PopulateParams(params{})).
		To(webservice.Controller(
			func(req *restful.Request) (body interface{}, err error) {

				converted, ok := webservice.Params(req).(*params)

				if !ok {
					return api.SubscriptionCreateResponse{}, errors.New("failed to parse body")
				}

				subscription, err := c.subscriptionService.CreateSubscription(req.Request.Context(), converted.Request)
				return c.subscriptionConverter.ToCreateResponse(subscription), err
			}))
}

func (c *subscriptionController) deleteSubscription(svc *restful.WebService) *restful.RouteBuilder {
	type params struct {
		SubscriptionId string `req:"path"`
	}

	return svc.DELETE(pathSuffixSubscriptionId).Operation("deleteSubscription").
		Doc("API to start deletion of subscription").
		Notes("API to start deletion of subscription").
		Do(webservice.StandardDelete).
		Do(managePermission).
		Param(pathParamSubscriptionId).
		Do(webservice.PopulateParams(params{})).
		To(webservice.Controller(
			func(req *restful.Request) (body interface{}, err error) {

				converted, ok := webservice.Params(req).(*params)
				if !ok {
					return nil, errors.New("failed to parse param")
				}

				delErr := c.subscriptionService.DeleteSubscription(req.Request.Context(), converted.SubscriptionId)
				return nil, delErr
			}))
}

func newSubscriptionController(ctx context.Context) webservice.RestController {
	service := controllerFromContext(ctx)
	if service == nil {
		service = &subscriptionController{
			subscriptionService:   newSubscriptionService(ctx),
			subscriptionConverter: subscriptionConverter{},
		}
	}
	return service
}

func init() {
	app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseBefore, func(ctx context.Context) error {
			if svc, err := webservice.WebServerFromContext(ctx).NewService(pathRoot); err != nil {
				return err
			} else {
				newSubscriptionController(ctx).Routes(svc)
			}
			return nil
		})
		return nil
	})
}
