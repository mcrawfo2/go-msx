// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestEndpointResponse_WithDocumentor(t *testing.T) {
	e := EndpointResponse{}
	f := e.WithDocumentor(TestDocumentor[EndpointResponse]{})
	assert.NotNil(t, f)
}

func TestEndpointResponse_WithEnvelope(t *testing.T) {
	e := EndpointResponse{}
	f := e.WithEnvelope(true)
	assert.True(t, f.Envelope)
}

func TestEndpointResponse_WithResponseCodes(t *testing.T) {
	e := EndpointResponse{}
	r := CreateResponseCodes
	f := e.WithResponseCodes(r)
	assert.Equal(t, r, f.Codes)
}

func TestEndpointResponse_WithSuccessPayload(t *testing.T) {
	e := EndpointResponse{}
	r := struct{}{}
	f := e.WithSuccessPayload(r)
	assert.Equal(t, types.OptionalOf[interface{}](r), f.Success.Payload)
}

func TestEndpointResponse_WithSuccessHeader(t *testing.T) {
	e := EndpointResponse{}
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("Success Header"),
	}
	f := e.WithSuccessHeader("success-header", h)
	assert.Equal(t, h, f.Success.Headers["success-header"])
}

func TestEndpointResponse_WithErrorPayload(t *testing.T) {
	e := EndpointResponse{}
	r := struct{}{}
	f := e.WithErrorPayload(r)
	assert.Equal(t, types.OptionalOf[interface{}](r), f.Error.Payload)
}

func TestEndpointResponse_WithErrorHeader(t *testing.T) {
	e := EndpointResponse{}
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("Error Header"),
	}
	f := e.WithErrorHeader("error-header", h)
	assert.Equal(t, h, f.Error.Headers["error-header"])
}

func TestEndpointResponse_WithMime(t *testing.T) {
	e := EndpointResponse{}
	r := MediaTypeFormUrlencoded
	f := e.WithMime(r)
	assert.Equal(t, r, f.Success.Mime)
	assert.Equal(t, r, f.Error.Mime)
}

func TestEndpointResponse_WithSuccess(t *testing.T) {
	e := EndpointResponse{}
	r := EndpointResponseContent{
		Mime: MediaTypeJson,
	}
	f := e.WithSuccess(r)
	assert.Equal(t, r, f.Success)
}

func TestEndpointResponse_WithError(t *testing.T) {
	e := EndpointResponse{}
	r := EndpointResponseContent{
		Mime: MediaTypeJson,
	}
	f := e.WithError(r)
	assert.Equal(t, r, f.Error)
}

func TestEndpointResponse_WithHeader(t *testing.T) {
	e := EndpointResponse{}
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("Error Header"),
	}
	f := e.WithHeader("header", h)
	assert.Equal(t, h, f.Success.Headers["header"])
	assert.Equal(t, h, f.Error.Headers["header"])
}

func TestEndpointResponse_HasBody(t *testing.T) {
	e := EndpointResponse{}
	r := EndpointResponseContent{
		Mime: MediaTypeJson,
	}
	f := e.WithSuccess(r)
	assert.True(t, f.HasBody())
}

func TestEndpointResponse_Produces(t *testing.T) {
	e := EndpointResponse{}
	r := EndpointResponseContent{
		Mime: MediaTypeJson,
	}
	f := e.WithSuccess(r)
	m := []string{MediaTypeJson}
	assert.Equal(t, m, f.Produces())
}

func TestEndpointResponse_withEnumCodes(t *testing.T) {
	e := EndpointResponse{}
	codes := "200,400,100,B"
	f := e.withEnumCodes(codes)
	assert.Equal(t, []int{200}, f.Codes.Success)
	assert.Equal(t, []int{400}, f.Codes.Error)
}

func TestEndpointResponse_WithOutputs(t *testing.T) {
	e := EndpointResponse{}
	f := e.WithOutputs(struct{}{})
	assert.NotNil(t, f.Port)
}

func TestEndpointResponse_WithOutputs_Invalid(t *testing.T) {
	e := EndpointResponse{}
	f := e.WithOutputs("")
	assert.Nil(t, f.Port)
}

