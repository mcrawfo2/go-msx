// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"reflect"
	"testing"
)

func TestEndpointParameterDocumentor_Document(t *testing.T) {
	tests := []struct {
		name      string
		doc       *EndpointParameterDocumentor
		parameter restops.EndpointRequestParameter
		want      *openapi3.ParameterOrRef
		wantErr   bool
	}{
		{
			name:    "Skip",
			doc:     new(EndpointParameterDocumentor).WithSkip(true),
			want:    nil,
			wantErr: false,
		},
		{
			name: "Parameter",
			doc: new(EndpointParameterDocumentor).WithParameter(
				new(openapi3.Parameter).
					WithDescription("Description")),
			parameter: restops.
				EndpointRequestParameter{Name: "Parameter"},
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name:        "Parameter",
					Description: types.NewStringPtr("Description"),
				},
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			doc: new(EndpointParameterDocumentor).WithMutator(
				func(p *openapi3.Parameter) {
					p.WithDescription("Description")
				}),
			parameter: restops.
				EndpointRequestParameter{Name: "Mutator"},
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name:        "Mutator",
					Description: types.NewStringPtr("Description"),
				},
			},
			wantErr: false,
		},
		{
			name: "Name",
			parameter: restops.
				EndpointRequestParameter{Name: "Name"},
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "Name",
				},
			},
			wantErr: false,
		},
		{
			name: "AllowEmptyValue",
			parameter: restops.
				EndpointRequestParameter{Name: "AllowEmptyValue"}.
				WithAllowEmptyValue(true),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name:            "AllowEmptyValue",
					AllowEmptyValue: types.NewBoolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "AllowReserved",
			parameter: restops.
				EndpointRequestParameter{Name: "AllowReserved"}.
				WithAllowReserved(true),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name:          "AllowReserved",
					AllowReserved: types.NewBoolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "Deprecated",
			parameter: restops.
				EndpointRequestParameter{Name: "Deprecated"}.
				WithDeprecated(true),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name:       "Deprecated",
					Deprecated: types.NewBoolPtr(true),
				},
			},
			wantErr: false,
		},
		{
			name: "Example",
			parameter: restops.
				EndpointRequestParameter{Name: "Example"}.
				WithExample("abc"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "Example",
					Example: func() *interface{} {
						var e interface{} = "abc"
						return &e
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "Reference",
			parameter: restops.
				EndpointRequestParameter{Name: "Reference"}.
				WithReference("reference"),
			want: &openapi3.ParameterOrRef{
				ParameterReference: &openapi3.ParameterReference{
					Ref: "#/components/parameters/reference",
				},
			},
			wantErr: false,
		},
		{
			name: "Payload",
			parameter: restops.EndpointRequestParameter{Name: "Payload"}.
				WithPayload(""),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "Payload",
					Schema: &openapi3.SchemaOrRef{
						Schema: func() *openapi3.Schema {
							s := StringSchema()
							s.ReflectType = reflect.TypeOf("")
							return s
						}(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeString",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeString"}.
				WithType("string").
				WithFormat("date-time"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeString",
					Schema: &openapi3.SchemaOrRef{
						Schema: StringSchema().WithFormat("date-time"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeNumber",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeNumber"}.
				WithType("number"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeNumber",
					Schema: &openapi3.SchemaOrRef{
						Schema: NumberSchema(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeInteger",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeInteger"}.
				WithType("integer"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeInteger",
					Schema: &openapi3.SchemaOrRef{
						Schema: IntegerSchema(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeBoolean",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeBoolean"}.
				WithType("boolean"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeBoolean",
					Schema: &openapi3.SchemaOrRef{
						Schema: BooleanSchema(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeObject",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeObject"}.
				WithType("object"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeObject",
					Schema: &openapi3.SchemaOrRef{
						Schema: ObjectSchema(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TypeArray",
			parameter: restops.
				EndpointRequestParameter{Name: "TypeArray"}.
				WithType("array"),
			want: &openapi3.ParameterOrRef{
				Parameter: &openapi3.Parameter{
					Name: "TypeArray",
					Schema: &openapi3.SchemaOrRef{
						Schema: ArraySchema(NewSchemaOrRef(new(openapi3.Schema))),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.doc
			if e == nil {
				e = new(EndpointParameterDocumentor)
			}

			err := e.Document(&tt.parameter)
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

func TestEndpointParameterDocumentor_DocType(t *testing.T) {
	doc := new(EndpointParameterDocumentor)
	assert.Equal(t, DocType, doc.DocType())
}

func TestEndpointParameterDocumentor_Result(t *testing.T) {
	param := restops.NewEndpointRequestParameter("Name", "query")
	doc := new(EndpointParameterDocumentor)
	err := doc.Document(&param)
	assert.NoError(t, err)
	por := doc.Result()
	assert.NotNil(t, por)
}
