// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"testing"
)

func TestPopulateEnumFromTags(t *testing.T) {
	tests := []struct {
		name    string
		s       *jsonschema.Schema
		tag     string
		want    []interface{}
		wantErr bool
	}{
		{
			name: "Strings",
			s:    NewSchemaPtr(jsonschema.String),
			tag:  `enum:"A,B,C,D"`,
			want: []interface{}{"A", "B", "C", "D"},
		},
		{
			name: "Integers",
			s:    NewSchemaPtr(jsonschema.Integer),
			tag:  `enum:"1,2,3,4"`,
			want: []interface{}{int64(1), int64(2), int64(3), int64(4)},
		},
		{
			name: "Booleans",
			s:    NewSchemaPtr(jsonschema.Boolean),
			tag:  `enum:"true,false"`,
			want: []interface{}{true, false},
		},
		{
			name: "Numbers",
			s:    NewSchemaPtr(jsonschema.Number),
			tag:  `enum:"1.1,1.2,1.3"`,
			want: []interface{}{1.1, 1.2, 1.3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := reflect.StructTag(tt.tag)
			err := PopulateEnumFromTags(tt.s, st)
			assert.Equal(t, tt.wantErr, err != nil)

			got := tt.s.Enum
			want := tt.want
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(want, got),
					testhelpers.Diff(want, got))
			}
		})
	}
}

func TestPopulateFieldsFromTags(t *testing.T) {
	tests := []struct {
		name    string
		s       *jsonschema.Schema
		tag     string
		want    *jsonschema.Schema
		wantErr bool
	}{
		{
			name: "Description",
			s:    NewSchemaPtr(jsonschema.String),
			tag:  `description:"some-description"`,
			want: NewSchemaPtr(jsonschema.String).WithDescription("some-description"),
		},
		{
			name: "Maximum",
			s:    NewSchemaPtr(jsonschema.Number),
			tag:  `maximum:"10"`,
			want: NewSchemaPtr(jsonschema.Number).WithMaximum(10),
		},
		{
			name: "Enum",
			s:    NewSchemaPtr(jsonschema.Integer),
			tag:  `enum:"1,2,3"`,
			want: NewSchemaPtr(jsonschema.Integer).WithEnum(int64(1), int64(2), int64(3)),
		},
		{
			name: "Const",
			s:    NewSchemaPtr(jsonschema.Boolean),
			tag:  `const:"false"`,
			want: NewSchemaPtr(jsonschema.Boolean).WithConst(false),
		},
		{
			name: "Default",
			s:    NewSchemaPtr(jsonschema.Boolean),
			tag:  `default:"true"`,
			want: NewSchemaPtr(jsonschema.Boolean).WithDefault(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := reflect.StructTag(tt.tag)
			err := PopulateFieldsFromTags(tt.s, st)
			assert.Equal(t, tt.wantErr, err != nil)

			got := tt.s
			want := tt.want
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(want, got),
					testhelpers.Diff(want, got))
			}
		})
	}
}

func TestPopulateInterfaceFieldsFromTags(t *testing.T) {
	tests := []struct {
		name    string
		s       *jsonschema.Schema
		tag     string
		want    *jsonschema.Schema
		wantErr bool
	}{
		{
			name: "Enum",
			s:    NewSchemaPtr(jsonschema.Integer),
			tag:  `enum:"1,2,3"`,
			want: NewSchemaPtr(jsonschema.Integer),
		},
		{
			name: "ConstBoolean",
			s:    NewSchemaPtr(jsonschema.Boolean),
			tag:  `const:"false"`,
			want: NewSchemaPtr(jsonschema.Boolean).WithConst(false),
		},
		{
			name: "DefaultBoolean",
			s:    NewSchemaPtr(jsonschema.Boolean),
			tag:  `default:"true"`,
			want: NewSchemaPtr(jsonschema.Boolean).WithDefault(true),
		},
		{
			name: "ConstString",
			s:    NewSchemaPtr(jsonschema.String),
			tag:  `const:"false"`,
			want: NewSchemaPtr(jsonschema.String).WithConst("false"),
		},
		{
			name: "DefaultString",
			s:    NewSchemaPtr(jsonschema.String),
			tag:  `default:"true"`,
			want: NewSchemaPtr(jsonschema.String).WithDefault("true"),
		},
		{
			name: "ConstNumber",
			s:    NewSchemaPtr(jsonschema.Number),
			tag:  `const:"3.14159"`,
			want: NewSchemaPtr(jsonschema.Number).WithConst(3.14159),
		},
		{
			name: "DefaultNumber",
			s:    NewSchemaPtr(jsonschema.Number),
			tag:  `default:"2.71828"`,
			want: NewSchemaPtr(jsonschema.Number).WithDefault(2.71828),
		},
		{
			name: "Integer",
			s:    NewSchemaPtr(jsonschema.Integer),
			tag:  `const:"360"`,
			want: NewSchemaPtr(jsonschema.Integer).WithConst(int64(360)),
		},
		{
			name: "DefaultNumber",
			s:    NewSchemaPtr(jsonschema.Integer),
			tag:  `default:"42"`,
			want: NewSchemaPtr(jsonschema.Integer).WithDefault(int64(42)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := reflect.StructTag(tt.tag)
			err := PopulateInterfaceFieldsFromTags(tt.s, st)
			assert.Equal(t, tt.wantErr, err != nil)

			got := tt.s
			want := tt.want
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(want, got),
					testhelpers.Diff(want, got))
			}
		})
	}
}
