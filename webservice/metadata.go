// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

// Route Metadata Accessors

func RouteBuilderWithTagDefinition(b *restful.RouteBuilder, tag spec.TagProps) {
	b.Metadata(MetadataTagDefinition, tag)
}

func TagDefinitionFromRoute(r restful.Route) (tag spec.TagProps, ok bool) {
	val, ok := r.Metadata[MetadataTagDefinition]
	if !ok {
		return
	}

	tag, ok = val.(spec.TagProps)
	return
}

func RouteBuilderWithPermissions(b *restful.RouteBuilder, permissions []string) {
	b.Metadata(MetadataPermissions, permissions)
}

func PermissionsFromRoute(r restful.Route) (permissions []string, ok bool) {
	val, ok := r.Metadata[MetadataPermissions]
	if !ok {
		return
	}

	permissions, ok = val.([]string)
	return
}

func RouteBuilderWithSuccessResponse(b *restful.RouteBuilder, body interface{}) {
	b.Metadata(MetadataSuccessResponse, body)
}

func RouteBuilderWithEnvelopedPayload(b *restful.RouteBuilder, payload interface{}) {
	b.Metadata(MetadataEnvelope, payload)
}

func EnvelopedPayloadFromRoute(r restful.Route) (payload interface{}, ok bool) {
	payload, ok = r.Metadata[MetadataEnvelope]
	return
}

func RouteBuilderWithErrorPayload(b *restful.RouteBuilder, payload interface{}) {
	b.Metadata(MetadataErrorPayload, payload)
}

func ErrorPayloadFromRoute(r restful.Route) (payload interface{}, ok bool) {
	payload, ok = r.Metadata[MetadataErrorPayload]
	return
}

func RouteInputParamsInjector(params interface{}) restfulcontext.RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		RouteBuilderWithInputParams(builder, params)
	}
}

func RouteBuilderWithInputParams(b *restful.RouteBuilder, params interface{}) {
	b.Metadata(MetadataParams, params)
}

func InputParamsFromRoute(route restful.Route) interface{} {
	return route.Metadata[MetadataParams]
}

func RouteBuilderWithRequest(b *restful.RouteBuilder, template interface{}) {
	b.Metadata(MetadataRequest, template)
}

func RouteBuilderWithDefaultReturnCode(b *restful.RouteBuilder, code int) {
	b.Metadata(MetadataDefaultReturnCode, code)
}

func RouteBuilderWithApiIgnore(b *restful.RouteBuilder, apiIgnore bool) {
	b.Metadata(MetadataApiIgnore, apiIgnore)
}

func ApiIgnoreFromRoute(r restful.Route) (isApiIgnored, ok bool) {
	val, ok := r.Metadata[MetadataApiIgnore]
	if !ok {
		return
	}

	isApiIgnored, ok = val.(bool)
	return
}

// Request Attribute Accessors

func RequestWithError(request *restful.Request, err error) *restful.Request {
	request.SetAttribute(AttributeError, err)
	return request
}

func ErrorFromRequest(request *restful.Request) error {
	val, _ := request.Attribute(AttributeError).(error)
	return val
}

func RequestWithErrorPayload(request *restful.Request, payload interface{}) {
	request.SetAttribute(AttributeErrorPayload, payload)
}

func ErrorPayloadFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeErrorPayload)
}

// Params returns the populated input struct for the request
// Deprecated
func Params(req *restful.Request) interface{} {
	return ParamsFromRequest(req)
}

// ParamsFromRequest returns the populated input struct for the request
func ParamsFromRequest(req *restful.Request) interface{} {
	return req.Attribute(AttributeParams)
}

func RequestWithParams(request *restful.Request, params interface{}) {
	request.SetAttribute(AttributeParams, params)
}

func RequestWithSilenceLog(request *restful.Request) {
	request.SetAttribute(AttributeSilenceLog, true)
}

func SilenceLogFromRequest(request *restful.Request) bool {
	val := request.Attribute(AttributeSilenceLog)
	silenceLog, ok := val.(bool)
	return ok && silenceLog
}
