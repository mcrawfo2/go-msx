package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/audit/auditlog"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/contexttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"net/http"
	"reflect"
	"testing"
)

func TestAcceptReturns(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(AcceptReturns).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(202, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(202)).
		WithRequestPredicate(webservicetest.RequestHasAttribute("DefaultReturnCode", 202)).
		Test(t)
}

func TestConsumesJson(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(ConsumesJson).
		WithRoutePredicate(webservicetest.RouteHasConsumes("application/json")).
		Test(t)
}

func TestConsumesTextPlain(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(ConsumesTextPlain).
		WithRoutePredicate(webservicetest.RouteHasConsumes("text/plain")).
		Test(t)
}

func TestCreateReturns(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(CreateReturns).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(201, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(201)).
		WithRequestPredicate(webservicetest.RequestHasAttribute("DefaultReturnCode", 201)).
		Test(t)
}

func TestDefaultReturns(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Ok",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(func(b *restful.RouteBuilder) {
					b.Do(DefaultReturns(100))
				}).
				WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(100)).
				WithRequestPredicate(webservicetest.RequestHasAttribute("DefaultReturnCode", 100)),
		},
		{
			name: "NotFound",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(func(b *restful.RouteBuilder) {
					b.Do(DefaultReturns(404))
				}).
				WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(404)).
				WithRequestPredicate(webservicetest.RequestHasAttribute("DefaultReturnCode", 404)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestErrorPayload(t *testing.T) {
	payload := types.NewOptionalStringFromString("bob").Ptr()

	new(RouteBuilderTest).
		WithRouteBuilderDo(func(b *restful.RouteBuilder) {
			b.Do(ErrorPayload(payload))
		}).
		WithRequestPredicate(webservicetest.RequestHasAttribute("ErrorPayload", payload)).
		Test(t)
}

func TestNoContentReturns(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(NoContentReturns).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(204, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(204)).
		WithRequestPredicate(webservicetest.RequestHasAttribute("DefaultReturnCode", 204)).
		Test(t)
}

func TestPaginatedResponsePayload(t *testing.T) {
	t.Skipped()
}

func TestParams(t *testing.T) {
	type params struct {
		A *int `req:"query"`
	}

	ts := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "OptionalAbsent",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(PopulateParams(new(params))).
				WithRouteTarget(func(req *restful.Request, resp *restful.Response) {
					_ = Params(req).(*params)
					resp.WriteHeader(222)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(222)),
		},
		{
			name: "OptionalPresent",
			test: new(RouteBuilderTest).
				WithRequestQueryParameter("a", "222").
				WithRouteBuilderDo(PopulateParams(new(params))).
				WithRouteTarget(func(req *restful.Request, resp *restful.Response) {
					args := Params(req).(*params)
					resp.WriteHeader(*args.A)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(222)),
		},
	}

	for _, tt := range ts {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestPermissions(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Permissions("MANAGE_SERVICES", "VIEW_SERVICES")).
		WithRoutePredicate(webservicetest.RouteHasPermissions("MANAGE_SERVICES", "VIEW_SERVICES")...).
		Test(t)
}

func TestPermissionsFilter(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Forbidden",
			test: new(RouteBuilderTest).
				WithContextInjector(securitytest.PermissionInjector("MANAGE_TESTS")).
				WithRouteFilter(PermissionsFilter("MANAGE_SERVICES", "VIEW_SERVICES")).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden)),
		},
		{
			name: "Allowed",
			test: new(RouteBuilderTest).
				WithContextInjector(securitytest.PermissionInjector("MANAGE_SERVICES")).
				WithRouteFilter(PermissionsFilter("MANAGE_SERVICES", "VIEW_SERVICES")).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestPopulateParams(t *testing.T) {
	t.Skipped()
}

func TestProducesJson(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(ProducesJson).
		WithRoutePredicate(webservicetest.RouteHasProduces("application/json")).
		Test(t)
}

func TestProducesTextPlain(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(ProducesTextPlain).
		WithRoutePredicate(webservicetest.RouteHasProduces("text/plain")).
		Test(t)
}

func TestResponsePayload(t *testing.T) {
	payload := types.NewOptionalStringFromString("bob").Ptr()

	new(RouteBuilderTest).
		WithRouteBuilderDo(ResponsePayload(payload)).
		WithRoutePredicate(webservicetest.RouteHasAnyWriteSample()).
		Test(t)
}

func TestResponseRawPayload(t *testing.T) {
	payload := types.NewOptionalStringFromString("bob").Ptr()
	var payloadInterface interface{} = payload

	new(RouteBuilderTest).
		WithRouteBuilderDo(ResponseRawPayload(payloadInterface)).
		WithRoutePredicate(webservicetest.RouteHasWriteSample(payloadInterface)).
		Test(t)
}

func TestResponseTypeName(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "BuiltIn",
			args: args{
				t: reflect.TypeOf("string"),
			},
			want:  "string",
			want1: true,
		},
		{
			name: "Struct",
			args: args{
				t: reflect.TypeOf(RouteParam{}),
			},
			want:  "webservice.RouteParam",
			want1: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ResponseTypeName(tt.args.t)
			if got != tt.want {
				t.Errorf("ResponseTypeName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ResponseTypeName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReturns(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "200",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(200)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(200)),
		},
		{
			name: "201",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(201)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(201)),
		},
		{
			name: "202",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(202)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(202)),
		},
		{
			name: "204",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(204)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(204)),
		},
		{
			name: "400",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(400)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(400)),
		},
		{
			name: "401",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(401)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(401)),
		},
		{
			name: "403",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(403)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(403)),
		},
		{
			name: "404",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(404)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(404)),
		},
		{
			name: "409",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(409)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(409)),
		},
		{
			name: "424",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(424)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(424)),
		},
		{
			name: "500",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(500)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(500)),
		},
		{
			name: "502",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(502)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(502)),
		},
		{
			name: "503",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(503)).
				WithRoutePredicate(webservicetest.RouteHasReturnCode(503)),
		},
		{
			name: "Multiple",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(Returns(200, 404)).
				WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 404)...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestReturns200(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns200).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(200)).
		Test(t)
}

