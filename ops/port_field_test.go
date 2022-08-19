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

func TestPortField_WithOptional(t *testing.T) {
	pf := NewPortField("test", "peer", "group", true, PortFieldType{}, []int{0})
	pf.WithOptional(true)
	assert.True(t, pf.Optional)
}

type MyEnumerable struct{}

func (m *MyEnumerable) Enum() []interface{} {
	return []interface{}{
		"A", "B", "C",
	}
}

func TestPortField_Enum(t *testing.T) {
	tests := []struct {
		name string
		pf   PortField
		want []interface{}
	}{
		{
			name: "Primitive",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"enum": "A,B,C",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: []interface{}{
				"A", "B", "C",
			},
		},
		{
			name: "Enumer",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(MyEnumerable{}),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(MyEnumerable{}),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			want: []interface{}{
				"A", "B", "C",
			},
		},
		{
			name: "EnumerOverride",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(MyEnumerable{}),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(MyEnumerable{}),
				},
				Options: map[string]string{
					"enum": "D,E,F",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: []interface{}{
				"A", "B", "C",
			},
		},
		{
			name: "None",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"enum": "",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pf.Enum()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Enum diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortField_BoolOption(t *testing.T) {
	tests := []struct {
		name       string
		pf         PortField
		want       bool
		wantExists bool
	}{
		{
			name: "ExistsTrue",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"opt": "true",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want:       true,
			wantExists: true,
		},
		{
			name: "ExistsFalse",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"opt": "false",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want:       false,
			wantExists: true,
		},
		{
			name: "NotExists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			want:       false,
			wantExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotExists := tt.pf.BoolOption("opt")
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("BoolOption want diff\n%s", testhelpers.Diff(tt.want, got))
			}
			if !reflect.DeepEqual(tt.wantExists, gotExists) {
				t.Errorf("BoolOption wantExists diff\n%s", testhelpers.Diff(tt.wantExists, got))
			}
		})
	}
}

func TestPortField_WithBoolOptionDefault(t *testing.T) {
	type args struct {
		option string
		value  bool
	}

	tests := []struct {
		name string
		pf   PortField
		args args
		want bool
	}{
		{
			name: "Exists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"opt": "true",
				},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  false,
			},
			want: true,
		},
		{
			name: "NotExists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := tt.pf
			pf.WithBoolOptionDefault(tt.args.option, tt.args.value)

			got, gotExists := pf.BoolOption(tt.args.option)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("WithBoolOptionDefault want=%+v got=%+v", tt.want, got)
			}

			if !reflect.DeepEqual(true, gotExists) {
				t.Errorf("WithBoolOptionDefault wantExists=%+v gotExists=%+v", true, gotExists)
			}
		})
	}
}

func TestPortField_WithOptionDefault(t *testing.T) {
	type args struct {
		option string
		value  string
	}

	tests := []struct {
		name string
		pf   PortField
		args args
		want string
	}{
		{
			name: "Exists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"opt": "alpha",
				},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  "broccoli",
			},
			want: "alpha",
		},
		{
			name: "NotExists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  "broccoli",
			},
			want: "broccoli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := tt.pf
			pf.WithOptionDefault(tt.args.option, tt.args.value)

			got, gotExists := pf.Options[tt.args.option]

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("WithOptionDefault want=%+v got=%+v", tt.want, got)
			}

			if !reflect.DeepEqual(true, gotExists) {
				t.Errorf("WithOptionDefault wantExists=%+v gotExists=%+v", true, gotExists)
			}
		})
	}
}

func TestPortField_WithOption(t *testing.T) {
	type args struct {
		option string
		value  string
	}

	tests := []struct {
		name string
		pf   PortField
		args args
		want string
	}{
		{
			name: "Exists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"opt": "alpha",
				},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  "broccoli",
			},
			want: "broccoli",
		},
		{
			name: "NotExists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "opt",
				value:  "broccoli",
			},
			want: "broccoli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := tt.pf
			pf.WithOption(tt.args.option, tt.args.value)

			got, gotExists := pf.Options[tt.args.option]

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("WithOption want=%+v got=%+v", tt.want, got)
			}

			if !reflect.DeepEqual(true, gotExists) {
				t.Errorf("WithOption wantExists=%+v gotExists=%+v", true, gotExists)
			}
		})
	}
}

func TestPortField_WithBaggageItem(t *testing.T) {
	type args struct {
		option interface{}
		value  interface{}
	}

	tests := []struct {
		name string
		pf   PortField
		args args
		want string
	}{
		{
			name: "Exists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{
					"bag": "alpha",
				},
			},
			args: args{
				option: "bag",
				value:  "broccoli",
			},
			want: "broccoli",
		},
		{
			name: "NotExists",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			args: args{
				option: "bag",
				value:  "broccoli",
			},
			want: "broccoli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := tt.pf
			pf.WithBaggageItem(tt.args.option, tt.args.value)

			got, gotExists := pf.Baggage[tt.args.option]

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("WithBaggageItem want=%+v got=%+v", tt.want, got)
			}

			if !reflect.DeepEqual(true, gotExists) {
				t.Errorf("WithBaggageItem wantExists=%+v gotExists=%+v", true, gotExists)
			}
		})
	}
}

func TestPortField_ExpectShape(t *testing.T) {
	tests := []struct {
		name    string
		pf      PortField
		shape   string
		wantErr bool
	}{
		{
			name: "Matching",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			shape:   FieldShapePrimitive,
			wantErr: false,
		},
		{
			name: "NotMatching",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeUnknown,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			shape:   FieldShapePrimitive,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := tt.pf

			err := pf.ExpectShape(tt.shape)

			if !reflect.DeepEqual(tt.wantErr, err != nil) {
				t.Errorf("ExpectShape wantErr=%+v gotErr=%+v", tt.wantErr, err)
			}
		})
	}
}

