// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
	"testing"
)

func TestSpringNamer_TypeInstanceName(t *testing.T) {
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
			if got := NewSpringTypeNamer().TypeInstanceName(tt.instance); got != tt.want {
				t.Error(testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestSpringNamer_TypeName(t *testing.T) {
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
			want:     "List«types.UUID»",
		},
		{
			name:     "Map[string]UUID",
			instance: map[string]types.UUID{},
			want:     "Map«string,types.UUID»",
		},
		{
			name:     "[]map[string]UUID",
			instance: []map[string]types.UUID{},
			want:     "List«Map«string,types.UUID»»",
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
			if got := NewSpringTypeNamer().TypeName(reflect.TypeOf(tt.instance)); got != tt.want {
				t.Error(testhelpers.Diff(tt.want, got))
			}
		})
	}
}
