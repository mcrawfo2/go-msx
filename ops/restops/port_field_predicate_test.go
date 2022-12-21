// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPortFieldIsErrorHeader(t *testing.T) {
	tests := []struct {
		name string
		pf   *ops.PortField
		want bool
	}{
		{
			name: "True",
			pf: &ops.PortField{
				Group: FieldGroupHttpHeader,
				Options: map[string]string{
					"error": "true",
				},
			},
			want: true,
		},
		{
			name: "False",
			pf: &ops.PortField{
				Group: FieldGroupHttpHeader,
				Options: map[string]string{
					"error": "false",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PortFieldIsErrorHeader(tt.pf)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPortFieldIsError(t *testing.T) {
	tests := []struct {
		name string
		pf   *ops.PortField
		want bool
	}{
		{
			name: "True",
			pf: &ops.PortField{
				Options: map[string]string{
					"error": "true",
				},
			},
			want: true,
		},
		{
			name: "False",
			pf: &ops.PortField{
				Options: map[string]string{
					"error": "false",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PortFieldIsError(tt.pf)
			assert.Equal(t, tt.want, got)
		})
	}
}