func TestPortField_Default(t *testing.T) {
	tests := []struct {
		name string
		pf   PortField
		want interface{}
	}{
		{
			name: "Primitive",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"default": "A",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}("A"),
		},
		{
			name: "Array",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeArray,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"default": "A,B,C",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(
				[]string{
					"A", "B", "C",
				}),
		},
		{
			name: "Other",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeObject,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"default": "A",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(nil),
		},
		{
			name: "Undefined",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeObject,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pf.Default()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Default diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortField_Const(t *testing.T) {
	tests := []struct {
		name string
		pf   PortField
		want interface{}
	}{
		{
			name: "Primitive",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"const": "A",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}("A"),
		},
		{
			name: "Array",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeArray,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"const": "A,B,C",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(
				[]string{
					"A", "B", "C",
				}),
		},
		{
			name: "Other",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeObject,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{
					"const": "A",
				},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(nil),
		},
		{
			name: "Undefined",
			pf: PortField{
				Name:     "A",
				Indices:  []int{0},
				Peer:     "a",
				Group:    FieldGroupInjector,
				Optional: false,
				Type: PortFieldType{
					Shape:        FieldShapeObject,
					Type:         reflect.TypeOf(""),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(""),
				},
				Options: map[string]string{},
				Baggage: map[interface{}]interface{}{},
			},
			want: interface{}(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pf.Const()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Const diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestPortFieldHasGroup(t *testing.T) {
	type args struct {
		group   string
		pfgroup string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Matching",
			args: args{
				group:   "A",
				pfgroup: "A",
			},
			want: true,
		},
		{
			name: "NotMatching",
			args: args{
				group:   "B",
				pfgroup: "A",
			},
			want: false,
		},
		{
			name: "Empty",
			args: args{
				group:   "",
				pfgroup: "A",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			matcher := PortFieldHasGroup(tt.args.group)

			pf := NewPortField("", "", tt.args.pfgroup, true, PortFieldType{}, []int{0})
			matches := matcher(pf)

			assert.Equalf(t, tt.want, matches, "PortFieldHasGroup(%v)", tt.args.group)
		})
	}
}

func TestPortFieldHasName(t *testing.T) {
	type args struct {
		name   string
		pfname string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Matching",
			args: args{
				name:   "A",
				pfname: "A",
			},
			want: true,
		},
		{
			name: "NotMatching",
			args: args{
				name:   "B",
				pfname: "A",
			},
			want: false,
		},
		{
			name: "Empty",
			args: args{
				name:   "",
				pfname: "A",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			matcher := PortFieldHasName(tt.args.name)

			pf := NewPortField(tt.args.pfname, "", "", true, PortFieldType{}, []int{0})
			matches := matcher(pf)

			assert.Equalf(t, tt.want, matches, "PortFieldHasName(%v)", tt.args.name)
		})
	}
}

func TestPortFieldHasPeer(t *testing.T) {
	type args struct {
		peer   string
		pfpeer string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Matching",
			args: args{
				peer:   "A",
				pfpeer: "A",
			},
			want: true,
		},
		{
			name: "NotMatching",
			args: args{
				peer:   "B",
				pfpeer: "A",
			},
			want: false,
		},
		{
			name: "Empty",
			args: args{
				peer:   "",
				pfpeer: "A",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			matcher := PortFieldHasPeer(tt.args.peer)

			pf := NewPortField("", tt.args.pfpeer, "", true, PortFieldType{}, []int{0})
			matches := matcher(pf)

			assert.Equalf(t, tt.want, matches, "PortFieldHasName(%v)", tt.args.peer)
		})
	}
}

func TestPortFields_All(t *testing.T) {
	type args struct {
		predicates []PortFieldPredicate
	}
	tests := []struct {
		name string
		f    PortFields
		args args
		want PortFields
	}{
		{
			name: "MatchingOne",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
				{
					Name:  "C",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: PortFields{
				{
					Name:  "B",
					Group: "GroupB",
				},
			},
		},
		{
			name: "MatchingMultiple",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
				{
					Name:  "C",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: PortFields{
				{
					Name:  "B",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
			},
		},
		{
			name: "NotMatching",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
				{
					Name:  "C",
					Group: "GroupA",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.All(tt.args.predicates...)
			assert.Equal(t, tt.want, got, testhelpers.Diff(tt.want, got))
		})
	}
}

func TestPortFields_First(t *testing.T) {
	type args struct {
		predicates []PortFieldPredicate
	}
	tests := []struct {
		name string
		f    PortFields
		args args
		want *PortField
	}{
		{
			name: "MatchingOne",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
				{
					Name:  "C",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: &PortField{
				Name:  "B",
				Group: "GroupB",
			},
		},
		{
			name: "MatchingMultiple",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
				{
					Name:  "C",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
				{
					Name:  "B",
					Group: "GroupB",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: &PortField{
				Name:  "B",
				Group: "GroupB",
			},
		},
		{
			name: "NotMatching",
			f: PortFields{
				{
					Name:  "A",
					Group: "GroupB",
				},
				{
					Name:  "B",
					Group: "GroupC",
				},
				{
					Name:  "C",
					Group: "GroupA",
				},
			},
			args: args{
				predicates: []PortFieldPredicate{
					PortFieldHasGroup("GroupB"),
					PortFieldHasName("B"),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.First(tt.args.predicates...)
			assert.Equal(t, tt.want, got, testhelpers.Diff(tt.want, got))
		})
	}
}