func TestEndpointResponse_withPortField(t *testing.T) {
	tests := []struct {
		name string
		pf   *ops.PortField
		want EndpointResponse
	}{
		{
			name: "SuccessBody",
			pf: &ops.PortField{
				Name:     "body",
				Indices:  []int{},
				Peer:     "body",
				Group:    FieldGroupHttpBody,
				Optional: false,
				PortType: PortTypeRequest,
				Type: ops.PortFieldType{
					Type:        reflect.TypeOf(struct{}{}),
					HandlerType: reflect.TypeOf(struct{}{}),
				},
				Options: map[string]string{
					"envelope": "true",
					"success":  "true",
					"mime":     "text/plain",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: EndpointResponse{
				Envelope: true,
				Success: EndpointResponseContent{
					Mime:    "text/plain",
					Payload: types.OptionalOf[interface{}](struct{}{}),
				},
			},
		},
		{
			name: "ErrorBody",
			pf: &ops.PortField{
				Name:     "body",
				Indices:  []int{},
				Peer:     "body",
				Group:    FieldGroupHttpBody,
				Optional: false,
				PortType: PortTypeRequest,
				Type: ops.PortFieldType{
					Type:        reflect.TypeOf(struct{}{}),
					HandlerType: reflect.TypeOf(struct{}{}),
				},
				Options: map[string]string{
					"envelope": "true",
					"error":    "true",
					"mime":     "text/plain",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: EndpointResponse{
				Envelope: true,
				Error: EndpointResponseContent{
					Mime:    "text/plain",
					Payload: types.OptionalOf[interface{}](struct{}{}),
				},
			},
		},
		{
			name: "Code",
			pf: &ops.PortField{
				Name:     "Code",
				Indices:  []int{},
				Peer:     "code",
				Group:    FieldGroupHttpCode,
				Optional: false,
				PortType: PortTypeResponse,
				Type: ops.PortFieldType{
					Type:        reflect.TypeOf(0),
					HandlerType: reflect.TypeOf(0),
				},
				Options: map[string]string{
					"enum": "200,404",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: EndpointResponse{
				Codes: EndpointResponseCodes{
					Success: []int{200},
					Error:   []int{404},
				},
			},
		},
		{
			name: "ErrorHeader",
			pf: &ops.PortField{
				Name:     "Header",
				Indices:  []int{},
				Peer:     "some-header",
				Group:    FieldGroupHttpHeader,
				Optional: false,
				PortType: PortTypeResponse,
				Type: ops.PortFieldType{
					Type:        reflect.TypeOf(""),
					HandlerType: reflect.TypeOf(""),
				},
				Options: map[string]string{
					"error":   "true",
					"success": "false",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: EndpointResponse{
				Error: EndpointResponseContent{
					Headers: map[string]EndpointResponseHeader{
						"some-header": {
							Explode:  types.NewBoolPtr(false),
							Required: types.NewBoolPtr(true),
							PortField: &ops.PortField{
								Name:     "Header",
								Indices:  []int{},
								Peer:     "some-header",
								Group:    FieldGroupHttpHeader,
								Optional: false,
								PortType: PortTypeResponse,
								Type: ops.PortFieldType{
									Type:        reflect.TypeOf(""),
									HandlerType: reflect.TypeOf(""),
								},
								Options: map[string]string{
									"error":   "true",
									"success": "false",
								},
								Baggage: map[interface{}]interface{}{},
							},
						},
					},
				},
			},
		},
		{
			name: "SuccessPaging",
			pf: &ops.PortField{
				Name:     "pageing",
				Indices:  []int{},
				Peer:     "pageResponse",
				Group:    FieldGroupHttpPaging,
				Optional: false,
				PortType: PortTypeResponse,
				Type: ops.PortFieldType{
					Type:        reflect.TypeOf(struct{}{}),
					HandlerType: reflect.TypeOf(struct{}{}),
				},
				Options: map[string]string{
					"envelope": "true",
					"success":  "true",
					"mime":     "text/plain",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: EndpointResponse{
				Success: EndpointResponseContent{
					Mime:   "text/plain",
					Paging: types.OptionalOf[interface{}](struct{}{}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := EndpointResponse{}
			s := r.withPortField(tt.pf)
			assert.True(t,
				reflect.DeepEqual(tt.want, s),
				testhelpers.Diff(tt.want, s))
		})
	}
}

func TestNewEndpointResponse(t *testing.T) {
	e := NewEndpointResponse()
	f := EndpointResponse{}
	assert.Equal(t, f, e)
}

func TestEndpointResponseCodes_DefaultCode(t *testing.T) {
	type fields struct {
		Success []int
		Error   []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "NoCodes",
			fields: fields{},
			want:   http.StatusOK,
		},
		{
			name: "SomeCodes",
			fields: fields{
				Success: []int{http.StatusCreated},
			},
			want: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EndpointResponseCodes{
				Success: tt.fields.Success,
				Error:   tt.fields.Error,
			}
			assert.Equalf(t, tt.want, c.DefaultCode(), "DefaultCode()")
		})
	}
}

func TestDefaultResponseCodes(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   EndpointResponseCodes
	}{
		{
			name:   http.MethodGet,
			method: http.MethodGet,
			want:   GetResponseCodes,
		},
		{
			name:   http.MethodPost,
			method: http.MethodPost,
			want:   CreateResponseCodes,
		},
		{
			name:   http.MethodPut,
			method: http.MethodPut,
			want:   UpdateResponseCodes,
		},
		{
			name:   http.MethodDelete,
			method: http.MethodDelete,
			want:   DeleteResponseCodes,
		},
		{
			name:   http.MethodPatch,
			method: http.MethodPatch,
			want:   UnknownVerbResponseCodes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultResponseCodes(tt.method)
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestEndpointResponseContent_WithHeader(t *testing.T) {
	e := EndpointResponseContent{}
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("Error Header"),
	}
	f := e.WithHeader("header", h)
	assert.Equal(t, h, f.Headers["header"])
}

func TestEndpointResponseContent_WithExample(t *testing.T) {
	e := EndpointResponseContent{}
	x := struct{}{}
	f := e.WithExample(x)
	assert.Equal(t, types.OptionalOf[interface{}](x), f.Example)
}

func TestEndpointResponseContent_WithPayload(t *testing.T) {
	e := EndpointResponseContent{}
	x := struct{}{}
	f := e.WithPayload(x)
	assert.Equal(t, types.OptionalOf[interface{}](x), f.Payload)
	assert.Equal(t, MediaTypeJson, f.Mime)
}

func TestEndpointResponseHeader_WithDescription(t *testing.T) {
	e := EndpointResponseHeader{}
	d := "description"
	f := e.WithDescription(d)
	assert.Equal(t, &d, f.Description)
}

func TestEndpointResponseHeader_WithRequired(t *testing.T) {
	e := EndpointResponseHeader{}
	d := true
	f := e.WithRequired(d)
	assert.Equal(t, &d, f.Required)
}

func TestEndpointResponseHeader_WithDeprecated(t *testing.T) {
	e := EndpointResponseHeader{}
	d := true
	f := e.WithDeprecated(d)
	assert.Equal(t, &d, f.Deprecated)
}

func TestEndpointResponseHeader_WithAllowEmptyValue(t *testing.T) {
	e := EndpointResponseHeader{}
	d := true
	f := e.WithAllowEmptyValue(d)
	assert.Equal(t, &d, f.AllowEmptyValue)
}

func TestEndpointResponseHeader_WithExplode(t *testing.T) {
	e := EndpointResponseHeader{}
	d := true
	f := e.WithExplode(d)
	assert.Equal(t, &d, f.Explode)
}

func TestEndpointResponseHeader_WithAllowReserved(t *testing.T) {
	e := EndpointResponseHeader{}
	d := true
	f := e.WithAllowReserved(d)
	assert.Equal(t, &d, f.AllowReserved)
}

func TestEndpointResponseHeader_WithExample(t *testing.T) {
	e := EndpointResponseHeader{}
	x := struct{}{}
	f := e.WithExample(x)
	assert.Equal(t, types.OptionalOf[interface{}](x), f.Example)
}

func TestEndpointResponseHeader_WithPayload(t *testing.T) {
	e := EndpointResponseHeader{}
	x := struct{}{}
	f := e.WithPayload(x)
	assert.Equal(t, types.OptionalOf[interface{}](x), f.Payload)
}

func TestEndpointResponseHeader_WithReference(t *testing.T) {
	e := EndpointResponseHeader{}
	x := "SomeHeader"
	f := e.WithReference(x)
	assert.Equal(t, &x, f.Reference)
}

func TestNewEndpointResponseHeader(t *testing.T) {
	e := NewEndpointResponseHeader()
	f := EndpointResponseHeader{}
	assert.Equal(t, f, e)
}
