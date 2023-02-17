// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPortReflector_postProcessField(t *testing.T) {
	tests := []struct {
		name       string
		portStruct interface{}
		want       *ops.PortField
	}{
		{
			name: "Path",
			portStruct: struct {
				PathParameter string `req:"path" optional:"true"`
			}{},
			want: &ops.PortField{
				Name:     "PathParameter",
				Indices:  []int{0},
				Peer:     "pathParameter",
				Group:    FieldGroupHttpPath,
				Optional: false,
				PortType: PortTypeRequest,
				Type: ops.PortFieldType{
					Shape:        ops.FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"required": "true",
				},
				Baggage: map[interface{}]interface{}{},
			},
		},
		{
			name: "Body",
			portStruct: struct {
				Body any `req:"body"`
			}{},
			want: &ops.PortField{
				Name:     "Body",
				Indices:  []int{0},
				Peer:     "body",
				Group:    FieldGroupHttpBody,
				Optional: true,
				PortType: PortTypeRequest,
				Type: ops.PortFieldType{
					Shape:        ops.FieldShapeContent,
					Type:         ops.AnyType,
					Indirections: 0,
					HandlerType:  ops.AnyType,
					Optional:     true,
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PortReflector{}.ReflectInputPort(reflect.TypeOf(tt.portStruct))
			assert.NoError(t, err)
			assert.True(t,
				reflect.DeepEqual(tt.want, got.Fields[0]),
				testhelpers.Diff(tt.want, got.Fields[0]))
		})
	}
}
