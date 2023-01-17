// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestOutputsPopulator_PopulateOutputs(t *testing.T) {
	observer := new(MockResponseObserver)
	observer.On("Success", 200).Return()

	type fields struct {
		Endpoint  *Endpoint
		Outputs   *interface{}
		Error     error
		Observer  ResponseObserver
		Encoder   ResponseEncoder
		Describer RequestDescriber
	}
	tests := []struct {
		name         string
		fields       fields
		wantResponse *http.Response
		wantErr      bool
	}{
		{
			name: "Success",
			fields: fields{
				Endpoint: types.May[*Endpoint](NewEndpoint(http.MethodGet, "a", "b", "c").
					WithHandler(func() {}).
					WithOutputs(struct {
						Code      int                      `resp:"code" enum:"200,400"`
						SetCookie string                   `resp:"header"`
						Body      string                   `resp:"body"`
						Paging    paging.PaginatedResponse `resp:"paging"`
						ErrorBody webservice.ErrorV8       `resp:"body" error:"true"`
					}{}).
					Build()),
				Observer: observer,
				Outputs: types.OptionalOf[interface{}](struct {
					Code      int                      `resp:"code" enum:"200,400"`
					SetCookie string                   `resp:"header"`
					Body      map[string]string        `resp:"body"`
					Paging    paging.PaginatedResponse `resp:"paging"`
					ErrorBody webservice.ErrorV8       `resp:"body" error:"true"`
				}{
					Code:      200,
					SetCookie: `sessionId=38afes7a8; Max-Age=2592000`,
					Body:      map[string]string{"key": "value"},
					Paging: paging.PaginatedResponse{
						Size:             1,
						NumberOfElements: 1,
					},
				}).ValueInterfacePtr(),
			},
			wantResponse: &http.Response{
				Status:     http.StatusText(200),
				StatusCode: 200,
				Header: http.Header{
					HeaderSetCookie:   {"sessionId=38afes7a8; Max-Age=2592000"},
					HeaderContentType: {ContentTypeJson},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte(
					`{"content":{"key":"value"},"hasNext":false,"size":1,"numberOfElements":1,"number":0,"pageable":{"page":0,"size":0,"sort":{"sorted":false},"pagingState":null}}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Response{Header: make(http.Header)}
			e := EndpointResponseEncoder{
				Sink: NewHttpResponseDataSink(r),
			}
			p := &OutputsPopulator{
				Endpoint:  tt.fields.Endpoint,
				Outputs:   tt.fields.Outputs,
				Error:     tt.fields.Error,
				Observer:  tt.fields.Observer,
				Encoder:   e,
				Describer: tt.fields.Describer,
			}

			err := p.PopulateOutputs()
			gotResponse := r

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(tt.wantResponse, gotResponse),
					testhelpers.Diff(tt.wantResponse, gotResponse))
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestOutputsPopulator_EvaluateResponseCode(t *testing.T) {
	type fields struct {
		Endpoint  *Endpoint
		Outputs   *interface{}
		Error     error
		Observer  ResponseObserver
		Encoder   ResponseEncoder
		Describer RequestDescriber
	}
	tests := []struct {
		name     string
		fields   fields
		wantCode int
		wantErr  bool
	}{
		{
			name: "SuccessDefaultCode",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{
						Codes: EndpointResponseCodes{
							Success: []int{http.StatusCreated},
						},
						Success: EndpointResponseContent{
							Mime: MediaTypeJson,
						},
					}),
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "SuccessSystemDefaultCode",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{
						Success: EndpointResponseContent{
							Mime: MediaTypeJson,
						},
					}),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Error",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c"),
				Error:    errors.New("some error"),
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "ErrorCodeProvider",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c"),
				Error: webservice.NewStatusCodeError(
					errors.New("some error"),
					http.StatusForbidden),
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "SuccessNoBody",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{
						Codes: EndpointResponseCodes{
							Success: []int{http.StatusCreated},
						},
					}),
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "SuccessPort",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{
						Codes: EndpointResponseCodes{
							Success: []int{http.StatusCreated},
						},
						Port: func() *ops.Port {
							type portStruct struct {
								Code int `resp:"code"`
							}
							port, _ := PortReflector{}.ReflectOutputPort(reflect.TypeOf(portStruct{}))
							return port
						}(),
					}),
				Outputs: func() *interface{} {
					type portStruct struct {
						Code int `resp:"code"`
					}
					var result interface{} = portStruct{
						Code: http.StatusAccepted,
					}
					return &result
				}(),
			},
			wantCode: http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OutputsPopulator{
				Endpoint:  tt.fields.Endpoint,
				Outputs:   tt.fields.Outputs,
				Error:     tt.fields.Error,
				Observer:  tt.fields.Observer,
				Encoder:   tt.fields.Encoder,
				Describer: tt.fields.Describer,
			}
			gotCode, err := p.EvaluateResponseCode()
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantCode, gotCode)
			}
		})
	}
}

func TestOutputsPopulator_EvaluateSuccessBody(t *testing.T) {
	type fields struct {
		Endpoint  *Endpoint
		Outputs   *interface{}
		Error     error
		Observer  ResponseObserver
		Encoder   ResponseEncoder
		Describer RequestDescriber
	}
	tests := []struct {
		name     string
		fields   fields
		code     int
		wantBody interface{}
		wantErr  bool
	}{
		{
			name: "SuccessOutputsBody",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							Body types.Empty `resp:"body"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](struct{ Body types.Empty }{}).
					ValueInterfacePtr(),
				Describer: TestRequestDescriber{},
			},
			code:     http.StatusOK,
			wantBody: types.Empty{},
		},
		{
			name: "SuccessStaticPayload",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithSuccessPayload(types.Empty{})),
				Describer: TestRequestDescriber{},
			},
			code:     http.StatusOK,
			wantBody: types.Empty{},
		},
		{
			name: "SuccessPagingOutputsBody",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							Body   []int                    `resp:"body"`
							Paging paging.PaginatedResponse `resp:"paging"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](struct {
					Body   []int
					Paging paging.PaginatedResponse
				}{
					Body: []int{0, 1, 2},
					Paging: paging.PaginatedResponse{
						Size:             20,
						NumberOfElements: 3,
						Pageable: paging.PageableResponse{
							Size: 20,
						},
					},
				}).
					ValueInterfacePtr(),
				Describer: TestRequestDescriber{},
			},
			code: http.StatusOK,
			wantBody: paging.PaginatedResponse{
				Content:          []int{0, 1, 2},
				HasNext:          false,
				Size:             20,
				NumberOfElements: 3,
				Number:           0,
				Pageable: paging.PageableResponse{
					Size: 20,
				},
			},
		},
		{
			name: "SuccessOutputsEnvelopeBody",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithOperationId("getById").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							Body types.Empty `resp:"body" envelope:"true"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](struct{ Body types.Empty }{}).
					ValueInterfacePtr(),
				Describer: TestRequestDescriber{
					ReturnPath: "/a/b/c",
					ReturnParameters: map[string]interface{}{
						"tenantId": "c",
					},
				},
			},
			code: http.StatusOK,
			wantBody: integration.MsxEnvelope{
				Command:    "getById",
				HttpStatus: "OK",
				Message:    "getById succeeded",
				Params: map[string]interface{}{
					"tenantId": "c",
				},
				Payload: types.Empty{},
				Success: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OutputsPopulator{
				Endpoint:  tt.fields.Endpoint,
				Outputs:   tt.fields.Outputs,
				Error:     tt.fields.Error,
				Observer:  tt.fields.Observer,
				Encoder:   tt.fields.Encoder,
				Describer: tt.fields.Describer,
			}
			gotBody, err := p.EvaluateSuccessBody(tt.code)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantBody, gotBody)
			}
		})
	}
}

