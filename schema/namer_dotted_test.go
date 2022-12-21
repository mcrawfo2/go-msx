// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDottedName_TypeInstanceName(t *testing.T) {
	tests := []struct {
		name     string
		instance interface{}
		want     string
	}{
		{
			name:     "UUID",
			instance: new(types.UUID),
			want:     "types.UUID",
		},
		{
			name: "Nil",
			want: VoidTypeName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDottedTypeNamer().TypeInstanceName(tt.instance); got != tt.want {
				t.Error(testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestDottedNamer_TypeName(t *testing.T) {
	tests := []struct {
		name     string
		instance interface{}
		want     string
	}{
		{
			name:     "UUID",
			instance: new(types.UUID),
			want:     "types.UUID",
		},
		{
			name:     "[]UUID",
			instance: []types.UUID{},
			want:     "types.UUID.List",
		},
		{
			name:     "map[string]UUID",
			instance: map[string]types.UUID{},
			want:     "types.UUID.Map",
		},
		{
			name:     "[]map[string]UUID",
			instance: []map[string]types.UUID{},
			want:     "types.UUID.Map.List",
		},
		{
			name:     "Time",
			instance: new(types.Time),
			want:     "types.Time",
		},
		{
			name:     "int",
			instance: 0,
			want:     "int",
		},
		{
			name: "Nil",
			want: VoidTypeName,
		},
		{
			name:     "Anonymous",
			instance: new(struct{ A string }),
			want:     ".anonymous1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDottedTypeNamer().TypeName(reflect.TypeOf(tt.instance)); got != tt.want {
				t.Error(testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestDottedTypeNamer_ParameterizedTypeName(t *testing.T) {
	tests := []struct {
		name        string
		wrapperType reflect.Type
		wrappedType reflect.Type
		want        string
	}{
		{
			name:        "EnvelopeV2",
			wrapperType: reflect.TypeOf(integration.MsxEnvelope{}),
			wrappedType: reflect.TypeOf(types.UUID{}),
			want:        "types.UUID.Envelope",
		},
		{
			name:        "PagingV2",
			wrapperType: reflect.TypeOf(paging.PaginatedResponse{}),
			wrappedType: reflect.TypeOf([]types.UUID{}),
			want:        "types.UUID.List.Page",
		},
		{
			name:        "PagingV8",
			wrapperType: reflect.TypeOf(paging.PaginatedResponseV8{}),
			wrappedType: reflect.TypeOf([]types.UUID{}),
			want:        "types.UUID.List.Page",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewDottedTypeNamer()
			got := n.ParameterizedTypeName(tt.wrapperType, tt.wrappedType)
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}
