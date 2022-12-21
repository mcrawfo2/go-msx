// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"reflect"
	"strings"
)

type SwaggerParam struct {
	Name             string
	Description      string
	DataType         string
	DataFormat       string
	In               string
	Required         bool
	AllowableValues  map[string]string
	AllowMultiple    bool
	DefaultValue     string
	CollectionFormat string
}

func SwaggerParamFromParamData(data restful.ParameterData) SwaggerParam {
	var in string

	switch data.Kind {
	case restful.QueryParameterKind:
		in = FieldGroupHttpQuery
	case restful.PathParameterKind:
		in = FieldGroupHttpPath
	case restful.FormParameterKind:
		in = FieldGroupHttpCookie
	case restful.HeaderParameterKind:
		in = FieldGroupHttpHeader
	case restful.BodyParameterKind:
		// not a param in our case
		in = FieldGroupHttpBody
	}

	return SwaggerParam{
		Name:             data.Name,
		Description:      data.Description,
		DataType:         data.DataType,
		DataFormat:       data.DataFormat,
		In:               in,
		Required:         data.Required,
		AllowableValues:  data.AllowableValues,
		AllowMultiple:    data.AllowMultiple,
		DefaultValue:     data.DefaultValue,
		CollectionFormat: data.CollectionFormat,
	}
}

// NewEndpointFromRoute transforms a restful.Route into an Endpoint
func NewEndpointFromRoute(route restful.Route, basePath string) *Endpoint {
	routePath := route.Path
	if basePath != "/" && strings.HasPrefix(routePath, basePath) {
		strippedPath := strings.TrimPrefix(routePath, basePath)
		if strippedPath == "" {
			strippedPath = "/"
		} else if strippedPath[0] != '/' {
			strippedPath = route.Path
		}
		routePath = strippedPath
	}

	result := NewEndpoint(route.Method, routePath).
		WithOperationId(route.Operation).
		WithTags(webservice.TagsFromRoute(route)...).
		WithPermissionAnyOf(func() []string {
			perms, _ := webservice.PermissionsFromRoute(route)
			return perms
		}()...)

	if route.Doc != "" {
		result.WithSummary(route.Doc)
	}

	if route.Notes != "" {
		result.WithDescription(route.Notes)
	}

	ctx := webservice.ContextWithRoute(context.Background(), &route)

	// Request
	legacyInputPortStruct := webservice.InputParamsFromRoute(route)
	if legacyInputPortStruct != nil {
		inputPort, err := webservice.GetRouteParams(ctx, legacyInputPortStruct)
		if err != nil {
			logger.WithError(err).Error("Failed to parse RouteParams")
		}

		for _, field := range inputPort.Fields {
			applyRouteParamToEndpoint(result, field)
		}
	}

	// Response
	if payload, ok := webservice.EnvelopedPayloadFromRoute(route); ok {
		result.Response = result.Response.WithEnvelope(true)
		result.Response = result.Response.WithSuccessPayload(payload)
	} else if payload, ok = route.Metadata[webservice.MetadataSuccessResponse]; ok {
		result.Response = result.Response.WithEnvelope(false)
		result.Response = result.Response.WithSuccessPayload(payload)
	} else if payload = route.ReadSample; payload != nil {
		result.Response = result.Response.WithEnvelope(false)
		result.Response = result.Response.WithSuccessPayload(payload)
	}

	if result.Response.Envelope {
		// Envelope contains the error
	} else if payload, ok := webservice.ErrorPayloadFromRoute(route); ok {
		result.Response = result.Response.WithErrorPayload(payload)
	} else {
		result.Response = result.Response.WithErrorPayload(integration.ErrorDTO{})
	}

	for code := range route.ResponseErrors {
		if code < 399 {
			result.Response.Codes.Success = append(result.Response.Codes.Success, code)
		} else {
			result.Response.Codes.Error = append(result.Response.Codes.Error, code)
		}
	}

	return result
}

func applyRouteParamToEndpoint(e *Endpoint, r *webservice.RouteParam) {

	switch r.Parameter.Kind {
	case restful.PathParameterKind, restful.QueryParameterKind, restful.HeaderParameterKind:
		swaggerParam := SwaggerParamFromParamData(r.Parameter)
		p := EndpointRequestParameterFromSwaggerParam(swaggerParam, r.Options)
		e.Request = e.Request.WithParameter(p)
	case restful.FormParameterKind:
		swaggerParam := SwaggerParamFromParamData(r.Parameter)
		ff := EndpointRequestBodyFormFieldFromSwaggerParam(swaggerParam, r.Options)
		e.Request.Body = e.Request.Body.WithFormField(ff)
	case restful.BodyParameterKind:
		example := types.Instantiate(r.Field.Type)

		e.Request = e.Request.WithBody(EndpointRequestBody{
			Description: r.Parameter.Description,
			Required:    r.Field.Type.Kind() != reflect.Ptr || r.Parameter.Required,
			Mime:        MediaTypeJson,
			Example:     types.OptionalOf(example),
			Payload:     types.OptionalOf(example),
		})
	}
}
