// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func testNewRequestValidator(t *testing.T, portStruct interface{}, dataSource RequestDataSource) RequestValidator {
	port, err := PortReflector{}.ReflectInputPort(reflect.TypeOf(portStruct))
	assert.NoError(t, err)
	assert.NotNil(t, port)

	decoder := NewRequestDecoder(dataSource)
	assert.NotNil(t, decoder)

	return NewRequestValidator(port, decoder)
}

func TestNewRequestValidator(t *testing.T) {
	type testStruct struct {
		Primitive string `req:"header"`
	}

	validator := testNewRequestValidator(t, testStruct{}, MockRequestDataSource{})
	assert.NotNil(t, validator.port)
	assert.NotNil(t, validator.decoder)
}

func TestRequestValidator_ValidateRequest(t *testing.T) {
	tests := []struct {
		name       string
		portStruct interface{}
		dataSource RequestDataSource
		schema     map[string]*jsonschema.Schema
		wantErr    bool
	}{
		{
			name: "SingleField",
			portStruct: struct {
				A string `req:"header"`
			}{},
			dataSource: MockRequestDataSource{
				headers: map[string][]string{
					"A": {"123"},
				},
			},
			schema: map[string]*jsonschema.Schema{
				"A": js.StringSchema().WithPattern(`\d+`),
			},
		},
		{
			name: "MissingOptional",
			portStruct: struct {
				A string `req:"header" optional:"true"`
			}{},
			dataSource: MockRequestDataSource{},
			schema: map[string]*jsonschema.Schema{
				"A": js.StringSchema().WithPattern(`\d+`),
			},
		},
		{
			name: "ValidationFailure",
			portStruct: struct {
				A string `req:"header"`
			}{},
			dataSource: MockRequestDataSource{},
			schema: map[string]*jsonschema.Schema{
				"A": js.StringSchema().WithPattern(`\d+`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterPortFieldValidationSchemaFunc(func(field *ops.PortField) (schema js.ValidationSchema, err error) {
				jsonSchema, ok := tt.schema[field.Peer]
				if ok {
					return js.NewValidationSchemaFromJsonSchema(jsonSchema)
				} else {
					return js.ValidationSchema{}, nil
				}
			})

			v := testNewRequestValidator(t, tt.portStruct, tt.dataSource)
			gotErr := v.ValidateRequest()
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestRequestValidator_GetFieldValue(t *testing.T) {
	type testStruct struct {
		Primitive string                  `req:"header"`
		Array     []string                `req:"query" explode:"true"`
		Object    map[string]string       `req:"header" explode:"true"`
		File      *multipart.FileHeader   `req:"form"`
		FileArray []*multipart.FileHeader `req:"form"`
		Body      types.Pojo              `req:"body"`
		Anything  map[string]interface{}  `req:"header" explode:"true"`
	}

	tests := []struct {
		name       string
		fieldName  string
		dataSource RequestDataSource
		schema     *jsonschema.Schema
		wantValue  interface{}
		wantErr    bool
	}{
		{
			name:      "Primitive",
			fieldName: "Primitive",
			dataSource: MockRequestDataSource{
				headers: map[string][]string{
					"Primitive": {"123"},
				},
			},
			schema:    js.StringSchema().WithPattern(`\d+`),
			wantValue: "123",
		},
		{
			name:      "Array",
			fieldName: "Array",
			dataSource: MockRequestDataSource{
				query: url.Values{
					"array": {"abc", "123"},
				},
			},
			schema:    js.ArraySchema(*js.StringSchema()),
			wantValue: []interface{}{"abc", "123"},
		},
		{
			name:      "Object",
			fieldName: "Object",
			dataSource: MockRequestDataSource{
				headers: http.Header{
					"Object": {"R=100,G=200,B=150"},
				},
			},
			wantValue: types.Pojo{
				"R": 100,
				"G": 200,
				"B": 150,
			},
			schema: js.ObjectSchema().
				WithPropertiesItem("R", js.IntegerSchema().ToSchemaOrBool()).
				WithPropertiesItem("G", js.IntegerSchema().ToSchemaOrBool()).
				WithPropertiesItem("B", js.IntegerSchema().ToSchemaOrBool()),
		},
		{
			name:      "ObjectEmpty",
			fieldName: "Object",
			dataSource: MockRequestDataSource{
				headers: http.Header{},
			},
			wantValue: nil,
			schema: js.ObjectSchema().
				WithPropertiesItem("R", js.IntegerSchema().ToSchemaOrBool()).
				WithPropertiesItem("G", js.IntegerSchema().ToSchemaOrBool()).
				WithPropertiesItem("B", js.IntegerSchema().ToSchemaOrBool()),
		},
		{
			name:      "File",
			fieldName: "File",
			dataSource: MockRequestDataSource{
				multipartForm: func() *multipart.Form {
					r := bytes.NewBufferString(`--boundary
Content-Disposition: form-data; name="file"; filename="example.txt"

abc123
--boundary--`)
					form, err := multipart.NewReader(r, "boundary").ReadForm(10 * 1024 * 1024)
					assert.NoError(t, err)
					return form
				}(),
			},
			schema:    js.StringSchema().WithFormat("binary"),
			wantValue: "abc123",
		},
		{
			name:      "FileArray",
			fieldName: "FileArray",
			dataSource: MockRequestDataSource{
				multipartForm: func() *multipart.Form {
					r := bytes.NewBufferString(`--boundary
Content-Disposition: form-data; name="fileArray"; filename="example.txt"

abc123
--boundary
Content-Disposition: form-data; name="fileArray"; filename="example2.txt"

def456
--boundary--`)
					form, err := multipart.NewReader(r, "boundary").ReadForm(10 * 1024 * 1024)
					assert.NoError(t, err)
					return form
				}(),
			},
			schema:    js.ArraySchema(*js.StringSchema().WithFormat("binary")),
			wantValue: []string{"abc123", "def456"},
		},
		{
			name:      "Content",
			fieldName: "Body",
			dataSource: MockRequestDataSource{
				body: []byte(`{"a":"123","b":"456"}`),
			},
			schema: js.ObjectSchema(),
			wantValue: map[string]interface{}{
				"a": "123",
				"b": "456",
			},
		},
		{
			name:      "Any",
			fieldName: "Anything",
			dataSource: MockRequestDataSource{
				headers: http.Header{
					"Anything": {"R=100,G=200,B=150"},
				},
			},
			wantValue: types.Pojo{
				"R": "100",
				"G": "200",
				"B": "150",
			},
			schema: js.ObjectSchema().
				WithAdditionalProperties(js.AnySchema().ToSchemaOrBool()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := testNewRequestValidator(t, testStruct{}, tt.dataSource)

			portField := v.port.Fields.First(func(p *ops.PortField) bool {
				return p.Name == tt.fieldName
			})
			assert.NotNil(t, portField)

			validationSchema, err := js.NewValidationSchemaFromJsonSchema(tt.schema)

			assert.NoError(t, err)
			gotValue, gotErr := v.GetFieldValue(portField, validationSchema)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.True(t,
					reflect.DeepEqual(tt.wantValue, gotValue),
					testhelpers.Diff(tt.wantValue, gotValue))
			}
		})
	}
}

func TestRequestValidator_GetPrimitiveElement(t *testing.T) {
	type args struct {
		value string
		types []string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "String",
			args: args{
				value: "123",
				types: []string{"string"},
			},
			want: "123",
		},
		{
			name: "Integer",
			args: args{
				value: "123",
				types: []string{"integer"},
			},
			want: 123,
		},
		{
			name: "Number",
			args: args{
				value: "123",
				types: []string{"number"},
			},
			want: float64(123),
		},
		{
			name: "Boolean",
			args: args{
				value: "true",
				types: []string{"boolean"},
			},
			want: true,
		},
		{
			name: "ArrayFailure",
			args: args{
				value: "123",
				types: []string{"array"},
			},
			wantErr: true,
		},
		{
			name: "ObjectFailure",
			args: args{
				value: "123",
				types: []string{"object"},
			},
			wantErr: true,
		},
		{
			name: "CustomFailure",
			args: args{
				value: "123",
				types: []string{"null"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := RequestValidator{}
			got, err := v.GetPrimitiveElement(tt.args.value, tt.args.types)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(tt.want, got),
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}
