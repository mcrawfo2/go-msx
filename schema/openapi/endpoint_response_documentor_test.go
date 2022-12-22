// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/json"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"reflect"
	"testing"
)

type testError struct {
	Message string
}

func (t *testError) ApplyError(err error) {
	t.Message = err.Error()
}

func testEnvelope(code int, operationId string) *interface{} {
	var throwable *integration.Throwable
	var payload interface{}
	var errors []string
	success := code <= 399
	message := "Successfully executed "
	if !success {
		message = "Failed to execute "
		throwable = new(integration.Throwable)
		throwable.Message = "Service returned " + http.StatusText(code)
		errors = append(errors, throwable.Message)
	} else {
		payload = ""
	}

	return types.OptionalOf(integration.MsxEnvelope{
		Command:    operationId,
		Params:     map[string]interface{}{},
		HttpStatus: strcase.ToScreamingSnake(http.StatusText(code)),
		Message:    message + operationId,
		Payload:    payload,
		Success:    success,
		Throwable:  throwable,
		Errors:     errors,
	}).ValueInterfacePtr()
}

func testEnvelopeErrorResponse(code int, operationId string) *openapi3.Response {
	return new(openapi3.Response).
		WithDescription(http.StatusText(code)).
		WithContentItem("application/json",
			openapi3.MediaType{
				Schema:  NewSchemaRefPtr("types.Void.Envelope"),
				Example: testEnvelope(code, operationId),
			})
}