func TestReturns201(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns201).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(201)).
		Test(t)
}

func TestReturns202(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns202).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(202)).
		Test(t)
}

func TestReturns204(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns204).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(204)).
		Test(t)
}

func TestReturns400(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns400).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(400)).
		Test(t)
}

func TestReturns401(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns401).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(401)).
		Test(t)
}

func TestReturns403(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns403).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(403)).
		Test(t)
}

func TestReturns404(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns404).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(404)).
		Test(t)
}

func TestReturns409(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns409).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(409)).
		Test(t)
}

func TestReturns424(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns424).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(424)).
		Test(t)
}

func TestReturns500(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns500).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(500)).
		Test(t)
}

func TestReturns502(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns502).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(502)).
		Test(t)
}

func TestReturns503(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Returns503).
		WithRoutePredicate(webservicetest.RouteHasReturnCode(503)).
		Test(t)
}

func TestRoutes(t *testing.T) {
	t.Skipped()
}

func TestStandardAccept(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardAccept).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(202, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(202)).
		WithRoutePredicate(webservicetest.RouteHasConsumes(MIME_JSON)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestStandardCreate(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardCreate).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(201, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(201)).
		WithRoutePredicate(webservicetest.RouteHasConsumes(MIME_JSON)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestStandardDelete(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardDelete).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(200)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestStandardList(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardList).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(200)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestStandardRetrieve(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardRetrieve).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(200)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestStandardReturns(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardReturns).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(200)).
		Test(t)
}

func TestStandardUpdate(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(StandardUpdate).
		WithRoutePredicate(webservicetest.RouteHasReturnCodes(200, 400, 401, 403)...).
		WithRoutePredicate(webservicetest.RouteHasDefaultReturnCode(200)).
		WithRoutePredicate(webservicetest.RouteHasConsumes(MIME_JSON)).
		WithRoutePredicate(webservicetest.RouteHasProduces(MIME_JSON)).
		Test(t)
}

func TestTagDefinition(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(TagDefinition("bob", "alice")).
		WithRoutePredicate(webservicetest.RouteHasTag("bob")).
		Test(t)
}

