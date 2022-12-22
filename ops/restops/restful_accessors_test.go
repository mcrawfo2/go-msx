// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Route

func TestRouteWithEndpoint(t *testing.T) {
	ws := new(restful.WebService)
	endpoint := new(Endpoint)
	route := ws.
		GET("/api/v1/entity").
		To(func(request *restful.Request, response *restful.Response) {}).
		Do(RouteWithEndpoint(endpoint)).
		Build()

	e := EndpointFromRoute(route)
	assert.Equal(t, endpoint, e)
}

// Request

func TestRequestWithEndpoint(t *testing.T) {
	endpoint := new(Endpoint)
	httpRequest := new(http.Request)
	request := restful.NewRequest(httpRequest)

	result := RequestWithEndpoint(request, endpoint)
	e := EndpointFromRequest(result)
	assert.Equal(t, endpoint, e)
}

func TestRequestWithInputs(t *testing.T) {
	inputs := &struct{}{}
	httpRequest := new(http.Request)
	request := restful.NewRequest(httpRequest)

	result := RequestWithInputs(request, inputs)
	e := InputsFromRequest(result)
	assert.Equal(t, inputs, e)
}

func TestRequestWithOutputs(t *testing.T) {
	outputs := &struct{}{}
	httpRequest := new(http.Request)
	request := restful.NewRequest(httpRequest)

	result := RequestWithOutputs(request, outputs)
	e := OutputsFromRequest(result)
	assert.Equal(t, outputs, e)
}

func TestRequestWithEndpointRequestDecoder(t *testing.T) {
	decoder := new(EndpointRequestDecoder)
	httpRequest := new(http.Request)
	request := restful.NewRequest(httpRequest)

	result := RequestWithEndpointRequestDecoder(request, decoder)
	e := EndpointRequestDecoderFromRequest(result)
	assert.Equal(t, decoder, e)
}

func TestRequestWithEndpointResponseEncoder(t *testing.T) {
	decoder := new(EndpointResponseEncoder)
	httpRequest := new(http.Request)
	request := restful.NewRequest(httpRequest)

	result := RequestWithEndpointResponseEncoder(request, decoder)
	e := EndpointResponseEncoderFromRequest(result)
	assert.Equal(t, decoder, e)
}
