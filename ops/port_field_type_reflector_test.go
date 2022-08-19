// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_File(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "MultipartFilePointer",
			arg:  reflect.TypeOf(new(multipart.FileHeader)),
			wantType: PortFieldType{
				Shape:        FieldShapeFile,
				Type:         reflect.TypeOf(new(multipart.FileHeader)),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(new(multipart.FileHeader)),
			},
			wantOptional: false,
		},
		{
			name: "MultipartFilePointerIndirect",
			arg:  reflect.TypeOf(new(*multipart.FileHeader)),
			wantType: PortFieldType{
				Shape:        FieldShapeFile,
				Type:         reflect.TypeOf(new(multipart.FileHeader)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(new(multipart.FileHeader)),
			},
			wantOptional: true,
		},
		{
			name: "Base64Bytes",
			arg:  reflect.TypeOf(types.Base64Bytes{}),
			wantType: PortFieldType{
				Shape:        FieldShapeFile,
				Type:         reflect.TypeOf(types.Base64Bytes{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(types.Base64Bytes{}),
			},
		},
		{
			name: "Base64BytesIndirect",
			arg:  reflect.TypeOf(new(types.Base64Bytes)),
			wantType: PortFieldType{
				Shape:        FieldShapeFile,
				Type:         reflect.TypeOf(types.Base64Bytes{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(types.Base64Bytes{}),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_FileArray(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "MultipartFilePointer",
			arg:  reflect.TypeOf([]*multipart.FileHeader{}),
			wantType: PortFieldType{
				Shape:        FieldShapeFileArray,
				Type:         reflect.TypeOf([]*multipart.FileHeader{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]*multipart.FileHeader{}),
			},
			wantOptional: true,
		},
		{
			name: "Base64Bytes",
			arg:  reflect.TypeOf([]types.Base64Bytes{}),
			wantType: PortFieldType{
				Shape:        FieldShapeFileArray,
				Type:         reflect.TypeOf([]types.Base64Bytes{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]types.Base64Bytes{}),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Array(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "IntArray",
			arg:  reflect.TypeOf([]int{}),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]int{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]int{}),
			},
			wantOptional: true,
		},
		{
			name: "IntArrayIndirect",
			arg:  reflect.TypeOf(new([]int)),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]int{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]int{}),
			},
			wantOptional: true,
		},
		{
			name: "ObjectArray",
			arg:  reflect.TypeOf([]struct{}{}),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]struct{}{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]struct{}{}),
			},
			wantOptional: true,
		},
		{
			name: "PointerArray",
			arg:  reflect.TypeOf([]*string{}),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]*string{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]*string{}),
			},
			wantOptional: true,
		},
		{
			name: "PointerArrayIndirect",
			arg:  reflect.TypeOf(new([]*string)),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]*string{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]*string{}),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Object(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "MapString",
			arg:  reflect.TypeOf(map[string]interface{}{}),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(map[string]interface{}{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(map[string]interface{}{}),
			},
			wantOptional: true,
		},
		{
			name: "MapStringIndirect",
			arg:  reflect.TypeOf(new(map[string]interface{})),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(map[string]interface{}{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(map[string]interface{}{}),
			},
			wantOptional: true,
		},
		{
			name: "Struct",
			arg:  reflect.TypeOf(struct{}{}),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(struct{}{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(struct{}{}),
			},
		},
		{
			name: "StructIndirect",
			arg:  reflect.TypeOf(new(struct{})),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(struct{}{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(struct{}{}),
			},
			wantOptional: true,
		},
		{
			name: "StructDoubleIndirect",
			arg:  reflect.TypeOf(new(*struct{})),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(struct{}{}),
				Indirections: 2,
				HandlerType:  reflect.TypeOf(struct{}{}),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

type MyReadCloser struct{}

func (m *MyReadCloser) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MyReadCloser) Close() error {
	//TODO implement me
	panic("implement me")
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Content(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "Content",
			arg:  reflect.TypeOf(Content{}),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         ContentType,
				Indirections: 0,
				HandlerType:  ContentType,
			},
			wantOptional: false,
		},
		{
			name: "ContentIndirect",
			arg:  reflect.TypeOf(new(Content)),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         ContentType,
				Indirections: 1,
				HandlerType:  ContentType,
			},
			wantOptional: true,
		},
		{
			name: "IoReadCloser",
			arg:  reflect.TypeOf(MyReadCloser{}),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         reflect.TypeOf(MyReadCloser{}),
				Indirections: 0,
				HandlerType:  IoReadCloserType,
			},
			wantOptional: false,
		},
		{
			name: "IoReadCloserIndirect",
			arg:  reflect.TypeOf(new(MyReadCloser)),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         reflect.TypeOf(MyReadCloser{}),
				Indirections: 1,
				HandlerType:  IoReadCloserType,
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Primitive(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "TypesTime",
			arg:  reflect.TypeOf(types.Time{}),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Time{}),
				Indirections: 0,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: false,
		},
		{
			name: "TypesTimeIndirect",
			arg:  reflect.TypeOf(new(types.Time)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Time{}),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: true,
		},
		{
			name: "TypesUuid",
			arg:  reflect.TypeOf(types.UUID{}),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.UUID{}),
				Indirections: 0,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: false,
		},
		{
			name: "TypesUuidIndirect",
			arg:  reflect.TypeOf(new(types.UUID)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.UUID{}),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: true,
		},
		{
			name: "TypesDuration",
			arg:  reflect.TypeOf(types.Duration(0)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Duration(0)),
				Indirections: 0,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: false,
		},
		{
			name: "TypesDurationIndirect",
			arg:  reflect.TypeOf(new(types.Duration)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Duration(0)),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
			},
			wantOptional: true,
		},
		{
			name: "ByteSlice",
			arg:  reflect.TypeOf([]byte{}),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]byte{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]byte{}),
			},
			wantOptional: false,
		},
		{
			name: "ByteSliceIndirect",
			arg:  reflect.TypeOf(new([]byte)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]byte{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]byte{}),
			},
			wantOptional: true,
		},
		{
			name: "JsonRawMessage",
			arg:  reflect.TypeOf(json.RawMessage{}),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(json.RawMessage{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]byte{}),
			},
			wantOptional: false,
		},
		{
			name: "JsonRawMessageIndirect",
			arg:  reflect.TypeOf(new(json.RawMessage)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(json.RawMessage{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]byte{}),
			},
			wantOptional: true,
		},
		{
			name: "RuneSlice",
			arg:  reflect.TypeOf([]rune{}),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]rune{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]rune{}),
			},
			wantOptional: false,
		},
		{
			name: "RuneSliceIndirect",
			arg:  reflect.TypeOf(new([]rune)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]rune{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]rune{}),
			},
			wantOptional: true,
		},
		{
			name: "IntIndirect",
			arg:  reflect.TypeOf(new(int)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int(0)),
			},
			wantOptional: true,
		},
		{
			name: "Int8Indirect",
			arg:  reflect.TypeOf(new(int8)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int8(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int8(0)),
			},
			wantOptional: true,
		},
		{
			name: "Int16Indirect",
			arg:  reflect.TypeOf(new(int16)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int16(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int16(0)),
			},
			wantOptional: true,
		},
		{
			name: "Int32Indirect",
			arg:  reflect.TypeOf(new(int32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int32(0)),
			},
			wantOptional: true,
		},
		{
			name: "Int64Indirect",
			arg:  reflect.TypeOf(new(int64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int64(0)),
			},
			wantOptional: true,
		},
		{
			name: "UintIndirect",
			arg:  reflect.TypeOf(new(uint)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint(0)),
			},
			wantOptional: true,
		},
		{
			name: "Uint8Indirect",
			arg:  reflect.TypeOf(new(uint8)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint8(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint8(0)),
			},
			wantOptional: true,
		},
		{
			name: "Uint16Indirect",
			arg:  reflect.TypeOf(new(uint16)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint16(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint16(0)),
			},
			wantOptional: true,
		},
		{
			name: "Uint32Indirect",
			arg:  reflect.TypeOf(new(uint32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint32(0)),
			},
			wantOptional: true,
		},
		{
			name: "Uint64Indirect",
			arg:  reflect.TypeOf(new(uint64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint64(0)),
			},
			wantOptional: true,
		},
		{
			name: "Float32Indirect",
			arg:  reflect.TypeOf(new(float32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(float32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(float32(0)),
			},
			wantOptional: true,
		},
		{
			name: "Float64Indirect",
			arg:  reflect.TypeOf(new(float64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(float64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(float64(0)),
			},
			wantOptional: true,
		},
		{
			name: "BoolIndirect",
			arg:  reflect.TypeOf(new(bool)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(bool(false)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(bool(false)),
			},
			wantOptional: true,
		},
		{
			name: "String",
			arg:  reflect.TypeOf(""),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(""),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(""),
			},
			wantOptional: false,
		},
		{
			name: "StringIndirect",
			arg:  reflect.TypeOf(new(string)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(""),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(""),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Unknown(t *testing.T) {
	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "Channel",
			arg:  reflect.TypeOf(make(chan struct{})),
			wantType: PortFieldType{
				Shape:        FieldShapeUnknown,
				Type:         reflect.TypeOf(make(chan struct{})),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(make(chan struct{})),
			},
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultPortFieldTypeReflector{}
			gotType, gotOptional := r.ReflectPortFieldType(tt.arg)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
			if gotOptional != tt.wantOptional {
				t.Errorf("ReflectPortFieldType() gotOptional = %v, want %v", gotOptional, tt.wantOptional)
			}
		})
	}
}
