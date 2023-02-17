// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
		},
		{
			name: "MultipartFilePointerIndirect",
			arg:  reflect.TypeOf(new(*multipart.FileHeader)),
			wantType: PortFieldType{
				Shape:        FieldShapeFile,
				Type:         reflect.TypeOf(new(multipart.FileHeader)),
				Indirections: 1,
				Optional:     true,
				HandlerType:  reflect.TypeOf(new(multipart.FileHeader)),
			},
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
				Optional:     true,
				HandlerType:  reflect.TypeOf(types.Base64Bytes{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionOut)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
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
				Optional:     true,
			},
		},
		{
			name: "Base64Bytes",
			arg:  reflect.TypeOf([]types.Base64Bytes{}),
			wantType: PortFieldType{
				Shape:        FieldShapeFileArray,
				Type:         reflect.TypeOf([]types.Base64Bytes{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]types.Base64Bytes{}),
				Optional:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionOut)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
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
				Optional:     true,
				Items: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapePrimitive,
						Type:        reflect.TypeOf(int(0)),
						HandlerType: reflect.TypeOf(int(0)),
					},
				},
			},
		},
		{
			name: "IntArrayIndirect",
			arg:  reflect.TypeOf(new([]int)),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]int{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]int{}),
				Optional:     true,
				Items: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapePrimitive,
						Type:        reflect.TypeOf(int(0)),
						HandlerType: reflect.TypeOf(int(0)),
					},
				},
			},
		},
		{
			name: "ObjectArray",
			arg:  reflect.TypeOf([]struct{}{}),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]struct{}{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]struct{}{}),
				Optional:     true,
				Items: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapeObject,
						Fields:      []PortFieldElementType{},
						Type:        reflect.TypeOf(struct{}{}),
						HandlerType: reflect.TypeOf(struct{}{}),
					},
				},
			},
		},
		{
			name: "PointerArray",
			arg:  reflect.TypeOf([]*string{}),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]*string{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf([]*string{}),
				Optional:     true,
				Items: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Optional:     true,
						Shape:        FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 1,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "PointerArrayIndirect",
			arg:  reflect.TypeOf(new([]*string)),
			wantType: PortFieldType{
				Shape:        FieldShapeArray,
				Type:         reflect.TypeOf([]*string{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]*string{}),
				Optional:     true,
				Items: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:        FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 1,
						Optional:     true,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionOut)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(tt.wantType, gotType))
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Object(t *testing.T) {
	type embeddedStruct struct {
		AnotherField string `inp:"test"`
	}

	type nestedStruct struct {
		SingleField string         `inp:"test"`
		SecondField embeddedStruct `inp:"test"`
	}

	type recursiveStruct struct {
		Field string           `inp:"test"`
		Next  *recursiveStruct `inp:"test"`
	}

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
				Optional:     true,
				Keys: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapePrimitive,
						Type:        reflect.TypeOf(""),
						HandlerType: reflect.TypeOf(""),
					},
				},
				Values: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapeAny,
						Type:        reflect.TypeOf((*interface{})(nil)).Elem(),
						HandlerType: reflect.TypeOf((*interface{})(nil)).Elem(),
						Optional:    true,
					},
				},
			},
		},
		{
			name: "MapStringIndirect",
			arg:  reflect.TypeOf(new(map[string]interface{})),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(map[string]interface{}{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(map[string]interface{}{}),
				Optional:     true,
				Keys: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapePrimitive,
						Type:        reflect.TypeOf(""),
						HandlerType: reflect.TypeOf(""),
					},
				},
				Values: &PortFieldElementType{
					PortFieldType: &PortFieldType{
						Shape:       FieldShapeAny,
						Type:        reflect.TypeOf((*interface{})(nil)).Elem(),
						HandlerType: reflect.TypeOf((*interface{})(nil)).Elem(),
						Optional:    true,
					},
				},
			},
		},
		{
			name: "Struct",
			arg:  reflect.TypeOf(struct{}{}),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(struct{}{}),
				Indirections: 0,
				HandlerType:  reflect.TypeOf(struct{}{}),
				Fields:       []PortFieldElementType{},
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
				Optional:     true,
				Fields:       []PortFieldElementType{},
			},
		},
		{
			name: "StructDoubleIndirect",
			arg:  reflect.TypeOf(new(*struct{})),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Type:         reflect.TypeOf(struct{}{}),
				Indirections: 2,
				HandlerType:  reflect.TypeOf(struct{}{}),
				Optional:     true,
				Fields:       []PortFieldElementType{},
			},
		},
		{
			name: "StructNested",
			arg:  reflect.TypeOf(new(nestedStruct)),
			wantType: PortFieldType{
				Shape:        FieldShapeObject,
				Indirections: 1,
				Type:         reflect.TypeOf(nestedStruct{}),
				HandlerType:  reflect.TypeOf(nestedStruct{}),
				Optional:     true,
				Fields: []PortFieldElementType{
					{
						Peer:    "singleField",
						Indices: []int{0},
						PortFieldType: &PortFieldType{
							Shape:       FieldShapePrimitive,
							Type:        reflect.TypeOf(""),
							HandlerType: reflect.TypeOf(""),
						},
					},
					{
						Peer:    "secondField",
						Indices: []int{1},
						PortFieldType: &PortFieldType{
							Shape:       FieldShapeObject,
							Type:        reflect.TypeOf(embeddedStruct{}),
							HandlerType: reflect.TypeOf(embeddedStruct{}),
							Fields: []PortFieldElementType{
								{
									Peer:    "anotherField",
									Indices: []int{0},
									PortFieldType: &PortFieldType{
										Shape:       FieldShapePrimitive,
										Type:        reflect.TypeOf(""),
										HandlerType: reflect.TypeOf(""),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StructRecursive",
			arg:  reflect.TypeOf(new(recursiveStruct)),
			wantType: func() PortFieldType {
				root := PortFieldType{
					Shape:        FieldShapeObject,
					Indirections: 1,
					Type:         reflect.TypeOf(recursiveStruct{}),
					HandlerType:  reflect.TypeOf(recursiveStruct{}),
					Optional:     true,
					Fields: []PortFieldElementType{
						{
							Peer:    "field",
							Indices: []int{0},
							PortFieldType: &PortFieldType{
								Shape:       FieldShapePrimitive,
								Type:        reflect.TypeOf(""),
								HandlerType: reflect.TypeOf(""),
							},
						},
						{
							Peer:    "next",
							Indices: []int{1},
						},
					},
				}

				root.Fields[1].PortFieldType = &root

				return root
			}(),
			wantOptional: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionIn)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
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
		name     string
		arg      reflect.Type
		wantType PortFieldType
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
		},
		{
			name: "ContentIndirect",
			arg:  reflect.TypeOf(new(Content)),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         ContentType,
				Indirections: 1,
				HandlerType:  ContentType,
				Optional:     true,
			},
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
		},
		{
			name: "IoReadCloserIndirect",
			arg:  reflect.TypeOf(new(MyReadCloser)),
			wantType: PortFieldType{
				Shape:        FieldShapeContent,
				Type:         reflect.TypeOf(MyReadCloser{}),
				Indirections: 1,
				Optional:     true,
				HandlerType:  IoReadCloserType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionIn)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Primitive(t *testing.T) {
	tests := []struct {
		name     string
		arg      reflect.Type
		wantType PortFieldType
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
		},
		{
			name: "TypesTimeIndirect",
			arg:  reflect.TypeOf(new(types.Time)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Time{}),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
				Optional:     true,
			},
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
		},
		{
			name: "TypesUuidIndirect",
			arg:  reflect.TypeOf(new(types.UUID)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.UUID{}),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
				Optional:     true,
			},
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
		},
		{
			name: "TypesDurationIndirect",
			arg:  reflect.TypeOf(new(types.Duration)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(types.Duration(0)),
				Indirections: 1,
				HandlerType:  TextUnmarshalerType,
				Optional:     true,
			},
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
		},
		{
			name: "ByteSliceIndirect",
			arg:  reflect.TypeOf(new([]byte)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]byte{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]byte{}),
				Optional:     true,
			},
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
		},
		{
			name: "JsonRawMessageIndirect",
			arg:  reflect.TypeOf(new(json.RawMessage)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(json.RawMessage{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]byte{}),
				Optional:     true,
			},
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
		},
		{
			name: "RuneSliceIndirect",
			arg:  reflect.TypeOf(new([]rune)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf([]rune{}),
				Indirections: 1,
				HandlerType:  reflect.TypeOf([]rune{}),
				Optional:     true,
			},
		},
		{
			name: "IntIndirect",
			arg:  reflect.TypeOf(new(int)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int(0)),
				Optional:     true,
			},
		},
		{
			name: "Int8Indirect",
			arg:  reflect.TypeOf(new(int8)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int8(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int8(0)),
				Optional:     true,
			},
		},
		{
			name: "Int16Indirect",
			arg:  reflect.TypeOf(new(int16)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int16(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int16(0)),
				Optional:     true,
			},
		},
		{
			name: "Int32Indirect",
			arg:  reflect.TypeOf(new(int32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int32(0)),
				Optional:     true,
			},
		},
		{
			name: "Int64Indirect",
			arg:  reflect.TypeOf(new(int64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(int64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(int64(0)),
				Optional:     true,
			},
		},
		{
			name: "UintIndirect",
			arg:  reflect.TypeOf(new(uint)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint(0)),
				Optional:     true,
			},
		},
		{
			name: "Uint8Indirect",
			arg:  reflect.TypeOf(new(uint8)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint8(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint8(0)),
				Optional:     true,
			},
		},
		{
			name: "Uint16Indirect",
			arg:  reflect.TypeOf(new(uint16)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint16(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint16(0)),
				Optional:     true,
			},
		},
		{
			name: "Uint32Indirect",
			arg:  reflect.TypeOf(new(uint32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint32(0)),
				Optional:     true,
			},
		},
		{
			name: "Uint64Indirect",
			arg:  reflect.TypeOf(new(uint64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(uint64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(uint64(0)),
				Optional:     true,
			},
		},
		{
			name: "Float32Indirect",
			arg:  reflect.TypeOf(new(float32)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(float32(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(float32(0)),
				Optional:     true,
			},
		},
		{
			name: "Float64Indirect",
			arg:  reflect.TypeOf(new(float64)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(float64(0)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(float64(0)),
				Optional:     true,
			},
		},
		{
			name: "BoolIndirect",
			arg:  reflect.TypeOf(new(bool)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(bool(false)),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(bool(false)),
				Optional:     true,
			},
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
		},
		{
			name: "StringIndirect",
			arg:  reflect.TypeOf(new(string)),
			wantType: PortFieldType{
				Shape:        FieldShapePrimitive,
				Type:         reflect.TypeOf(""),
				Indirections: 1,
				HandlerType:  reflect.TypeOf(""),
				Optional:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionIn)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Error(t *testing.T) {
	tests := []struct {
		name     string
		arg      reflect.Type
		wantType PortFieldType
		wantErr  bool
	}{
		{
			name:    "Channel",
			arg:     reflect.TypeOf(make(chan struct{})),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionIn)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				if !reflect.DeepEqual(gotType, &tt.wantType) {
					t.Errorf("ReflectPortFieldType() diff\n%s",
						testhelpers.Diff(&tt.wantType, gotType))
				}
			}
		})
	}
}

func TestDefaultPortFieldTypeReflector_ReflectPortFieldType_Any(t *testing.T) {
	var anything any
	var anythingType = reflect.TypeOf(&anything).Elem()

	tests := []struct {
		name         string
		arg          reflect.Type
		wantType     PortFieldType
		wantOptional bool
	}{
		{
			name: "Interface",
			arg:  anythingType,
			wantType: PortFieldType{
				Shape:        FieldShapeAny,
				Type:         anythingType,
				Indirections: 0,
				HandlerType:  anythingType,
				Optional:     true,
			},
		},
		{
			name: "IndirectInterface",
			arg:  reflect.TypeOf(&anything),
			wantType: PortFieldType{
				Shape:        FieldShapeAny,
				Type:         anythingType,
				Indirections: 1,
				HandlerType:  anythingType,
				Optional:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewDefaultPortFieldTypeReflector(PortDirectionIn)
			gotType, gotErr := r.ReflectPortFieldType(tt.arg)
			assert.Nil(t, gotErr)
			if !reflect.DeepEqual(gotType, &tt.wantType) {
				t.Errorf("ReflectPortFieldType() diff\n%s",
					testhelpers.Diff(&tt.wantType, gotType))
			}
		})
	}
}