func TestEndpointResponseDocumentor(t *testing.T) {
	errorContent := map[string]openapi3.MediaType{
		"application/json": {
			Schema: NewSchemaRefPtr("webservice.ErrorV8"),
		},
	}

	tests := []struct {
		name     string
		doc      *EndpointResponseDocumentor
		response restops.EndpointResponse
		want     map[string]*openapi3.Response
		wantErr  bool
	}{
		{
			name:    "Skip",
			doc:     new(EndpointResponseDocumentor).WithSkip(true),
			want:    nil,
			wantErr: false,
		},
		{
			name: "Response",
			doc: new(EndpointResponseDocumentor).WithResponses(new(openapi3.Responses).
				WithMapOfResponseOrRefValuesItem("200",
					openapi3.ResponseOrRef{
						Response: new(openapi3.Response).
							WithDescription("OK"),
					})),
			response: restops.EndpointResponse{
				Success: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Payload: types.OptionalOf[interface{}](""),
				},
				Codes: restops.GetResponseCodes,
			},
			want: map[string]*openapi3.Response{
				"200": new(openapi3.Response).
					WithDescription("OK").
					WithContentItem(webservice.MIME_JSON,
						openapi3.MediaType{
							Schema: func() *openapi3.SchemaOrRef {
								s := NewSchemaOrRefPtr(StringSchema())
								s.Schema.ReflectType = reflect.TypeOf("")
								return s
							}(),
						}),
				"400": new(openapi3.Response).
					WithDescription("Bad Request").
					WithContent(errorContent),
				"401": new(openapi3.Response).
					WithDescription("Unauthorized").
					WithContent(errorContent),
				"403": new(openapi3.Response).
					WithDescription("Forbidden").
					WithContent(errorContent),
				"404": new(openapi3.Response).
					WithDescription("Not Found").
					WithContent(errorContent),
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			doc: new(EndpointResponseDocumentor).WithMutator(
				func(p *openapi3.Responses) {
					p.WithMapOfAnythingItem("x-msx-type", "integration.MsxEnvelope")
					p.WithMapOfResponseOrRefValuesItem("200",
						openapi3.ResponseOrRef{
							Response: new(openapi3.Response).
								WithDescription("Description").
								WithContentItem(webservice.MIME_JSON, openapi3.MediaType{
									Schema: NewSchemaOrRefPtr(StringSchema()),
								}),
						})
				}),
			want: map[string]*openapi3.Response{
				"200": new(openapi3.Response).
					WithDescription("Description").
					WithContentItem(webservice.MIME_JSON, openapi3.MediaType{
						Schema: NewSchemaOrRefPtr(StringSchema()),
					}),
			},
			wantErr: false,
		},
		{
			name: "Example",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{200},
				},
				Success: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Payload: types.OptionalOf[interface{}](0),
					Example: types.OptionalOf[interface{}](1),
				},
			},
			want: map[string]*openapi3.Response{
				"200": new(openapi3.Response).WithContentItem(
					webservice.MIME_JSON,
					openapi3.MediaType{
						Schema: func() *openapi3.SchemaOrRef {
							sr := NewSchemaOrRefPtr(IntegerSchema())
							sr.Schema.ReflectType = reflect.TypeOf(0)
							return sr
						}(),
						Example: types.OptionalOf(1).ValueInterfacePtr(),
					}).
					WithDescription("OK"),
			},
			wantErr: false,
		},
		{
			name: "Header",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{204},
				},
				Success: restops.EndpointResponseContent{
					Headers: map[string]restops.EndpointResponseHeader{
						"Retry-After": {
							Description: types.NewStringPtr("Seconds to wait before retrying"),
							Required:    types.NewBoolPtr(true),
							Payload:     types.OptionalOf[interface{}](29),
							Example:     types.OptionalOf[interface{}](30),
						},
					},
				},
			},
			want: map[string]*openapi3.Response{
				"204": new(openapi3.Response).
					WithDescription("No Content").
					WithHeaders(map[string]openapi3.HeaderOrRef{
						"Retry-After": {
							Header: func() *openapi3.Header {
								h := new(openapi3.Header).
									WithDescription("Seconds to wait before retrying").
									WithRequired(true).
									WithSchema(NewSchemaOrRef(IntegerSchema()))
								h.Schema.Schema.ReflectType = reflect.TypeOf(0)
								h.Example = types.OptionalOf(30).ValueInterfacePtr()
								return h
							}(),
						},
					}),
			},
		},
		{
			name: "HeaderAny",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{204},
				},
				Success: restops.EndpointResponseContent{
					Headers: map[string]restops.EndpointResponseHeader{
						"Retry-After": {
							Description: types.NewStringPtr("Any type"),
						},
					},
				},
			},
			want: map[string]*openapi3.Response{
				"204": new(openapi3.Response).
					WithDescription("No Content").
					WithHeaders(map[string]openapi3.HeaderOrRef{
						"Retry-After": openapi3.HeaderOrRef{
							Header: new(openapi3.Header).
								WithDescription("Any type").
								WithSchema(NewSchemaOrRef(AnySchema())),
						},
					}),
			},
		},
		{
			name: "HeaderReference",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{204},
				},
				Success: restops.EndpointResponseContent{
					Headers: map[string]restops.EndpointResponseHeader{
						"Retry-After": {
							Payload:   types.OptionalOf[interface{}](29),
							Reference: types.NewStringPtr("RetryAfter"),
						},
					},
				},
			},
			want: map[string]*openapi3.Response{
				"204": new(openapi3.Response).
					WithDescription("No Content").
					WithHeaders(map[string]openapi3.HeaderOrRef{
						"Retry-After": NewHeaderRef("RetryAfter"),
					}),
			},
		},
		{
			name: "Envelope",
			doc: new(EndpointResponseDocumentor).
				WithEndpoint(new(restops.Endpoint).
					WithOperationId("Envelope")),
			response: restops.EndpointResponse{
				Success: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Payload: types.OptionalOf[interface{}](""),
				},
				Envelope: true,
				Codes:    restops.GetResponseCodes,
			},
			want: map[string]*openapi3.Response{
				"200": new(openapi3.Response).
					WithDescription("OK").
					WithContentItem(webservice.MIME_JSON,
						openapi3.MediaType{
							Schema:  NewSchemaRefPtr("string.Envelope"),
							Example: testEnvelope(200, "Envelope"),
						}),
				"400": testEnvelopeErrorResponse(400, "Envelope"),
				"401": testEnvelopeErrorResponse(401, "Envelope"),
				"403": testEnvelopeErrorResponse(403, "Envelope"),
				"404": testEnvelopeErrorResponse(404, "Envelope"),
			},
			wantErr: false,
		},
		{
			name: "Paging",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{201},
					Error:   []int{400},
				},
				Success: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Paging:  types.OptionalOf[interface{}](paging.PaginatedResponseV8{}),
					Payload: types.OptionalOf[interface{}]([]int{}),
				},
			},
			want: map[string]*openapi3.Response{
				"201": new(openapi3.Response).
					WithDescription("Created").
					WithContentItem(webservice.MIME_JSON,
						openapi3.MediaType{
							Schema: NewSchemaRefPtr("int.List.Page"),
							Example: types.OptionalOf[interface{}](&paging.PaginatedResponseV8{
								Contents: []int{},
							}).ValueInterfacePtr(),
						}),
				"400": new(openapi3.Response).
					WithDescription("Bad Request").
					WithContent(errorContent),
			},
		},
		{
			name: "CustomError",
			response: restops.EndpointResponse{
				Codes: restops.EndpointResponseCodes{
					Success: []int{200},
					Error:   []int{400},
				},
				Success: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Payload: types.OptionalOf[interface{}]([]int{}),
				},
				Error: restops.EndpointResponseContent{
					Mime:    webservice.MIME_JSON,
					Payload: types.OptionalOf[interface{}](testError{}),
					Example: types.OptionalOf[interface{}](testError{}),
				},
			},
			want: map[string]*openapi3.Response{
				"200": new(openapi3.Response).
					WithDescription("OK").
					WithContentItem(webservice.MIME_JSON,
						openapi3.MediaType{
							Schema: NewSchemaOrRefPtr(ArraySchema(NewSchemaOrRef(IntegerSchema()))),
						}),
				"400": new(openapi3.Response).
					WithDescription("Bad Request").
					WithContent(map[string]openapi3.MediaType{
						webservice.MIME_JSON: {
							Schema:  NewSchemaRefPtr("openapi.testError"),
							Example: types.OptionalOf(testError{Message: "Bad Request"}).ValueInterfacePtr(),
						},
					}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.doc
			if e == nil {
				e = new(EndpointResponseDocumentor)
			}

			err := e.Document(&tt.response)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			got := e.Result()
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t,
				len(tt.want),
				len(got.MapOfResponseOrRefValues))
			for k, v := range got.MapOfResponseOrRefValues {
				wantBytes, _ := json.Marshal(tt.want[k])
				gotBytes, _ := json.Marshal(v.Response)

				assert.Equal(t,
					wantBytes, gotBytes,
					testhelpers.Diff(tt.want[k], v.Response))
			}
		})
	}
}

func TestEndpointResponseDocumentor_DocType(t *testing.T) {
	doc := new(EndpointResponseDocumentor)
	assert.Equal(t, DocType, doc.DocType())
}

func TestEndpointResponseDocumentor_Result(t *testing.T) {
	param := restops.EndpointResponse{
		Success: restops.EndpointResponseContent{
			Mime:    webservice.MIME_JSON,
			Payload: types.OptionalOf[interface{}](""),
		},
	}
	doc := new(EndpointResponseDocumentor)
	err := doc.Document(&param)
	assert.NoError(t, err)
	por := doc.Result()
	assert.NotNil(t, por)
}
