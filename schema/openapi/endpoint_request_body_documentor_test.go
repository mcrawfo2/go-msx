// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"reflect"
	"testing"
)

func TestEndpointRequestBodyDocumentor_Document(t *testing.T) {
	tests := []struct {
		name        string
		doc         *EndpointRequestBodyDocumentor
		requestBody restops.EndpointRequestBody
		want        *openapi3.RequestBody
		wantErr     bool
	}{
		{
			name:    "Skip",
			doc:     new(EndpointRequestBodyDocumentor).WithSkip(true),
			want:    nil,
			wantErr: false,
		},
		{
			name: "RequestBody",
			doc: new(EndpointRequestBodyDocumentor).WithRequestBody(
				new(openapi3.RequestBody).
					WithDescription("Description")),
			requestBody: restops.EndpointRequestBody{
				Mime:     webservice.MIME_JSON,
				Payload:  types.OptionalOf[interface{}]([]types.UUID{}),
				Required: true,
			},
			want: &openapi3.RequestBody{
				Description: types.NewStringPtr("Description"),
				Content: map[string]openapi3.MediaType{
					webservice.MIME_JSON: {
						Schema: func() *openapi3.SchemaOrRef {
							sr := NewSchemaOrRefPtr(ArraySchema(NewSchemaRef("UUID")))
							sr.Schema.ReflectType = reflect.TypeOf([]types.UUID{})
							return sr
						}(),
					},
				},
				Required: types.NewBoolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			doc: new(EndpointRequestBodyDocumentor).WithMutator(
				func(p *openapi3.RequestBody) {
					p.WithDescription("Description")
				}),
			requestBody: restops.EndpointRequestBody{
				Mime:    webservice.MIME_JSON,
				Payload: types.OptionalOf[interface{}](""),
			},
			want: &openapi3.RequestBody{
				Description: types.NewStringPtr("Description"),
				Content: map[string]openapi3.MediaType{
					webservice.MIME_JSON: {
						Schema: func() *openapi3.SchemaOrRef {
							sr := NewSchemaOrRefPtr(StringSchema())
							sr.Schema.ReflectType = reflect.TypeOf("")
							return sr
						}(),
					},
				},
				Required: types.NewBoolPtr(false),
			},
			wantErr: false,
		},
		{
			name:        "NoPayload",
			requestBody: restops.EndpointRequestBody{},
			wantErr:     false,
		},
		{
			name: "Example",
			requestBody: restops.EndpointRequestBody{
				Mime:    webservice.MIME_JSON,
				Payload: types.OptionalOf[interface{}](""),
				Example: types.OptionalOf[interface{}]("abc"),
			},
			want: &openapi3.RequestBody{
				Content: map[string]openapi3.MediaType{
					webservice.MIME_JSON: {
						Schema: func() *openapi3.SchemaOrRef {
							sr := NewSchemaOrRefPtr(StringSchema())
							sr.Schema.ReflectType = reflect.TypeOf("")
							return sr
						}(),
						Example: func() *interface{} {
							var v interface{} = "abc"
							return &v
						}(),
					},
				},
				Required: types.NewBoolPtr(false),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.doc
			if e == nil {
				e = new(EndpointRequestBodyDocumentor)
			}

			err := e.Document(&tt.requestBody)
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

			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestEndpointRequestBodyDocumentor_DocType(t *testing.T) {
	doc := new(EndpointRequestBodyDocumentor)
	assert.Equal(t, DocType, doc.DocType())
}

func TestEndpointRequestBodyDocumentor_Result(t *testing.T) {
	param := restops.EndpointRequestBody{
		Description: "description",
		Required:    true,
		Payload:     types.OptionalOf[interface{}]([]types.UUID{}),
	}
	doc := new(EndpointRequestBodyDocumentor)
	err := doc.Document(&param)
	assert.NoError(t, err)
	por := doc.Result()
	assert.NotNil(t, por)
}
