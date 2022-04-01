// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
	"testing"
)

func Test_anyLiterals(t *testing.T) {
	type args struct {
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Handle null",
			args: args{
				values: []interface{}{
					nil,
				},
			},
			want: Null().GoString(),
		},
		{
			name: "Handle literal",
			args: args{
				values: []interface{}{
					"asc",
				},
			},
			want: "\"asc\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Statement(anyLiterals(tt.args.values)); !reflect.DeepEqual(got.GoString(), tt.want) {
				t.Errorf("anyLiterals = %v, want %v", got.GoString(), tt.want)
			}
		})
	}
}
