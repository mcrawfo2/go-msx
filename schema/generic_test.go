// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetJsonFieldName(t *testing.T) {
	tests := []struct {
		name        string
		structField reflect.StructField
		want        types.Optional[string]
	}{
		{
			name: "Found",
			structField: reflect.StructField{
				Tag: `json:"found"`,
			},
			want: types.OptionalOf("found"),
		},
		{
			name: "Split",
			structField: reflect.StructField{
				Tag: `json:"split,omitempty"`,
			},
			want: types.OptionalOf("split"),
		},
		{
			name:        "NotFound",
			structField: reflect.StructField{},
			want:        types.OptionalEmpty[string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetJsonFieldName(tt.structField), "GetJsonFieldName(%v)", tt.structField)
		})
	}
}

func TestFindParameterizedStructField(t *testing.T) {
	type A struct {
		Field interface{} `inject:"payload"`
	}

	type B struct {
		A
	}

	type C struct {
		Struct A
	}

	tests := []struct {
		name       string
		structType reflect.Type
		wantResult []int
		wantName   string
		wantErr    bool
	}{
		{
			name:       "Flat",
			structType: reflect.TypeOf(A{}),
			wantResult: []int{0},
			wantName:   "payload",
		},
		{
			name:       "Anonymous",
			structType: reflect.TypeOf(B{}),
			wantResult: []int{0, 0},
			wantName:   "payload",
		},
		{
			name:       "Missing",
			structType: reflect.TypeOf(struct{}{}),
		},
		{
			name:       "NestedMissing",
			structType: reflect.TypeOf(C{}),
		},
		{
			name:       "Failure",
			structType: reflect.TypeOf(map[string]string{}),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotName, err := FindParameterizedStructField(tt.structType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else if err != nil {
				assert.NoError(t, err)
				return
			}
			assert.Equalf(t, tt.wantResult, gotResult, "FindParameterizedStructField(%v)", tt.structType)
			assert.Equalf(t, tt.wantName, gotName, "FindParameterizedStructField(%v)", tt.structType)
		})
	}
}

func TestNewParameterizedStruct(t *testing.T) {
	payloadInstance := map[string]string{}
	payloadType := reflect.TypeOf(payloadInstance)
	envelopeType := reflect.TypeOf(integration.MsxEnvelope{})

	found := false
	st := NewParameterizedStruct(envelopeType, payloadInstance)
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		if sf.Tag.Get("inject") != "Envelope" {
			continue
		}

		found = true
		assert.Equal(t, payloadType, sf.Type)
	}
	assert.True(t, found)
}