// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
	"testing"
)

func TestPortReflector_ReflectPortStruct_Type(t *testing.T) {
	type fields struct {
		FieldGroups        map[string]FieldGroup
		FieldPostProcessor PortFieldPostProcessorFunc
		FieldTypeReflector PortFieldTypeReflector
	}
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "SingleField",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test" extra:"test2"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Peer:     "singleField",
						Group:    "test",
						Indices:  []int{0},
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"extra": "test2",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "IgnoredNoTagField",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					FirstField  string `inp:"test" extra:"test2"`
					SecondField string `extra:"test3"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "FirstField",
						Peer:     "firstField",
						Group:    "test",
						Indices:  []int{0},
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"extra": "test2",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "IgnoredNameField",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					FirstField string `inp:"-"`
				}{}),
			},
			want: &Port{
				Type:   "inp",
				Fields: nil,
			},
			wantErr: false,
		},
		{
			name: "EmptyStruct",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st:  reflect.TypeOf(struct{}{}),
			},
			want: &Port{
				Type:   "inp",
				Fields: nil,
			},
			wantErr: false,
		},
		{
			name: "InvalidGroup",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					FirstField string `inp:"test2"`
				}{}),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PortReflector{
				FieldGroups:        tt.fields.FieldGroups,
				FieldPostProcessor: tt.fields.FieldPostProcessor,
				FieldTypeReflector: tt.fields.FieldTypeReflector,
			}
			if tt.want != nil {
				tt.want.StructType = tt.args.st
			}
			got, err := r.ReflectPortStruct(tt.args.typ, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReflectPortStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ReflectPortStruct() diff\n%s",
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortReflector_ReflectPortStruct_Name(t *testing.T) {
	type fields struct {
		FieldGroups        map[string]FieldGroup
		FieldPostProcessor PortFieldPostProcessorFunc
		FieldTypeReflector PortFieldTypeReflector
	}
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "DefaultNaming",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CustomName",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test=bob"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "bob",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PortReflector{
				FieldGroups:        tt.fields.FieldGroups,
				FieldPostProcessor: tt.fields.FieldPostProcessor,
				FieldTypeReflector: tt.fields.FieldTypeReflector,
			}
			if tt.want != nil {
				tt.want.StructType = tt.args.st
			}
			got, err := r.ReflectPortStruct(tt.args.typ, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReflectPortStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ReflectPortStruct() diff\n%s",
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortReflector_ReflectPortStruct_PrimaryTag(t *testing.T) {
	type fields struct {
		FieldGroups        map[string]FieldGroup
		FieldPostProcessor PortFieldPostProcessorFunc
		FieldTypeReflector PortFieldTypeReflector
	}
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "DefaultOptionValue",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test,custom"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"custom": "true",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CustomOptionValue",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test,custom=value"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"custom": "value",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PortReflector{
				FieldGroups:        tt.fields.FieldGroups,
				FieldPostProcessor: tt.fields.FieldPostProcessor,
				FieldTypeReflector: tt.fields.FieldTypeReflector,
			}
			if tt.want != nil {
				tt.want.StructType = tt.args.st
			}
			got, err := r.ReflectPortStruct(tt.args.typ, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReflectPortStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ReflectPortStruct() diff\n%s",
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortReflector_ReflectPortStruct_Optional(t *testing.T) {
	type fields struct {
		FieldGroups        map[string]FieldGroup
		FieldPostProcessor PortFieldPostProcessorFunc
		FieldTypeReflector PortFieldTypeReflector
	}
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "Unspecified",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PrimaryRequired",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test,required"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"required": "true",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PrimaryOptional",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test,optional"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: true,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"required": "false",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "SecondaryRequired",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test" required:"true"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"required": "true",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "SecondaryOptional",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test" optional:"true"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: true,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{
							"optional": "true",
						},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PortReflector{
				FieldGroups:        tt.fields.FieldGroups,
				FieldPostProcessor: tt.fields.FieldPostProcessor,
				FieldTypeReflector: tt.fields.FieldTypeReflector,
			}
			if tt.want != nil {
				tt.want.StructType = tt.args.st
			}
			got, err := r.ReflectPortStruct(tt.args.typ, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReflectPortStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ReflectPortStruct() diff\n%s",
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortReflector_ReflectPortStruct_Indices(t *testing.T) {
	type embeddedStruct struct {
		AnotherField string `inp:"test"`
	}

	type fields struct {
		FieldGroups        map[string]FieldGroup
		FieldPostProcessor PortFieldPostProcessorFunc
		FieldTypeReflector PortFieldTypeReflector
	}
	type args struct {
		typ string
		st  reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		wantErr bool
	}{
		{
			name: "Flat",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField  string `inp:"test"`
					AnotherField string `inp:"test"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
					{
						Name:     "AnotherField",
						Indices:  []int{1},
						Peer:     "anotherField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Nested",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string         `inp:"test"`
					SecondField embeddedStruct `inp:"test"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
					{
						Name:     "SecondField",
						Indices:  []int{1},
						Peer:     "secondField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapeObject,
							Type:        reflect.TypeOf(embeddedStruct{}),
							HandlerType: reflect.TypeOf(embeddedStruct{}),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "NestedPointer",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string          `inp:"test"`
					SecondField *embeddedStruct `inp:"test"`
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
					{
						Name:     "SecondField",
						Indices:  []int{1},
						Peer:     "secondField",
						Group:    "test",
						Optional: true,
						Type: PortFieldType{
							Shape:        FieldShapeObject,
							Type:         reflect.TypeOf(embeddedStruct{}),
							Indirections: 1,
							HandlerType:  reflect.TypeOf(embeddedStruct{}),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "AnonymousStruct",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test"`
					embeddedStruct
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
					{
						Name:     "AnotherField",
						Indices:  []int{1, 0},
						Peer:     "anotherField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "AnonymousStructPtr",
			fields: fields{
				FieldGroups: map[string]FieldGroup{
					FieldGroupTest: {
						Cardinality: types.CardinalityZeroToMany(),
					},
				},
			},
			args: args{
				typ: "inp",
				st: reflect.TypeOf(struct {
					SingleField string `inp:"test"`
					*embeddedStruct
				}{}),
			},
			want: &Port{
				Type: "inp",
				Fields: []*PortField{
					{
						Name:     "SingleField",
						Indices:  []int{0},
						Peer:     "singleField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
					{
						Name:     "AnotherField",
						Indices:  []int{1, 0},
						Peer:     "anotherField",
						Group:    "test",
						Optional: false,
						Type: PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
						Options: map[string]string{},
						Baggage: map[interface{}]interface{}{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PortReflector{
				FieldGroups:        tt.fields.FieldGroups,
				FieldPostProcessor: tt.fields.FieldPostProcessor,
				FieldTypeReflector: tt.fields.FieldTypeReflector,
			}
			if tt.want != nil {
				tt.want.StructType = tt.args.st
			}
			got, err := r.ReflectPortStruct(tt.args.typ, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReflectPortStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"ReflectPortStruct() diff\n%s",
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}
