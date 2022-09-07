// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPortFieldExtractor_ExtractPrimitive_Text(t *testing.T) {
	type primitives struct {
		A string       `test:"extractor"`
		B *string      `test:"extractor"`
		C **string     `test:"extractor"`
		G []byte       `test:"extractor"`
		H *[]byte      `test:"extractor"`
		D **[]byte     `test:"extractor"`
		I []rune       `test:"extractor"`
		J *[]rune      `test:"extractor"`
		E **[]rune     `test:"extractor"`
		K MyTextType   `test:"extractor"`
		L *MyTextType  `test:"extractor"`
		F **MyTextType `test:"extractor"`
		M *string      `test:"extractor,optional"`
		N *string      `test:"extractor,default=abc"`
		O *string      `test:"extractor,const=abc"`
		P *string      `test:"extractor,required"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		field   *PortField
		value   primitives
		want    string
		wantErr bool
	}{
		{
			name:  "String",
			field: port.Fields.First(PortFieldHasName("A")),
			value: primitives{
				A: "abc",
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "StringIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: primitives{
				B: types.NewStringPtr("abc"),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "StringDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: primitives{
				C: func() **string {
					v := types.NewStringPtr("abc")
					return &v
				}(),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "ByteSlice",
			field: port.Fields.First(PortFieldHasName("G")),
			value: primitives{
				G: []byte("abc"),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "ByteSliceIndirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: primitives{
				H: types.NewByteSlicePtr([]byte("abc")),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "ByteSliceDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("D")),
			value: primitives{
				D: func() **[]byte {
					v := types.NewByteSlicePtr([]byte("abc"))
					return &v
				}(),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "RuneSlice",
			field: port.Fields.First(PortFieldHasName("I")),
			value: primitives{
				I: []rune("abc"),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "RuneSliceIndirect",
			field: port.Fields.First(PortFieldHasName("J")),
			value: primitives{
				J: types.NewRuneSlicePtr([]rune("abc")),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "RuneSliceDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: primitives{
				E: func() **[]rune {
					v := types.NewRuneSlicePtr([]rune("abc"))
					return &v
				}(),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "TextMarshaler",
			field: port.Fields.First(PortFieldHasName("K")),
			value: primitives{
				K: "abc",
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "TextMarshalerIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: primitives{
				L: NewMyTextType("abc"),
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name:  "TextMarshalerDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: primitives{
				F: func() **MyTextType {
					v := NewMyTextType("abc")
					return &v
				}(),
			},
			want:    "abc",
			wantErr: false,
		},

		{
			name:    "StringOptional",
			field:   port.Fields.First(PortFieldHasName("M")),
			wantErr: false,
		},
		{
			name:    "StringDefault",
			field:   port.Fields.First(PortFieldHasName("N")),
			want:    "abc",
			wantErr: false,
		},
		{
			name:    "StringConst",
			field:   port.Fields.First(PortFieldHasName("O")),
			want:    "abc",
			wantErr: false,
		},
		{
			name:    "StringRequiredError",
			field:   port.Fields.First(PortFieldHasName("P")),
			value:   primitives{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewPortFieldExtractor(tt.field, &tt.value)

			gotValue, err := i.ExtractPrimitive()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := types.OptionalOf(tt.want)
				if tt.want == "" {
					want = types.OptionalEmpty[string]()
				}
				if !reflect.DeepEqual(want, gotValue) {
					t.Errorf("ExtractPrimitive() diff\n%s",
						testhelpers.Diff(want, gotValue))
				}
			}
		})
	}
}

func TestPortFieldExtractor_ExtractPrimitive_Int(t *testing.T) {
	type primitives struct {
		A int     `test:"extractor"`
		B *int    `test:"extractor"`
		C **int   `test:"extractor"`
		D int8    `test:"extractor"`
		E *int8   `test:"extractor"`
		F **int8  `test:"extractor"`
		G int16   `test:"extractor"`
		H *int16  `test:"extractor"`
		I **int16 `test:"extractor"`
		J int32   `test:"extractor"`
		K *int32  `test:"extractor"`
		L **int32 `test:"extractor"`
		M int64   `test:"extractor"`
		N *int64  `test:"extractor"`
		O **int64 `test:"extractor"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
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
		value   primitives
		want    string
		wantErr bool
	}{
		{
			name:  "Int",
			field: port.Fields.First(PortFieldHasName("A")),
			value: primitives{
				A: int(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "IntIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: primitives{
				B: types.NewIntPtr(9),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "IntDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: primitives{
				C: func() **int {
					v := types.NewIntPtr(10)
					return &v
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Int8",
			field: port.Fields.First(PortFieldHasName("D")),
			value: primitives{
				D: int8(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Int8Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: primitives{
				E: func() *int8 {
					v := int8(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Int8DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: primitives{
				F: func() **int8 {
					v := int8(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Int16",
			field: port.Fields.First(PortFieldHasName("G")),
			value: primitives{
				G: int16(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Int16Indirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: primitives{
				H: func() *int16 {
					v := int16(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Int16DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("I")),
			value: primitives{
				I: func() **int16 {
					v := int16(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Int32",
			field: port.Fields.First(PortFieldHasName("J")),
			value: primitives{
				J: int32(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Int32Indirect",
			field: port.Fields.First(PortFieldHasName("K")),
			value: primitives{
				K: func() *int32 {
					v := int32(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Int32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: primitives{
				L: func() **int32 {
					v := int32(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		}, {
			name:  "Int64",
			field: port.Fields.First(PortFieldHasName("M")),
			value: primitives{
				M: int64(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Int64Indirect",
			field: port.Fields.First(PortFieldHasName("N")),
			value: primitives{
				N: func() *int64 {
					v := int64(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Int64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("O")),
			value: primitives{
				O: func() **int64 {
					v := int64(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewPortFieldExtractor(tt.field, &tt.value)

			gotValue, err := i.ExtractPrimitive()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := types.OptionalOf(tt.want)
				if !reflect.DeepEqual(want, gotValue) {
					t.Errorf("ExtractPrimitive() diff\n%s",
						testhelpers.Diff(want, gotValue))
				}
			}
		})
	}
}

func TestPortFieldExtractor_ExtractPrimitive_Uint(t *testing.T) {
	type primitives struct {
		A uint     `test:"extractor"`
		B *uint    `test:"extractor"`
		C **uint   `test:"extractor"`
		D uint8    `test:"extractor"`
		E *uint8   `test:"extractor"`
		F **uint8  `test:"extractor"`
		G uint16   `test:"extractor"`
		H *uint16  `test:"extractor"`
		I **uint16 `test:"extractor"`
		J uint32   `test:"extractor"`
		K *uint32  `test:"extractor"`
		L **uint32 `test:"extractor"`
		M uint64   `test:"extractor"`
		N *uint64  `test:"extractor"`
		O **uint64 `test:"extractor"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
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
		value   primitives
		want    string
		wantErr bool
	}{
		{
			name:  "Uint",
			field: port.Fields.First(PortFieldHasName("A")),
			value: primitives{
				A: 8,
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "UintIndirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: primitives{
				B: types.NewUintPtr(9),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "UintDoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: primitives{
				C: func() **uint {
					v := types.NewUintPtr(10)
					return &v
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Uint8",
			field: port.Fields.First(PortFieldHasName("D")),
			value: primitives{
				D: uint8(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Uint8Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: primitives{
				E: func() *uint8 {
					v := uint8(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Uint8DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: primitives{
				F: func() **uint8 {
					v := uint8(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Uint16",
			field: port.Fields.First(PortFieldHasName("G")),
			value: primitives{
				G: uint16(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Uint16Indirect",
			field: port.Fields.First(PortFieldHasName("H")),
			value: primitives{
				H: func() *uint16 {
					v := uint16(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Uint16DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("I")),
			value: primitives{
				I: func() **uint16 {
					v := uint16(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Uint32",
			field: port.Fields.First(PortFieldHasName("J")),
			value: primitives{
				J: uint32(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Uint32Indirect",
			field: port.Fields.First(PortFieldHasName("K")),
			value: primitives{
				K: func() *uint32 {
					v := uint32(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Uint32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("L")),
			value: primitives{
				L: func() **uint32 {
					v := uint32(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Uint64",
			field: port.Fields.First(PortFieldHasName("M")),
			value: primitives{
				M: uint64(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Uint64Indirect",
			field: port.Fields.First(PortFieldHasName("N")),
			value: primitives{
				N: func() *uint64 {
					v := uint64(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Uint64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("O")),
			value: primitives{
				O: func() **uint64 {
					v := uint64(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewPortFieldExtractor(tt.field, &tt.value)

			gotValue, err := i.ExtractPrimitive()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := types.OptionalOf(tt.want)
				if !reflect.DeepEqual(want, gotValue) {
					t.Errorf("ExtractPrimitive() diff\n%s",
						testhelpers.Diff(want, gotValue))
				}
			}
		})
	}
}

func TestPortFieldExtractor_ExtractPrimitive_Float(t *testing.T) {
	type primitives struct {
		A float32   `test:"extractor"`
		B *float32  `test:"extractor"`
		C **float32 `test:"extractor"`
		D float64   `test:"extractor"`
		E *float64  `test:"extractor"`
		F **float64 `test:"extractor"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
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
		value   primitives
		want    string
		wantErr bool
	}{
		{
			name:  "Float32",
			field: port.Fields.First(PortFieldHasName("A")),
			value: primitives{
				A: float32(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Float32Indirect",
			field: port.Fields.First(PortFieldHasName("B")),
			value: primitives{
				B: func() *float32 {
					v := float32(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Float32DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("C")),
			value: primitives{
				C: func() **float32 {
					v := float32(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
		{
			name:  "Float64",
			field: port.Fields.First(PortFieldHasName("D")),
			value: primitives{
				D: float64(8),
			},
			want:    "8",
			wantErr: false,
		},
		{
			name:  "Float64Indirect",
			field: port.Fields.First(PortFieldHasName("E")),
			value: primitives{
				E: func() *float64 {
					v := float64(9)
					return &v
				}(),
			},
			want:    "9",
			wantErr: false,
		},
		{
			name:  "Float64DoubleIndirect",
			field: port.Fields.First(PortFieldHasName("F")),
			value: primitives{
				F: func() **float64 {
					v := float64(10)
					w := &v
					return &w
				}(),
			},
			want:    "10",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewPortFieldExtractor(tt.field, &tt.value)

			gotValue, err := i.ExtractPrimitive()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := types.OptionalOf(tt.want)
				if !reflect.DeepEqual(want, gotValue) {
					t.Errorf("ExtractPrimitive() diff\n%s",
						testhelpers.Diff(want, gotValue))
				}
			}
		})
	}
}

func TestPortFieldExtractor_ExtractPrimitive_Bool(t *testing.T) {
	type primitives struct {
		C bool  `test:"extractor"`
		D *bool `test:"extractor"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
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
		value   primitives
		want    string
		wantErr bool
	}{
		{
			name:  "Boolean",
			field: port.Fields.First(PortFieldHasName("C")),
			value: primitives{
				C: true,
			},
			want:    "true",
			wantErr: false,
		},
		{
			name:  "BooleanIndirect",
			field: port.Fields.First(PortFieldHasName("D")),
			value: primitives{
				D: types.NewBoolPtr(true),
			},
			want:    "true",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewPortFieldExtractor(tt.field, &tt.value)

			gotValue, err := i.ExtractPrimitive()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPrimitive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				want := types.OptionalOf(tt.want)
				if !reflect.DeepEqual(want, gotValue) {
					t.Errorf("ExtractPrimitive() diff\n%s",
						testhelpers.Diff(want, gotValue))
				}
			}
		})
	}
}

func TestPortFieldExtractor_ExtractPrimitive_Error(t *testing.T) {
	type primitives struct {
		A chan struct{} `test:"extractor"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupExtractor: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	_, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(primitives{}))
	assert.Error(t, err)
}

func TestPortFieldExtractor_ExtractRawValue(t *testing.T) {
	type fields struct {
		portField    *PortField
		outputs      interface{}
		outputsValue reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := PortFieldExtractor{
				portField:    tt.fields.portField,
				outputs:      tt.fields.outputs,
				outputsValue: tt.fields.outputsValue,
			}
			assert.Equalf(t, tt.want, i.ExtractRawValue(), "ExtractRawValue()")
		})
	}
}

func TestPortFieldExtractor_ExtractValue(t *testing.T) {
	type fields struct {
		portField    *PortField
		outputs      interface{}
		outputsValue reflect.Value
	}
	tests := []struct {
		name    string
		fields  fields
		wantFv  reflect.Value
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := PortFieldExtractor{
				portField:    tt.fields.portField,
				outputs:      tt.fields.outputs,
				outputsValue: tt.fields.outputsValue,
			}
			gotFv, err := i.ExtractValue()
			if !tt.wantErr(t, err, fmt.Sprintf("ExtractValue()")) {
				return
			}
			assert.Equalf(t, tt.wantFv, gotFv, "ExtractValue()")
		})
	}
}