func TestOutputsPopulator_EvaluateErrorBody(t *testing.T) {
	type fields struct {
		Endpoint  *Endpoint
		Outputs   *interface{}
		Error     error
		Observer  ResponseObserver
		Encoder   ResponseEncoder
		Describer RequestDescriber
	}
	tests := []struct {
		name     string
		fields   fields
		code     int
		wantBody interface{}
		wantErr  bool
	}{
		{
			name: "ErrorOutputsBody",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							Body types.Empty `resp:"body"`
						}{})),
				Error: errors.New("some error"),
				Outputs: types.
					OptionalOf[interface{}](struct{ Body types.Empty }{}).
					ValueInterfacePtr(),
				Describer: TestRequestDescriber{},
			},
			code: http.StatusOK,
			wantBody: &webservice.ErrorV8{
				Code:    "UNKNOWN",
				Message: "some error",
			},
		},
		{
			name: "ErrorEnvelope",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithOperationId("getById").
					WithResponse(EndpointResponse{}.WithEnvelope(true)),
				Error: errors.New("some error"),
				Describer: TestRequestDescriber{
					ReturnParameters: map[string]interface{}{
						"factory": "worker",
					},
				},
			},
			code: http.StatusNotFound,
			wantBody: integration.MsxEnvelope{
				Command: "getById",
				Errors: []string{
					"some error",
				},
				HttpStatus: "NOT_FOUND",
				Message:    "some error",
				Params: map[string]interface{}{
					"factory": "worker",
				},
				Payload: nil,
				Throwable: &integration.Throwable{
					Message: "some error",
				},
			},
		},
		{
			name: "ErrorApplier",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.WithErrorPayload(webservice.ErrorV8{})),
				Error: webservice.NewCodedError("BXB", errors.New("some-error")),
			},
			code: http.StatusNotFound,
			wantBody: &webservice.ErrorV8{
				Code:    "BXB",
				Message: "some-error",
			},
		},
		{
			name: "ErrorRaw",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.WithErrorPayload(integration.ErrorDTO{})),
				Error: errors.New("some-error"),
				Describer: TestRequestDescriber{
					ReturnPath: "/a/b/c",
				},
			},
			code: http.StatusNotFound,
			wantBody: &integration.ErrorDTO{
				Code:       "404",
				Message:    "some-error",
				Path:       "/a/b/c",
				HttpStatus: "NOT_FOUND",
				Timestamp:  "",
			},
		},
		{
			name: "BadError",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.WithErrorPayload(struct{}{})),
				Error: errors.New("some-error"),
			},
			code:    http.StatusNotFound,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OutputsPopulator{
				Endpoint:  tt.fields.Endpoint,
				Outputs:   tt.fields.Outputs,
				Error:     tt.fields.Error,
				Observer:  tt.fields.Observer,
				Encoder:   tt.fields.Encoder,
				Describer: tt.fields.Describer,
			}
			gotBody, err := p.EvaluateErrorBody(tt.code)
			assert.Equal(t, tt.wantErr, err != nil, err)
			if !tt.wantErr {
				if wantDtoBody, ok := tt.wantBody.(*integration.ErrorDTO); ok {
					if gotDtoBody, ok := gotBody.(*integration.ErrorDTO); ok {
						wantDtoBody.Timestamp = gotDtoBody.Timestamp
					}
				}

				assert.True(t,
					reflect.DeepEqual(tt.wantBody, gotBody),
					testhelpers.Diff(tt.wantBody, gotBody))
			}
		})
	}
}

