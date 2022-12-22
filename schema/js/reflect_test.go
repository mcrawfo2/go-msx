// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"testing"
)

type TestStruct struct{}

type TestNamedStruct struct {
	A *string `required:"true"`
	B string  `optional:"true"`
}

func (t TestNamedStruct) JSONSchemaDefName() string {
	return "NamedStruct"
}

func (t TestNamedStruct) Example() interface{} {
	return map[string]interface{}{
		"a": "eg",
		"b": nil,
	}
}

func TestDefNameInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		verify func(*jsonschema.ReflectContext)
	}{
		{
			name: "NotExposer",
			verify: func(rc *jsonschema.ReflectContext) {
				fn := rc.DefName
				got := fn(reflect.TypeOf(TestStruct{}), "TestStruct")
				assert.Equal(t, "js.TestStruct", got)
			},
		},
		{
			name: "Exposer",
			verify: func(rc *jsonschema.ReflectContext) {
				fn := rc.DefName
				got := fn(reflect.TypeOf(TestNamedStruct{}), "TestNamedStruct")
				assert.Equal(t, "NamedStruct", got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &jsonschema.ReflectContext{}
			DefNameInterceptor()(rc)
			tt.verify(rc)
		})
	}
}

func TestTypeTitleInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		verify func(*jsonschema.ReflectContext)
	}{
		{
			name: "WithId",
			verify: func(rc *jsonschema.ReflectContext) {
				s := &jsonschema.Schema{}
				s.WithID("schema-id")
				fn := rc.InterceptType
				cont, err := fn(reflect.ValueOf(TestNamedStruct{}), s)
				assert.False(t, cont)
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(s.ID, s.Title),
					testhelpers.Diff(s.ID, s.Title))
			},
		},
		{
			name: "NoId",
			verify: func(rc *jsonschema.ReflectContext) {
				s := &jsonschema.Schema{}
				fn := rc.InterceptType
				cont, err := fn(reflect.ValueOf(TestNamedStruct{}), s)
				want := types.NewStringPtr("TestNamedStruct")
				assert.False(t, cont)
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(want, s.Title),
					testhelpers.Diff(want, s.Title))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &jsonschema.ReflectContext{}
			TypeTitleInterceptor()(rc)
			tt.verify(rc)
		})
	}
}

func TestFindRequiredJsonFields(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  []string
	}{
		{
			name: "AllFields",
			value: struct {
				A string
			}{},
			want: []string{"a"},
		},
		{
			name: "SkipField",
			value: struct {
				A string `json:"-"`
				B string
			}{},
			want: []string{"b"},
		},
		{
			name: "ExplicitRequiredField",
			value: struct {
				A string  `required:"false"`
				B *string `required:"true"`
			}{},
			want: []string{"b"},
		},
		{
			name: "ExplicitOptionalField",
			value: struct {
				A string  `optional:"true"`
				B *string `optional:"false"`
			}{},
			want: []string{"b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRequiredJsonFields(reflect.TypeOf(tt.value))
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestStructRequiredInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		verify func(*jsonschema.ReflectContext)
	}{
		{
			name: "Required",
			verify: func(rc *jsonschema.ReflectContext) {
				s := &jsonschema.Schema{}
				fn := rc.InterceptType
				cont, err := fn(reflect.ValueOf(TestNamedStruct{}), s)
				want := []string{"a"}
				got := s.Required
				assert.False(t, cont)
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(want, got),
					testhelpers.Diff(want, got))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &jsonschema.ReflectContext{}
			StructRequiredInterceptor()(rc)
			tt.verify(rc)
		})
	}
}

func TestExampleInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		verify func(*jsonschema.ReflectContext)
	}{
		{
			name: "Required",
			verify: func(rc *jsonschema.ReflectContext) {
				s := &jsonschema.Schema{}
				fn := rc.InterceptType
				cont, err := fn(reflect.ValueOf(TestNamedStruct{}), s)
				assert.False(t, cont)
				assert.NoError(t, err)

				want := TestNamedStruct{}.Example()
				assert.Len(t, s.Examples, 1)
				got := s.Examples[0]

				assert.True(t,
					reflect.DeepEqual(want, got),
					testhelpers.Diff(want, got))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &jsonschema.ReflectContext{}
			ExampleInterceptor()(rc)
			tt.verify(rc)
		})
	}
}
