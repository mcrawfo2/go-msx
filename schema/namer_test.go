// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"reflect"
	"testing"
)

func TestGetNamedTypeName(t *testing.T) {
	tests := []struct {
		name         string
		instanceType reflect.Type
		want         string
		wantOk       bool
	}{
		{
			name:         "int",
			instanceType: reflect.TypeOf(0),
			want:         "",
			wantOk:       false,
		},
		{
			name:         "Anonymous",
			instanceType: reflect.TypeOf(struct{ A string }{}),
			want:         "",
			wantOk:       false,
		},
		{
			name:         "Named",
			instanceType: reflect.TypeOf(integration.MsxEnvelope{}),
			want:         "integration.MsxEnvelope",
			wantOk:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := GetNamedTypeName(tt.instanceType)
			if got != tt.want {
				t.Errorf("GetNamedTypeName() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetNamedTypeName() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
