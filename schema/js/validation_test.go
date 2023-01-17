// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js_test

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"sort"
	"testing"
)

func TestNewValidationSchema(t *testing.T) {
	schema := &jsv.Schema{
		Ref: &jsv.Schema{
			Types: []string{
				string(jsonschema.Null),
				string(jsonschema.Array),
				string(jsonschema.Object),
				string(jsonschema.String),
				string(jsonschema.Integer),
				string(jsonschema.Number),
				string(jsonschema.Boolean),
			},
		},
	}

	validationSchema := js.NewValidationSchema(schema)
	assert.NotNil(t, validationSchema)
}

func TestValidationSchema_Validate(t *testing.T) {
	tests := []struct {
		name    string
		schema  *jsv.Schema
		value   interface{}
		wantErr bool
	}{
		{
			name:    "Schema",
			schema:  &jsv.Schema{Types: []string{"null"}},
			value:   nil,
			wantErr: false,
		},
		{
			name:    "NoSchema",
			schema:  nil,
			value:   nil,
			wantErr: false,
		},
		{
			name:    "ValidationError",
			schema:  &jsv.Schema{Types: []string{"string"}},
			value:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := js.NewValidationSchema(tt.schema)
			err := vs.Validate(tt.value)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestValidationSchema_Types(t *testing.T) {
	tests := []struct {
		name   string
		schema *jsv.Schema
		want   []string
	}{
		{
			name:   "All",
			schema: nil,
			want: []string{
				"array",
				"boolean",
				"integer",
				"null",
				"number",
				"object",
				"string",
			},
		},
		{
			name:   "Intersection",
			schema: &jsv.Schema{Types: []string{"string"}},
			want:   []string{"string"},
		},
		{
			name: "AllOf",
			schema: &jsv.Schema{
				Types: []string{"string", "int", "null"},
				AllOf: []*jsv.Schema{
					{Types: []string{"string", "null"}},
				},
			},
			want: []string{"null", "string"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := js.NewValidationSchema(tt.schema)
			got := vs.Types()
			sort.Strings(got)
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}
