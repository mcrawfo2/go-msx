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

func TestArraySchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Array)
	wantItems := NewSchemaPtr(jsonschema.Object).ToSchemaOrBool()
	want.WithItems(jsonschema.Items{
		SchemaArray: []jsonschema.SchemaOrBool{wantItems},
	})
	got := ArraySchema(*NewSchemaPtr(jsonschema.Object))
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestBooleanSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Boolean)
	got := BooleanSchema()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestIntegerSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Integer)
	got := IntegerSchema()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestMapSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Object)
	want.WithAdditionalProperties(NewSchemaPtr(jsonschema.Object).ToSchemaOrBool())
	got := MapSchema(NewSchemaPtr(jsonschema.Object))
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestNewSchemaPtr(t *testing.T) {
	tests := []struct {
		name string
		st   jsonschema.SimpleType
	}{
		{
			name: "Array",
			st:   jsonschema.Array,
		},
		{
			name: "Boolean",
			st:   jsonschema.Boolean,
		},
		{
			name: "Integer",
			st:   jsonschema.Integer,
		},
		{
			name: "Null",
			st:   jsonschema.Null,
		},
		{
			name: "Number",
			st:   jsonschema.Number,
		},
		{
			name: "Object",
			st:   jsonschema.Object,
		},
		{
			name: "String",
			st:   jsonschema.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSchemaPtr(tt.st)
			want := new(jsonschema.Schema).WithType(*NewType(tt.st))
			assert.True(t,
				reflect.DeepEqual(got, want),
				testhelpers.Diff(got, want))
		})
	}
}

func TestNewType(t *testing.T) {
	tests := []struct {
		name string
		st   jsonschema.SimpleType
	}{
		{
			name: "Array",
			st:   jsonschema.Array,
		},
		{
			name: "Boolean",
			st:   jsonschema.Boolean,
		},
		{
			name: "Integer",
			st:   jsonschema.Integer,
		},
		{
			name: "Null",
			st:   jsonschema.Null,
		},
		{
			name: "Number",
			st:   jsonschema.Number,
		},
		{
			name: "Object",
			st:   jsonschema.Object,
		},
		{
			name: "String",
			st:   jsonschema.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stp := &tt.st
			want := stp.Type()
			got := NewType(tt.st)

			assert.True(t,
				reflect.DeepEqual(got, &want),
				testhelpers.Diff(&want, got))
		})
	}
}

func TestNumberSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Number)
	got := NumberSchema()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestObjectSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.Object)
	got := ObjectSchema()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringSchema(t *testing.T) {
	want := NewSchemaPtr(jsonschema.String)
	got := StringSchema()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
