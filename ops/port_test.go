// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewPort(t *testing.T) {
	structType := reflect.TypeOf(struct{}{})
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{
				typ: "valid",
				st:  structType,
			},
			want: &Port{
				Type:       "valid",
				StructType: structType,
			},
		},
		{
			name: "EmptyType",
			args: args{
				typ: "",
				st:  structType,
			},
			wantErr: true,
		},
		{
			name: "NilStructType",
			args: args{
				typ: "nilstruct",
				st:  nil,
			},
			wantErr: true,
		},
		{
			name: "NonStructType",
			args: args{
				typ: "nonstruct",
				st:  reflect.TypeOf(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewPort(tt.args.typ, tt.args.st)
			assert.Equalf(t, tt.want, got, "%s", testhelpers.Diff(tt.want, got))
			assert.Equalf(t, tt.wantErr, gotErr != nil, "%s", testhelpers.Diff(tt.want, got))
		})
	}
}

func TestPort_NewStruct(t *testing.T) {
	var st = struct{}{}

	type fields struct {
		Type       string
		StructType reflect.Type
		Fields     PortFields
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "Success",
			fields: fields{
				Type:       "valid",
				StructType: reflect.TypeOf(struct{}{}),
				Fields:     nil,
			},
			want: &st,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Port{
				Type:       tt.fields.Type,
				StructType: tt.fields.StructType,
				Fields:     tt.fields.Fields,
			}
			gotPtr := p.NewStruct()
			if tt.want == nil && gotPtr != nil {
				assert.Equalf(t, tt.want, gotPtr, "NewStruct()")
			}
			got := *gotPtr
			assert.Equalf(t, tt.want, got, testhelpers.Diff(tt.want, got))
		})
	}
}

func TestPort_WithField(t *testing.T) {
	type fields struct {
		Type       string
		StructType reflect.Type
		Fields     PortFields
	}
	type args struct {
		f *PortField
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Port
	}{
		{
			name: "ValidField",
			fields: fields{
				Type:       "valid",
				StructType: reflect.TypeOf(struct{}{}),
				Fields:     nil,
			},
			args: args{
				f: &PortField{},
			},
			want: &Port{
				Type:       "valid",
				StructType: reflect.TypeOf(struct{}{}),
				Fields:     []*PortField{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Port{
				Type:       tt.fields.Type,
				StructType: tt.fields.StructType,
				Fields:     tt.fields.Fields,
			}
			assert.Equalf(t, tt.want, p.WithField(tt.args.f), "%s",
				testhelpers.Diff(tt.want, p))
		})
	}
}
