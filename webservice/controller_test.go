// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestContextController(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(ContextController(
					func(ctx context.Context) (interface{}, error) {
						return nil, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("success", true)),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(ContextController(
					func(ctx context.Context) (interface{}, error) {
						return struct{}{}, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("success", true)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJson("responseObject", func(result gjson.Result) bool {
					return result.Raw == "{}"
				})),
		},
		{
			name: "Error",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(ContextController(
					func(ctx context.Context) (interface{}, error) {
						return nil, errors.New("some error")
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(400)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValueType("responseObject", gjson.Null)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestController(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(Controller(
					func(req *restful.Request) (interface{}, error) {
						return nil, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("success", true)),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(Controller(
					func(req *restful.Request) (interface{}, error) {
						return struct{}{}, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("success", true)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJson("responseObject", func(result gjson.Result) bool {
					return result.Raw == "{}"
				})),
		},
		{
			name: "Error",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(Controller(
					func(req *restful.Request) (interface{}, error) {
						return nil, errors.New("some error")
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(400)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValueType("responseObject", gjson.Null)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestEnsureSlash(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Slash",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestPath("/swagger/").
				WithRouteTarget(EnsureSlash(func(writer http.ResponseWriter, request *http.Request) {
					writer.WriteHeader(200)
				})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "NoSlash",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestPath("/swagger").
				WithRouteTarget(EnsureSlash(func(writer http.ResponseWriter, request *http.Request) {
					writer.WriteHeader(200)
				})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(302)).
				WithResponsePredicate(webservicetest.ResponseHasHeader("Location", "/swagger/")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestHttpHandlerController(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(HttpHandlerController(
					func(resp http.ResponseWriter, req *http.Request) {
						resp.WriteHeader(204)
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(HttpHandlerController(
					func(resp http.ResponseWriter, req *http.Request) {
						resp.Header().Set("Content-Type", MIME_JSON)
						resp.WriteHeader(200)
						_, _ = resp.Write([]byte(`{"example":"value"}`))
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasHeader(headerNameContentType, MIME_JSON)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("example", "value")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRawContextController(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawContextController(
					func(ctx context.Context) (interface{}, error) {
						return nil, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawContextController(
					func(ctx context.Context) (interface{}, error) {
						return map[string]string{
							"key": "value",
						}, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("key", "value")),
		},
		{
			name: "Error",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawContextController(
					func(ctx context.Context) (interface{}, error) {
						return nil, errors.New("some error")
					})).
				WithRouteBuilderDo(ErrorPayload(new(integration.ErrorDTO))).
				WithResponsePredicate(webservicetest.ResponseHasStatus(400)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("code", "400")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRawController(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawController(
					func(req *restful.Request) (interface{}, error) {
						return nil, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawController(
					func(req *restful.Request) (interface{}, error) {
						return map[string]string{
							"key": "value",
						}, nil
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("key", "value")),
		},
		{
			name: "Error",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(RawController(
					func(req *restful.Request) (interface{}, error) {
						return nil, errors.New("some error")
					})).
				WithRouteBuilderDo(ErrorPayload(new(integration.ErrorDTO))).
				WithResponsePredicate(webservicetest.ResponseHasStatus(400)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue("code", "400")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestStaticFileAlias(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Rewrite",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestPath("/swagger").
				WithRouteTarget(StaticFileAlias(
					StaticAlias{
						Path: "/swagger",
						File: "/swagger-ui.html",
					},
					func(resp http.ResponseWriter, req *http.Request) {
						assert.Equal(t, "/swagger-ui.html", req.URL.Path)
						resp.WriteHeader(204)
					})).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}
