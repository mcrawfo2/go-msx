// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"reflect"
	"testing"
)

type MyTextType string

func (m MyTextType) Validate() error {
	if string(m) == "error" {
		return errors.New(string(m))
	}
	return nil
}

func (m *MyTextType) UnmarshalText(data string) error {
	*m = MyTextType(data)
	return nil
}

func (m MyTextType) MarshalText() (string, error) {
	return string(m), nil
}

func NewMyTextType(data string) *MyTextType {
	v := MyTextType(data)
	return &v
}

type ErrorReader struct{}

func (e ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Failure")
}

func (e ErrorReader) Close() error {
	return nil
}

func TestPortFieldInjector_InjectPrimitive_Text(t *testing.T) {
	type primitives struct {
		A string       `test:"injector"`
		B *string      `test:"injector"`
		C **string     `test:"injector"`
		G []byte       `test:"injector"`
		H *[]byte      `test:"injector"`
		D **[]byte     `test:"injector"`
		I []rune       `test:"injector"`
		J *[]rune      `test:"injector"`
		E **[]rune     `test:"injector"`
		K MyTextType   `test:"injector"`
		L *MyTextType  `test:"injector"`
		F **MyTextType `test:"injector"`
	}

	pr := PortReflector{
		Direction: PortDirectionIn,
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldPostProcessor: nil,
		FieldTypeReflector: DefaultPortFieldTypeReflector{
			Direction: PortDirectionIn,
		},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))

	tests := []struct {
		name    string
		field   *PortField
		value   string
		want    primitives
		wantErr bool
	}{
		{
			name:  "String",
			field: port.Fields.First(PortFieldHasName("A")),
			value: "abc",
			want: primitives{
				A: "abc",
			},
			wantErr: false,
		},
		{
			name:  "StringIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: "abc",
			want: primitives{
				B: types.NewStringPtr("abc"),
			},
			wantErr: false,
		},
		{
			name:  "StringDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "abc",
			want: primitives{
				C: func() **string {
					v := types.NewStringPtr("abc")
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "ByteSlice",
			field: port.Fields.First(PortFieldHasName("G")),
			value: "abc",
			want: primitives{
				G: []byte("abc"),
			},
			wantErr: false,
		},
		{
			name:  "ByteSliceIndirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: "abc",
			want: primitives{
				H: types.NewByteSlicePtr([]byte("abc")),
			},
			wantErr: false,
		},
		{
			name:  "ByteSliceDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("D")),
			value: "abc",
			want: primitives{
				D: func() **[]byte {
					v := types.NewByteSlicePtr([]byte("abc"))
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "RuneSlice",
			field: port.Fields.First(PortFieldHasName("I")),
			value: "abc",
			want: primitives{
				I: []rune("abc"),
			},
			wantErr: false,
		},
		{
			name:  "RuneSliceIndirect",
			field: port.Fields.First(PortFieldHasName("J")),
			value: "abc",
			want: primitives{
				J: types.NewRuneSlicePtr([]rune("abc")),
			},
			wantErr: false,
		},
		{
			name:  "RuneSliceDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: "abc",
			want: primitives{
				E: func() **[]rune {
					v := types.NewRuneSlicePtr([]rune("abc"))
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "TextUnmarshaler",
			field: port.Fields.First(PortFieldHasName("K")),
			value: "abc",
			want: primitives{
				K: "abc",
			},
			wantErr: false,
		},
		{
			name:  "TextUnmarshalerIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: "abc",
			want: primitives{
				L: NewMyTextType("abc"),
			},
			wantErr: false,
		},
		{
			name:  "TextUnmarshalerDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: "abc",
			want: primitives{
				F: func() **MyTextType {
					v := NewMyTextType("abc")
					return &v
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectPrimitive(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectPrimitive() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectPrimitive_Int(t *testing.T) {
	type primitives struct {
		A int     `test:"injector"`
		B *int    `test:"injector"`
		C **int   `test:"injector"`
		D int8    `test:"injector"`
		E *int8   `test:"injector"`
		F **int8  `test:"injector"`
		G int16   `test:"injector"`
		H *int16  `test:"injector"`
		I **int16 `test:"injector"`
		J int32   `test:"injector"`
		K *int32  `test:"injector"`
		L **int32 `test:"injector"`
		M int64   `test:"injector"`
		N *int64  `test:"injector"`
		O **int64 `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))

	tests := []struct {
		name    string
		field   *PortField
		value   string
		want    primitives
		wantErr bool
	}{
		{
			name:  "Int",
			field: port.Fields.First(PortFieldHasName("A")),
			value: "8",
			want: primitives{
				A: int(8),
			},
			wantErr: false,
		},
		{
			name:    "IntInvalid",
			field:   port.Fields.First(PortFieldHasName("A")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "IntIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: "9",
			want: primitives{
				B: types.NewIntPtr(9),
			},
			wantErr: false,
		},
		{
			name:  "IntDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "10",
			want: primitives{
				C: func() **int {
					v := types.NewIntPtr(10)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int8",
			field: port.Fields.First(PortFieldHasName("D")),
			value: "8",
			want: primitives{
				D: int8(8),
			},
			wantErr: false,
		},
		{
			name:    "Int8Invalid",
			field:   port.Fields.First(PortFieldHasName("D")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Int8Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: "9",
			want: primitives{
				E: func() *int8 {
					v := int8(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int8DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: "10",
			want: primitives{
				F: func() **int8 {
					v := int8(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int16",
			field: port.Fields.First(PortFieldHasName("G")),
			value: "8",
			want: primitives{
				G: int16(8),
			},
			wantErr: false,
		},
		{
			name:    "Int16Invalid",
			field:   port.Fields.First(PortFieldHasName("G")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Int16Indirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: "9",
			want: primitives{
				H: func() *int16 {
					v := int16(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int16DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("I")),
			value: "10",
			want: primitives{
				I: func() **int16 {
					v := int16(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int32",
			field: port.Fields.First(PortFieldHasName("J")),
			value: "8",
			want: primitives{
				J: int32(8),
			},
			wantErr: false,
		},
		{
			name:    "Int32Invalid",
			field:   port.Fields.First(PortFieldHasName("J")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Int32Indirect",
			field: port.Fields.First(PortFieldHasName("K")),
			value: "9",
			want: primitives{
				K: func() *int32 {
					v := int32(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: "10",
			want: primitives{
				L: func() **int32 {
					v := int32(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		}, {
			name:  "Int64",
			field: port.Fields.First(PortFieldHasName("M")),
			value: "8",
			want: primitives{
				M: int64(8),
			},
			wantErr: false,
		},
		{
			name:    "Int64Invalid",
			field:   port.Fields.First(PortFieldHasName("M")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Int64Indirect",
			field: port.Fields.First(PortFieldHasName("N")),
			value: "9",
			want: primitives{
				N: func() *int64 {
					v := int64(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Int64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("O")),
			value: "10",
			want: primitives{
				O: func() **int64 {
					v := int64(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectPrimitive(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectPrimitive() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectPrimitive_Uint(t *testing.T) {
	type primitives struct {
		A uint     `test:"injector"`
		B *uint    `test:"injector"`
		C **uint   `test:"injector"`
		D uint8    `test:"injector"`
		E *uint8   `test:"injector"`
		F **uint8  `test:"injector"`
		G uint16   `test:"injector"`
		H *uint16  `test:"injector"`
		I **uint16 `test:"injector"`
		J uint32   `test:"injector"`
		K *uint32  `test:"injector"`
		L **uint32 `test:"injector"`
		M uint64   `test:"injector"`
		N *uint64  `test:"injector"`
		O **uint64 `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))

	tests := []struct {
		name    string
		field   *PortField
		value   string
		want    primitives
		wantErr bool
	}{
		{
			name:  "Uint",
			field: port.Fields.First(PortFieldHasName("A")),
			value: "8",
			want: primitives{
				A: uint(8),
			},
			wantErr: false,
		},
		{
			name:    "UintInvalid",
			field:   port.Fields.First(PortFieldHasName("A")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "UintIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: "9",
			want: primitives{
				B: types.NewUintPtr(9),
			},
			wantErr: false,
		},
		{
			name:  "UintDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "10",
			want: primitives{
				C: func() **uint {
					v := types.NewUintPtr(10)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint8",
			field: port.Fields.First(PortFieldHasName("D")),
			value: "8",
			want: primitives{
				D: uint8(8),
			},
			wantErr: false,
		},
		{
			name:    "Uint8Invalid",
			field:   port.Fields.First(PortFieldHasName("D")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Uint8Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: "9",
			want: primitives{
				E: func() *uint8 {
					v := uint8(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint8DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: "10",
			want: primitives{
				F: func() **uint8 {
					v := uint8(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint16",
			field: port.Fields.First(PortFieldHasName("G")),
			value: "8",
			want: primitives{
				G: uint16(8),
			},
			wantErr: false,
		},
		{
			name:    "Uint16Invalid",
			field:   port.Fields.First(PortFieldHasName("G")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Uint16Indirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: "9",
			want: primitives{
				H: func() *uint16 {
					v := uint16(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint16DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("I")),
			value: "10",
			want: primitives{
				I: func() **uint16 {
					v := uint16(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint32",
			field: port.Fields.First(PortFieldHasName("J")),
			value: "8",
			want: primitives{
				J: uint32(8),
			},
			wantErr: false,
		},
		{
			name:    "Uint32Invalid",
			field:   port.Fields.First(PortFieldHasName("J")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Uint32Indirect",
			field: port.Fields.First(PortFieldHasName("K")),
			value: "9",
			want: primitives{
				K: func() *uint32 {
					v := uint32(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: "10",
			want: primitives{
				L: func() **uint32 {
					v := uint32(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		}, {
			name:  "Uint64",
			field: port.Fields.First(PortFieldHasName("M")),
			value: "8",
			want: primitives{
				M: uint64(8),
			},
			wantErr: false,
		},
		{
			name:    "Uint64Invalid",
			field:   port.Fields.First(PortFieldHasName("M")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Uint64Indirect",
			field: port.Fields.First(PortFieldHasName("N")),
			value: "9",
			want: primitives{
				N: func() *uint64 {
					v := uint64(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Uint64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("O")),
			value: "10",
			want: primitives{
				O: func() **uint64 {
					v := uint64(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectPrimitive(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectPrimitive() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectPrimitive_Float(t *testing.T) {
	type primitives struct {
		A float32   `test:"injector"`
		B *float32  `test:"injector"`
		C **float32 `test:"injector"`
		D float64   `test:"injector"`
		E *float64  `test:"injector"`
		F **float64 `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))

	tests := []struct {
		name    string
		field   *PortField
		value   string
		want    primitives
		wantErr bool
	}{
		{
			name:  "Float32",
			field: port.Fields.First(PortFieldHasName("A")),
			value: "8",
			want: primitives{
				A: float32(8),
			},
			wantErr: false,
		},
		{
			name:    "Float32Invalid",
			field:   port.Fields.First(PortFieldHasName("A")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Float32Indirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: "9",
			want: primitives{
				B: func() *float32 {
					v := float32(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Float32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "10",
			want: primitives{
				C: func() **float32 {
					v := float32(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Float64",
			field: port.Fields.First(PortFieldHasName("D")),
			value: "8",
			want: primitives{
				D: float64(8),
			},
			wantErr: false,
		},
		{
			name:    "Float64Invalid",
			field:   port.Fields.First(PortFieldHasName("D")),
			value:   "s",
			wantErr: true,
		},
		{
			name:  "Float64Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: "9",
			want: primitives{
				E: func() *float64 {
					v := float64(9)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Float64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: "10",
			want: primitives{
				F: func() **float64 {
					v := float64(10)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectPrimitive(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectPrimitive() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectPrimitive_Bool(t *testing.T) {
	type primitives struct {
		C bool  `test:"injector"`
		D *bool `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))

	tests := []struct {
		name    string
		field   *PortField
		value   string
		want    primitives
		wantErr bool
	}{
		{
			name:  "Boolean",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "true",
			want: primitives{
				C: true,
			},
			wantErr: false,
		},
		{
			name:  "BooleanInvalid",
			field: port.Fields.First(PortFieldHasName("C")),
			value: "something",
			want: primitives{
				C: false,
			},
			wantErr: true,
		},
		{
			name:  "BooleanIndirect",
			field: port.Fields.First(PortFieldHasName("D")),
			value: "true",
			want: primitives{
				D: types.NewBoolPtr(true),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectPrimitive(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectPrimitive() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectPrimitive_Error(t *testing.T) {
	type primitives struct {
		A chan struct{} `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	_, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))
	assert.Error(t, err)
}

func TestPortFieldInjector_InjectContent_Bytes(t *testing.T) {
	type contents struct {
		A []byte            `test:"injector"`
		B *[]byte           `test:"injector"`
		C **[]byte          `test:"injector"`
		D json.RawMessage   `test:"injector"`
		E *json.RawMessage  `test:"injector"`
		F **json.RawMessage `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(contents{}))

	tests := []struct {
		name    string
		field   *PortField
		value   Content
		want    contents
		wantErr bool
	}{
		{
			name:  "ByteArray",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				A: []byte(`"abc"`),
			},
			wantErr: false,
		},
		{
			name:  "ByteArrayIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				B: func() *[]byte {
					v := []byte(`"abc"`)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "ByteArrayDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				C: func() **[]byte {
					v := []byte(`"abc"`)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "ByteArrayReaderError",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromReadCloser(
				NewContentOptions("text/plain"),
				new(ErrorReader),
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs contents
			tt.field.Type.Shape = FieldShapeContent
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectContent(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectContent() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectContent_Runes(t *testing.T) {
	type contents struct {
		A []rune   `test:"injector"`
		B *[]rune  `test:"injector"`
		C **[]rune `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(contents{}))

	tests := []struct {
		name    string
		field   *PortField
		value   Content
		want    contents
		wantErr bool
	}{
		{
			name:  "RuneArray",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				A: []rune(`"abc"`),
			},
			wantErr: false,
		},
		{
			name:  "RuneArrayIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				B: func() *[]rune {
					v := []rune(`"abc"`)
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "RuneArrayDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				C: func() **[]rune {
					v := []rune(`"abc"`)
					w := &v
					return &w
				}(),
			},
			wantErr: false,
		},
		{
			name:  "RuneArrayReaderError",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromReadCloser(
				NewContentOptions("text/plain"),
				new(ErrorReader),
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs contents
			tt.field.Type.Shape = FieldShapeContent
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectContent(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectContent() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectContent_Content(t *testing.T) {
	type contents struct {
		A Content   `test:"injector"`
		B *Content  `test:"injector"`
		C **Content `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapeContent),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, _ := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(contents{}))

	testContent := Content{
		present: true,
		options: ContentOptions{MimeType: "text/plain"},
		source:  io.NopCloser(bytes.NewBufferString(`"abc"`)),
	}

	tests := []struct {
		name    string
		field   *PortField
		value   Content
		want    contents
		wantErr bool
	}{
		{
			name:  "Content",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				A: testContent,
			},
			wantErr: false,
		},
		{
			name:  "ContentIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				B: &testContent,
			},
			wantErr: false,
		},
		{
			name:  "ContentDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				C: func() **Content {
					v := &testContent
					return &v
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs contents
			tt.field.Type.Shape = FieldShapeContent
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectContent(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectContent() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectContent_Reader(t *testing.T) {
	type contents struct {
		A io.ReadCloser  `test:"injector"`
		B *io.ReadCloser `test:"injector"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupInjector: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapeContent),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(contents{}))
	assert.NoError(t, err)

	testContent := Content{
		present: true,
		options: ContentOptions{MimeType: "text/plain"},
		source:  io.NopCloser(bytes.NewBufferString(`"abc"`)),
	}

	tests := []struct {
		name    string
		field   *PortField
		value   Content
		want    contents
		wantErr bool
	}{
		{
			name:  "Reader",
			field: port.Fields.First(PortFieldHasName("A")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				A: func() io.ReadCloser {
					v, _ := testContent.Reader()
					return v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "ReaderIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte(`"abc"`),
			),
			want: contents{
				B: func() *io.ReadCloser {
					v, _ := testContent.Reader()
					return &v
				}(),
			},
			wantErr: false,
		},
		{
			name:  "ReaderIndirectEmpty",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte{},
			),
			want: contents{
				B: func() *io.ReadCloser {
					v := io.NopCloser(bytes.NewBufferString(""))
					return &v
				}(),
			},
		},
		{
			name:  "ReaderIndirectMissing",
			field: port.Fields.First(PortFieldHasName("B")),
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				nil,
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs contents
			tt.field.Type.Shape = FieldShapeContent
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectContent(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectContent() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}

func TestPortFieldInjector_InjectContent_Error(t *testing.T) {
	type primitives struct {
		A chan struct{} `test:"injector"`
	}

	tests := []struct {
		name    string
		field   *PortField
		value   Content
		want    primitives
		wantErr bool
	}{
		{
			name: "IncorrectShape",
			field: &PortField{
				Name:    "A",
				Indices: []int{0},
				Peer:    "a",
				Group:   FieldGroupInjector,
				Type: PortFieldType{
					Shape:        FieldShapePrimitive,
					Type:         reflect.TypeOf(make(chan struct{})),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(make(chan struct{})),
				},
			},
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte("true"),
			),
			wantErr: true,
		},
		{
			name: "MarshalingError",
			field: &PortField{
				Name:    "A",
				Indices: []int{0},
				Peer:    "a",
				Group:   FieldGroupInjector,
				Type: PortFieldType{
					Shape:        FieldShapeContent,
					Type:         reflect.TypeOf(make(chan struct{})),
					Indirections: 0,
					HandlerType:  reflect.TypeOf(make(chan struct{})),
				},
			},
			value: NewContentFromBytes(
				NewContentOptions("text/plain"),
				[]byte("true"),
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputs primitives
			i := NewPortFieldInjector(tt.field, &inputs)

			if err := i.InjectContent(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("InjectContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(tt.want, inputs) {
					t.Errorf("InjectContent() diff\n%s",
						testhelpers.Diff(tt.want, &inputs))
				}
			}
		})
	}
}