func TestOutputsPopulator_EvaluateHeaders(t *testing.T) {
	type fields struct {
		Endpoint  *Endpoint
		Outputs   *interface{}
		Error     error
		Observer  ResponseObserver
		Encoder   ResponseEncoder
		Describer RequestDescriber
	}
	tests := []struct {
		name       string
		fields     fields
		code       int
		wantHeader http.Header
		wantErr    bool
	}{
		{
			name: "SuccessHeaderPrimitive",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							CustomHeader string `resp:"header"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](
					struct{ CustomHeader string }{
						CustomHeader: "custom-header",
					}).
					ValueInterfacePtr(),
			},
			code: http.StatusOK,
			wantHeader: http.Header{
				"Custom-Header": {"custom-header"},
			},
		},
		{
			name: "SuccessHeaderArray",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							CustomHeader []string `resp:"header" style:"deepObject"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](
					struct{ CustomHeader []string }{
						CustomHeader: []string{"custom-header", "part-two"},
					}).
					ValueInterfacePtr(),
			},
			code: http.StatusOK,
			wantHeader: http.Header{
				"Custom-Header": {"custom-header,part-two"},
			},
		},
		{
			name: "SuccessHeaderObject",
			fields: fields{
				Endpoint: NewEndpoint(http.MethodGet, "a", "b", "c").
					WithResponse(EndpointResponse{}.
						WithOutputs(struct {
							CustomHeader map[string]string `resp:"header" explode:"true"`
						}{})),
				Outputs: types.
					OptionalOf[interface{}](
					struct{ CustomHeader map[string]string }{
						CustomHeader: map[string]string{
							"A": "custom-header",
							"B": "part-two",
						},
					}).
					ValueInterfacePtr(),
			},
			code: http.StatusOK,
			wantHeader: http.Header{
				"Custom-Header": {"A=custom-header,B=part-two"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Response{Header: make(http.Header)}
			e := EndpointResponseEncoder{
				Sink: NewHttpResponseDataSink(r),
			}
			p := &OutputsPopulator{
				Endpoint:  tt.fields.Endpoint,
				Outputs:   tt.fields.Outputs,
				Error:     tt.fields.Error,
				Observer:  tt.fields.Observer,
				Encoder:   e,
				Describer: tt.fields.Describer,
			}

			err := p.PopulateHeaders()
			gotHeader := r.Header

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(tt.wantHeader, gotHeader),
					testhelpers.Diff(tt.wantHeader, gotHeader))
			} else {
				assert.Error(t, err)
			}
		})
	}
}
