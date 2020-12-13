package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/audit/auditlog"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

const (
	HeaderNameAuthorization     = "Authorization"
	MetadataKeyResponseEnvelope = "MSX_RESPONSE_ENVELOPE"
	MetadataKeyResponsePayload  = "MSX_RESPONSE_PAYLOAD"
	MetadataTagDefinition       = "TagDefinition"
	MetadataPermissions         = "Permissions"
	AttributeDefaultReturnCode  = "DefaultReturnCode"
	AttributeErrorPayload       = "ErrorPayload"
	AttributeError              = "Error"
	AttributeParams             = "Params"
	AttributeStandard           = "Standard"
	AttributeParamsValidator    = "ParamsValidator"
)

var (
	logger = log.NewLogger("msx.webservice")
)

func StandardList(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func StandardRetrieve(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func StandardCreate(b *restful.RouteBuilder) {
	b.Do(CreateReturns, ProducesJson, ConsumesJson)
}

func StandardUpdate(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson, ConsumesJson)
}

func StandardDelete(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func StandardAccept(b *restful.RouteBuilder) {
	b.Do(AcceptReturns, ProducesJson, ConsumesJson)
}

func ResponseTypeName(t reflect.Type) (string, bool) {
	typeName := types.GetTypeName(t, true)
	return typeName, typeName != ""
}

func newGenericResponse(structType reflect.Type, structFieldName string, payloadInstance interface{}) interface{} {
	responseType := types.NewParameterizedStruct(
		structType,
		structFieldName,
		payloadInstance)

	return reflect.New(responseType).Interface()
}

func ResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.MsxEnvelope))
	return func(b *restful.RouteBuilder) {
		example := newGenericResponse(
			reflect.TypeOf(integration.MsxEnvelope{}),
			"Payload",
			payload)
		b.DefaultReturns("Success", example)
		b.Writes(example)
		b.Do(errorPayloadFn)
	}
}

func PaginatedResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.MsxEnvelope))
	return func(b *restful.RouteBuilder) {
		paginatedPayload := newGenericResponse(
			reflect.TypeOf(paging.PaginatedResponse{}),
			"Content",
			payload)
		envelopedPayload := newGenericResponse(
			reflect.TypeOf(integration.MsxEnvelope{}),
			"Payload",
			paginatedPayload)
		b.DefaultReturns("Success", envelopedPayload)
		b.Writes(envelopedPayload)
		b.Do(errorPayloadFn)
	}
}

func ResponseRawPayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.ErrorDTO))
	return func(b *restful.RouteBuilder) {
		b.DefaultReturns("Success", payload)
		if payload != nil {
			b.Writes(payload)
		}
		b.Do(errorPayloadFn)
	}
}

func ErrorPayload(payload interface{}) func(*restful.RouteBuilder) {
	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			request.SetAttribute(AttributeErrorPayload, payload)
			chain.ProcessFilter(request, response)
		})
	}
}

func StandardReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 400, 401, 403), DefaultReturns(200))
}

func CreateReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 201, 400, 401, 403), DefaultReturns(201))
}

func AcceptReturns(b *restful.RouteBuilder) {
	b.Do(Returns(202, 400, 401, 403), DefaultReturns(202))
}

func NoContentReturns(b *restful.RouteBuilder) {
	b.Do(Returns(204, 400, 401, 403), DefaultReturns(204))
}

func ProducesJson(b *restful.RouteBuilder) {
	b.Produces(MIME_JSON)
}

func ConsumesJson(b *restful.RouteBuilder) {
	b.Consumes(MIME_JSON)
}

func ProducesTextPlain(b *restful.RouteBuilder) {
	b.Produces(MIME_TEXT_PLAIN)
}

func ConsumesTextPlain(b *restful.RouteBuilder) {
	b.Consumes(MIME_TEXT_PLAIN)
}

func DefaultReturns(code int) RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		b.Metadata(AttributeDefaultReturnCode, code)
		b.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			request.SetAttribute(AttributeDefaultReturnCode, code)
			chain.ProcessFilter(request, response)
		})
	}
}

func securityContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := httprequest.ExtractToken(req.Request)
	if err != nil && err != httprequest.ErrNotFound {
		WriteError(req, resp, http.StatusUnauthorized, err)
		return
	}

	if err == nil {
		userContext, err := security.NewUserContextFromToken(req.Request.Context(), token)
		if err != nil {
			WriteError(req, resp, http.StatusUnauthorized, err)
			return
		}

		ctx := security.ContextWithUserContext(req.Request.Context(), userContext)
		req.Request = req.Request.WithContext(ctx)
	}

	chain.ProcessFilter(req, resp)
}

func authenticationFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authenticationProvider := AuthenticationProviderFromContext(req.Request.Context())
	if authenticationProvider != nil {
		err := authenticationProvider.Authenticate(req)
		if err != nil {
			WriteError(req, resp, http.StatusUnauthorized, err)
			return
		}
	}

	chain.ProcessFilter(req, resp)
}

func auditContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	server := WebServerFromContext(req.Request.Context())
	auditDetails := auditlog.ExtractRequestDetails(req, server.cfg.Host, server.cfg.Port)
	ctx := auditlog.ContextWithRequestDetails(req.Request.Context(), auditDetails)
	req.Request = req.Request.WithContext(ctx)
	chain.ProcessFilter(req, resp)
}

func Permissions(anyOf ...string) RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		b.Metadata(MetadataPermissions, anyOf)
		b.Filter(PermissionsFilter(anyOf...))
	}
}

// Deprecated
func PermissionsFilter(anyOf ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var ctx = req.Request.Context()
		// Temporarily allow system user
		var userContext = security.UserContextFromContext(ctx)
		if userContext.UserName != "system" {
			if err := rbac.HasPermission(ctx, anyOf); err != nil {
				logger.WithError(err).WithField("perms", anyOf).Error("Permission denied")
				WriteError(req, resp, http.StatusForbidden, err)
				return
			}
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
			WriteError(req, resp, http.StatusBadRequest, err)
			return
		}

		tenantUuid, err := types.ParseUUID(tenantId)
		if err != nil {
			WriteError(req, resp, http.StatusBadRequest, err)
		}

		ctx := req.Request.Context()
		if err := rbac.HasTenant(ctx, tenantUuid); err != nil {
			ctx = log.ExtendContext(ctx, log.LogContext{
				"tenant": tenantId,
			})
			req.Request = req.Request.WithContext(ctx)
			WriteError(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

type RouteBuilderFunc func(*restful.RouteBuilder)

func Returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", nil)
}
func Returns201(b *restful.RouteBuilder) {
	b.Returns(http.StatusCreated, "Created", nil)
}
func Returns202(b *restful.RouteBuilder) {
	b.Returns(http.StatusAccepted, "Accepted", nil)
}
func Returns204(b *restful.RouteBuilder) {
	b.Returns(http.StatusNoContent, "No Content", nil)
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
func Returns502(b *restful.RouteBuilder) {
	b.Returns(http.StatusBadGateway, "Bad Gateway", nil)
}
func Returns503(b *restful.RouteBuilder) {
	b.Returns(http.StatusServiceUnavailable, "Service Unavailable", nil)
}

func Returns(statuses ...int) RouteBuilderFunc {
	var statusFuncs []RouteBuilderFunc
	for _, status := range statuses {
		switch status {
		case 200:
			statusFuncs = append(statusFuncs, Returns200)
		case 201:
			statusFuncs = append(statusFuncs, Returns201)
		case 202:
			statusFuncs = append(statusFuncs, Returns202)
		case 204:
			statusFuncs = append(statusFuncs, Returns204)
		case 400:
			statusFuncs = append(statusFuncs, Returns400)
		case 401:
			statusFuncs = append(statusFuncs, Returns401)
		case 403:
			statusFuncs = append(statusFuncs, Returns403)
		case 404:
			statusFuncs = append(statusFuncs, Returns404)
		case 409:
			statusFuncs = append(statusFuncs, Returns409)
		case 424:
			statusFuncs = append(statusFuncs, Returns424)
		case 500:
			statusFuncs = append(statusFuncs, Returns500)
		case 502:
			statusFuncs = append(statusFuncs, Returns502)
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

func TagDefinition(name, description string) RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		b.Metadata(restfulspec.KeyOpenAPITags, []string{name})
		b.Metadata(MetadataTagDefinition, spec.TagProps{
			Name:        name,
			Description: description,
		})
	}
}

type RestController interface {
	Routes(svc *restful.WebService)
}

type RouteFunction func(svc *restful.WebService) *restful.RouteBuilder

func PopulateParams(template interface{}) RouteBuilderFunc {
	templateType := reflect.TypeOf(template)
	if templateType.Kind() == reflect.Ptr {
		templateType = templateType.Elem()
	}

	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			// Instantiate a new object of the same type as template
			target := reflect.New(templateType).Interface()

			// Populate the target
			if err := Populate(req, target); err != nil {
				WriteError(req, response, 400, err)
				return
			}

			req.SetAttribute(AttributeParams, target)

			chain.ProcessFilter(req, response)
		})
	}
}

func ValidateParams(fn ValidatorFunction) RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			err := validate.Validate(requestValidator{fn: fn, req: req})
			if err != nil {
				WriteError(req, response, 400, err)
				return
			}

			chain.ProcessFilter(req, response)
		})
	}
}

type ValidatorFunction func(req *restful.Request) (err error)

type requestValidator struct {
	req *restful.Request
	fn  ValidatorFunction
}

func (r requestValidator) Validate() error {
	return r.fn(r.req)
}

func Params(req *restful.Request) interface{} {
	return req.Attribute(AttributeParams)
}

func Routes(svc *restful.WebService, tag RouteBuilderFunc, routeFunctions ...RouteFunction) {
	for _, routeFunction := range routeFunctions {
		routeBuilder := routeFunction(svc)
		routeBuilder.Do(tag)
		svc.Route(routeBuilder)
	}
}
