// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import "testing"

func TestStaticAlias_Alias(t *testing.T) {
	type fields struct {
		Path string
		File string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Swagger",
			fields: fields{
				Path: "/swagger",
				File: "/swagger-ui.html",
			},
			args: args{
				path: "/swagger",
			},
			want: "/swagger-ui.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := StaticAlias{
				Path: tt.fields.Path,
				File: tt.fields.File,
			}
			if got := a.Alias(tt.args.path); got != tt.want {
				t.Errorf("Alias() = %v, want %v", got, tt.want)
			}
		})
	}
}
