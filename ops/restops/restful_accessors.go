// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
)

const (
	MetadataKeyEndpoint = "Endpoint"

	AttributeKeyEndpoint                = "Endpoint"
	AttributeKeyInputs                  = "Inputs"
	AttributeKeyOutputs                 = "Outputs"
	AttributeKeyEndpointRequestDecoder  = "EndpointRequestDecoder"
	AttributeKeyEndpointResponseEncoder = "EndpointResponseEncoder"
)

// Route

func RouteWithEndpoint(e *Endpoint) restfulcontext.RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		builder.Metadata(MetadataKeyEndpoint, e)
	}
}

func EndpointFromRoute(route restful.Route) *Endpoint {
	if val, ok := route.Metadata[MetadataKeyEndpoint]; ok {
		return val.(*Endpoint)
	}
	return new(Endpoint)
}

// Request

func RequestWithEndpoint(request *restful.Request, e *Endpoint) *restful.Request {
	request.SetAttribute(AttributeKeyEndpoint, e)
	return request
}

func EndpointFromRequest(request *restful.Request) *Endpoint {
	return request.Attribute(AttributeKeyEndpoint).(*Endpoint)
}

func RequestWithInputs(request *restful.Request, inputs interface{}) *restful.Request {
	request.SetAttribute(AttributeKeyInputs, inputs)
	return request
}

func InputsFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeKeyInputs)
}

func RequestWithOutputs(request *restful.Request, outputs interface{}) *restful.Request {
	request.SetAttribute(AttributeKeyOutputs, outputs)
	return request
}

func OutputsFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeKeyOutputs)
}

func RequestWithEndpointRequestDecoder(request *restful.Request, decoder ops.InputDecoder) *restful.Request {
	request.SetAttribute(AttributeKeyEndpointRequestDecoder, decoder)
	return request
}

func EndpointRequestDecoderFromRequest(request *restful.Request) ops.InputDecoder {
	return request.Attribute(AttributeKeyEndpointRequestDecoder).(ops.InputDecoder)
}

func RequestWithEndpointResponseEncoder(request *restful.Request, encoder ResponseEncoder) *restful.Request {
	request.SetAttribute(AttributeKeyEndpointResponseEncoder, encoder)
	return request
}

func EndpointResponseEncoderFromRequest(request *restful.Request) ResponseEncoder {
	return request.Attribute(AttributeKeyEndpointResponseEncoder).(ResponseEncoder)
}
