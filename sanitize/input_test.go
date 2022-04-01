// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sanitize

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
	"testing"
)

func TestInput(t *testing.T) {

	type inputStringsStruct struct {
		A string
		B *string
	}

	type inputDeepStruct struct {
		A string
		B inputStringsStruct
	}

	type args struct {
		value interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Pointer",
			args: args{
				value: types.NewStringPtr("a<b>c</b>d"),
			},
			want: types.NewStringPtr("acd"),
		},
		{
			name: "Struct",
			args: args{
				value: &inputStringsStruct{
					A: "a<b>c</b>d",
					B: types.NewStringPtr("a<b>c</b>d"),
				},
			},
			want: &inputStringsStruct{
				A: "acd",
				B: types.NewStringPtr("acd"),
			},
		},
		{
			name: "DeepStruct",
			args: args{
				value: &inputDeepStruct{
					A: "d<a>c</a>b",
					B: inputStringsStruct{
						A: "a<b>c</b>d",
						B: types.NewStringPtr("a<b>c</b>d"),
					},
				},
			},
			want: &inputDeepStruct{
				A: "dcb",
				B: inputStringsStruct{
					A:"acd",
					B: types.NewStringPtr("acd"),
				},
			},
		},
		{
			name: "DeepMapStruct",
			args: args{
				value: map[string]interface{}{
					"A": "d<a>c</a>b",
					"B": map[string]interface{}{
						"A": "a<b>c</b>d",
						"B": types.NewStringPtr("d<b>c</b>a"),
						"C": inputStringsStruct{
							A: "a<b>c</b>d",
							B: types.NewStringPtr("d</b>c</b>a"),
						},
					},
				},
			},
			want: map[string]interface{}{
				"A": "dcb",
				"B": map[string]interface{}{
					"A":"acd",
					"B": types.NewStringPtr("dca"),
					"C": inputStringsStruct{
						A: "acd",
						B: types.NewStringPtr("dca"),
					},
				},
			},
		},
		{
			name: "Slice",
			args: args{
				value: []string{
					"a<b>c</b>d",
					"d</b>c</b>a",
				},
			},
			want: []string{
				"acd",
				"dca",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := Options{Xss:true}
			if err := Input(tt.args.value, options); (err != nil) != tt.wantErr {
				t.Errorf("Input() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := tt.want
				if !reflect.DeepEqual(tt.args.value, want) {
					t.Error(testhelpers.Diff(want, tt.args.value))
				}
			}
		})
	}
}