func TestTenantFilter(t *testing.T) {
	tenantParam := restful.QueryParameter("tenantId", "Tenant Id")
	tenantId, _ := types.NewUUID()

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Forbidden",
			test: new(RouteBuilderTest).
				WithContextInjector(securitytest.TenantAssignmentInjector()).
				WithRouteFilter(TenantFilter(tenantParam)).
				WithRequestQueryParameter("tenantId", tenantId.String()).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden)),
		},
		{
			name: "Allowed",
			test: new(RouteBuilderTest).
				WithContextInjector(securitytest.TenantAssignmentInjector(tenantId)).
				WithRouteFilter(TenantFilter(tenantParam)).
				WithRequestQueryParameter("tenantId", tenantId.String()).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestValidateParams(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "BadRequest",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(ValidateParams(func(req *restful.Request) (err error) {
					return errors.New("Invalid")
				})).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusBadRequest)),
		},
		{
			name: "Ok",
			test: new(RouteBuilderTest).
				WithRouteBuilderDo(ValidateParams(func(req *restful.Request) (err error) {
					return nil
				})).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func Test_auditContextFilter(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteFilter(auditContextFilter).
		WithRouteTargetReturn(http.StatusOK).
		WithContextInjector(func(ctx context.Context) context.Context {
			return ContextWithWebServerValue(ctx, &WebServer{
				cfg: &WebServerConfig{},
			})
		}).
		WithContextPredicate(contexttest.ContextPredicate{
			Description: "Context has auditContext",
			Matches: func(ctx context.Context) bool {
				auditContext := auditlog.RequestAuditFromContext(ctx)
				return auditContext != nil
			},
		}).
		Test(t)
}

func Test_authenticationFilter(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Unauthorized",
			test: new(RouteBuilderTest).
				WithContextInjector(func(ctx context.Context) context.Context {
					mockAuthenticationProvider := new(MockAuthenticationProvider)
					mockAuthenticationProvider.
						On("Authenticate", mock.AnythingOfType("*restful.Request")).
						Return(errors.New("authentication failed"))

					return ContextWithSecurityProvider(ctx, mockAuthenticationProvider)
				}).
				WithRouteFilter(authenticationFilter).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusUnauthorized)),
		},
		{
			name: "Authorized",
			test: new(RouteBuilderTest).
				WithContextInjector(func(ctx context.Context) context.Context {
					mockAuthenticationProvider := new(MockAuthenticationProvider)
					mockAuthenticationProvider.
						On("Authenticate", mock.AnythingOfType("*restful.Request")).
						Return(nil)

					return ContextWithSecurityProvider(ctx, mockAuthenticationProvider)
				}).
				WithRouteFilter(authenticationFilter).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func Test_getParameter(t *testing.T) {
	t.Skipped()
}

func Test_newGenericResponse(t *testing.T) {
	t.Skipped()
}

func Test_requestValidator_Validate(t *testing.T) {
	type fields struct {
		req *restful.Request
		fn  ValidatorFunction
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				req: nil,
				fn: func(req *restful.Request) (err error) {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "BadRequest",
			fields: fields{
				req: nil,
				fn: func(req *restful.Request) (err error) {
					return errors.New("something wrong")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := requestValidator{
				req: tt.fields.req,
				fn:  tt.fields.fn,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_securityContextFilter(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoToken",
			test: new(RouteBuilderTest).
				WithRouteFilter(tokenUserContextFilter).
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)).
				WithContextPredicate(contexttest.ContextHasNamedUserContext("anonymous")),
		},
		{
			name: "MalformedHeader",
			test: new(RouteBuilderTest).
				WithRouteFilter(tokenUserContextFilter).
				WithRequestHeader("Authorization", "malformed").
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusUnauthorized)),
		},
		{
			name: "BearerToken",
			test: new(RouteBuilderTest).
				WithRouteFilter(tokenUserContextFilter).
				WithRouteBuilderDo(func(_ *restful.RouteBuilder) {
					mockTokenProvider := new(security.MockTokenProvider)
					mockTokenProvider.
						On("UserContextFromToken", mock.AnythingOfType("*context.valueCtx"), "abc123").
						Return(&security.UserContext{
							UserName: "tester",
							Roles:    []string{"TESTER"},
							Token:    "abc123",
						}, nil)
					security.SetTokenProvider(mockTokenProvider)
				}).
				WithRequestHeader("Authorization", "Bearer abc123").
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)).
				WithContextPredicate(contexttest.ContextHasNamedUserContext("tester")),
		},
		{
			name: "BadBearerToken",
			test: new(RouteBuilderTest).
				WithRouteFilter(tokenUserContextFilter).
				WithRouteBuilderDo(func(_ *restful.RouteBuilder) {
					mockTokenProvider := new(security.MockTokenProvider)
					mockTokenProvider.
						On("UserContextFromToken", mock.AnythingOfType("*context.valueCtx"), "abc123").
						Return(nil, errors.New("bad token"))
					security.SetTokenProvider(mockTokenProvider)
				}).
				WithRequestHeader("Authorization", "Bearer abc123").
				WithRouteTargetReturn(http.StatusOK).
				WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusUnauthorized)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}
