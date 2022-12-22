// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type EndpointsProducerFunc func() (Endpoints, error)

func (f EndpointsProducerFunc) Endpoints() (Endpoints, error) {
	return f()
}

func (f EndpointsProducerFunc) EndpointTransformers() EndpointTransformers {
	return nil
}

func TestRestfulEndpointRegisterer(t *testing.T) {
	const rootPath = "/api"
	endpointCount := 0

	ws := new(restful.WebService)

	server := NewMockRestfulWebServer(t)
	server.EXPECT().
		NewService(rootPath).
		Run(func(root string) {
			endpointCount++
			assert.Equal(t, rootPath, root)
		}).
		Return(ws, nil)

	registerer := NewRestfulEndpointRegisterer(server)

	err := registerer.RegisterEndpoints(EndpointsProducerFunc(func() (Endpoints, error) {
		return Endpoints{
			new(Endpoint),
			new(Endpoint),
		}, nil
	}))
	assert.NoError(t, err)

	assert.Equal(t, 1, endpointCount)
}

func TestContextEndpointRegisterer(t *testing.T) {
	ctx := context.Background()
	server := new(webservice.WebServer)
	ctx = webservice.ContextWithWebServerValue(ctx, server)
	want := &RestfulEndpointRegisterer{
		w: server,
	}
	got := ContextEndpointRegisterer(ctx)

	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}

func TestInjectRequestEndpointFilter(t *testing.T) {
	e := new(Endpoint)
	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectRequestEndpointFilter(e)).
		WithRequestPredicate(webservicetest.RequestHasAttribute(AttributeKeyEndpoint, e)).
		Test(t)
}

func TestInjectEndpointRequestDecoder(t *testing.T) {
	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectEndpointRequestDecoder).
		WithRequestPredicate(webservicetest.RequestHasAttributeAnyValue(AttributeKeyEndpointRequestDecoder)).
		Test(t)
}

func TestInjectEndpointResponseEncoder(t *testing.T) {
	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectEndpointResponseEncoder).
		WithRequestPredicate(webservicetest.RequestHasAttributeAnyValue(AttributeKeyEndpointResponseEncoder)).
		Test(t)
}

func TestEndpointController(t *testing.T) {
	called := false
	e, err := NewEndpoint("GET", "/api/v1/entity").
		WithOperationId("listEntities").
		WithRequestParameter(EndpointRequestParameter{
			Name: "entityIds",
			In:   FieldGroupHttpQuery,
			Type: types.NewStringPtr("array"),
		}).
		WithHandler(func(ctx context.Context, request *restful.Request) error {
			called = true
			decoder := EndpointRequestDecoderFromRequest(request)
			assert.NotNil(t, decoder)
			return nil
		}).
		Build()
	assert.NoError(t, err)

	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectRequestEndpointFilter(e)).
		WithRouteFilter(InjectEndpointRequestDecoder).
		WithRouteFilter(InjectEndpointResponseEncoder).
		WithRouteFilter(EndpointRequestFilter).
		WithRouteFilter(EndpointResponseFilter).
		WithRouteTarget(EndpointController(e)).
		WithRequestPredicate(webservicetest.RequestHasAttributeAnyValue(AttributeKeyEndpointRequestDecoder)).
		WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNoContent)).
		WithResponsePredicate(webservicetest.ResponseHasNoBody()).
		Test(t)

	assert.True(t, called)
}

func TestEndpointHandler_Restful(t *testing.T) {
	called := false
	e, err := NewEndpoint("GET", "/api/v1/entity").
		WithOperationId("listEntities").
		WithHandler(func(request *restful.Request, response *restful.Response) error {
			called = true
			return nil
		}).
		Build()
	assert.NoError(t, err)

	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectRequestEndpointFilter(e)).
		WithRouteFilter(InjectEndpointRequestDecoder).
		WithRouteFilter(InjectEndpointResponseEncoder).
		WithRouteFilter(EndpointRequestFilter).
		WithRouteFilter(EndpointResponseFilter).
		WithRouteTarget(EndpointController(e)).
		WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNoContent)).
		WithResponsePredicate(webservicetest.ResponseHasNoBody()).
		Test(t)

	assert.True(t, called)
}

func TestEndpointHandler_Http(t *testing.T) {
	called := false
	e, err := NewEndpoint("GET", "/api/v1/entity").
		WithOperationId("listEntities").
		WithHandler(func(request *http.Request, response http.ResponseWriter) error {
			called = true
			return nil
		}).
		Build()
	assert.NoError(t, err)

	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectRequestEndpointFilter(e)).
		WithRouteFilter(InjectEndpointRequestDecoder).
		WithRouteFilter(InjectEndpointResponseEncoder).
		WithRouteFilter(EndpointRequestFilter).
		WithRouteFilter(EndpointResponseFilter).
		WithRouteTarget(EndpointController(e)).
		WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNoContent)).
		WithResponsePredicate(webservicetest.ResponseHasNoBody()).
		Test(t)

	assert.True(t, called)
}

func TestEndpointHandler_Endpoint(t *testing.T) {
	called := false
	e, err := NewEndpoint("GET", "/api/v1/entity").
		WithOperationId("listEntities").
		WithHandler(func(decoder ops.InputDecoder, encoder ResponseEncoder) error {
			called = true
			return nil
		}).
		Build()
	assert.NoError(t, err)

	new(webservicetest.RouteBuilderTest).
		WithRouteFilter(InjectRequestEndpointFilter(e)).
		WithRouteFilter(InjectEndpointRequestDecoder).
		WithRouteFilter(InjectEndpointResponseEncoder).
		WithRouteFilter(EndpointRequestFilter).
		WithRouteFilter(EndpointResponseFilter).
		WithRouteTarget(EndpointController(e)).
		WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNoContent)).
		WithResponsePredicate(webservicetest.ResponseHasNoBody()).
		Test(t)

	assert.True(t, called)
}
