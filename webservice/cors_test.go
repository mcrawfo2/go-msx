package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestActivateCors(t *testing.T) {
	container := testCorsContainer()
	corsConfig := CorsConfig{Enabled: true}
	ActivateCors(container, corsConfig)

	new(RouteBuilderTest).
		WithRequestMethod("OPTIONS").
		WithRequestPath("/api/v1/service").
		WithRequestHeader("Origin", "localhost").
		WithRequestHeader("Access-Control-Request-Method", "GET").
		WithRouteTarget(HttpHandlerController(container.ServeHTTP)).
		WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Origin")).
		WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Method")).
		WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Headers")).
		WithResponsePredicate(webservicetest.ResponseHasHeader("Access-Control-Allow-Methods", "PATCH,POST,GET,PUT,DELETE,HEAD,OPTIONS,TRACE")).
		WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "GET,OPTIONS,POST")).
		WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
		Test(t)
}

func Test_corsFilter(t *testing.T) {
	container := testCorsContainer()
	corsConfig := CorsConfig{Enabled: true}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Root",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/service").
				WithRequestHeader("Origin", "localhost").
				WithRequestHeader("Access-Control-Request-Method", "GET").
				WithRouteFilter(corsFilter(container, newCors(container, corsConfig))).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Origin")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Method")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Headers")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Access-Control-Allow-Methods", "PATCH,POST,GET,PUT,DELETE,HEAD,OPTIONS,TRACE")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "GET,OPTIONS,POST")).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "Parameter",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/service/my-service-id").
				WithRequestHeader("Origin", "localhost").
				WithRequestHeader("Access-Control-Request-Method", "GET").
				WithRouteFilter(corsFilter(container, newCors(container, corsConfig))).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Origin")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Method")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Vary", "Access-Control-Request-Headers")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Access-Control-Allow-Methods", "PATCH,POST,GET,PUT,DELETE,HEAD,OPTIONS,TRACE")).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "DELETE,GET,OPTIONS,PUT")).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "NoEndpoints",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/services").
				WithRequestHeader("Origin", "localhost").
				WithRequestHeader("Access-Control-Request-Method", "GET").
				WithRouteFilter(corsFilter(container, newCors(container, corsConfig))).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "NonOptions",
			test: new(RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestPath("/api/v1/services").
				WithRequestHeader("Origin", "localhost").
				WithRouteTargetReturn(201).
				WithRouteFilter(corsFilter(container, newCors(container, corsConfig))).
				WithResponsePredicate(webservicetest.ResponseHasStatus(201)),
		},
		{
			name: "NonCors",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/service/my-service-id").
				WithRequestHeader("Origin", "localhost").
				WithRouteFilter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
					optionsHandler(container, request, response)
				}).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "DELETE,GET,OPTIONS,PUT")).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func Test_findAllowedMethods(t *testing.T) {
	container := testCorsContainer()

	type args struct {
		requested string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Root",
			args: args{
				requested: "/api/v1/service",
			},
			want: []string{"GET", "OPTIONS", "POST"},
		},
		{
			name: "Parameter",
			args: args{
				requested: "/api/v1/service/my-service-id",
			},
			want: []string{"DELETE", "GET", "OPTIONS", "PUT"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAllowedMethods(container, tt.args.requested); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllowedMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAllowHeader(t *testing.T) {
	type args struct {
		allowedMethodsResult string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple",
			args: args{
				allowedMethodsResult: "A,B,C,X,Y,Z",
			},
			want: "A,B,C,OPTIONS,X,Y,Z",
		},
		{
			name: "Duplicate",
			args: args{
				allowedMethodsResult: "GET,POST,PUT,PUT,POST,GET",
			},
			want: "GET,OPTIONS,POST,PUT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAllowHeader(tt.args.allowedMethodsResult); got != tt.want {
				t.Errorf("getAllowHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchesPath(t *testing.T) {
	type args struct {
		requested []string
		available []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ConstantMatch",
			args: args{
				requested: []string{"api", "v1", "dummy"},
				available: []string{"api", "v1", "dummy"},
			},
			want: true,
		},
		{
			name: "ParameterMatch",
			args: args{
				requested: []string{"api", "v1", "dummy", "my-service-id"},
				available: []string{"api", "v1", "dummy", "{serviceId}"},
			},
			want: true,
		},
		{
			name: "NoMatch",
			args: args{
				requested: []string{"api", "v1", "dummy"},
				available: []string{"api", "v1", "nomatch"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesPath(tt.args.requested, tt.args.available); got != tt.want {
				t.Errorf("matchesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCors(t *testing.T) {
	container := restful.NewContainer()
	corsConfig := CorsConfig{Enabled: true}
	cors := newCors(container, corsConfig)
	assert.NotEmpty(t, cors.AllowedDomains)
	assert.NotEmpty(t, cors.AllowedHeaders)
	assert.Equal(t, container, cors.Container)
}

func Test_optionsHandler(t *testing.T) {
	container := testCorsContainer()

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Root",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/service").
				WithRouteFilter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
					optionsHandler(container, request, response)
				}).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "GET,OPTIONS,POST")).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
		{
			name: "Parameter",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/service/my-service-id").
				WithRouteFilter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
					optionsHandler(container, request, response)
				}).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Allow", "DELETE,GET,OPTIONS,PUT")).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
		{
			name: "NoEndpoints",
			test: new(RouteBuilderTest).
				WithRequestMethod("OPTIONS").
				WithRequestPath("/api/v1/services").
				WithRequestHeader("Origin", "localhost").
				WithRouteFilter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
					optionsHandler(container, request, response)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(404)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func Test_tokenizePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Root",
			args: args{
				path: "/",
			},
			want: nil,
		},
		{
			name: "LeadingSlash",
			args: args{
				path: "/api/v1/simple",
			},
			want: []string{"api", "v1", "simple"},
		},
		{
			name: "TrailingSlash",
			args: args{
				path: "/api/v1/simple/",
			},
			want: []string{"api", "v1", "simple"},
		},
		{
			name: "PathParameter",
			args: args{
				path: "api/v1/simple/{id}",
			},
			want: []string{"api", "v1", "simple", "{id}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenizePath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testCorsContainer() *restful.Container {
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Path("/api/v1/service")
	ws.Route(ws.GET("").To(RouteBuilderTest{}.defaultTarget))
	ws.Route(ws.GET("/{serviceId}").To(RouteBuilderTest{}.defaultTarget))
	ws.Route(ws.POST("").To(RouteBuilderTest{}.defaultTarget))
	ws.Route(ws.PUT("/{serviceId}").To(RouteBuilderTest{}.defaultTarget))
	ws.Route(ws.DELETE("/{serviceId}").To(RouteBuilderTest{}.defaultTarget))
	container.Add(ws)
	return container
}
