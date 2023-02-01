// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"reflect"
	"strings"
	"testing"
)

type TestMapDecoder struct {
	Values          map[string]string
	ContentType     string
	ContentEncoding []string
}

func (t TestMapDecoder) DecodePrimitive(pf *PortField) (result types.Optional[string], err error) {
	value, ok := t.Values[pf.Peer]
	if !ok {
		return
	}

	result = types.OptionalOf[string](value)
	return
}

func (t TestMapDecoder) DecodeContent(pf *PortField) (content Content, err error) {
	if pf.Optional {
		return
	}

	return NewContentFromBytes(
		ContentOptions{
			MimeType: t.ContentType,
			Encoding: t.ContentEncoding,
		},
		[]byte(t.Values[pf.Peer])), nil
}

func (t TestMapDecoder) DecodeArray(pf *PortField) (result []string, err error) {
	value, ok := t.Values[pf.Peer]
	if !ok {
		return
	}

	result = strings.Split(value, ",")
	return
}

func (t TestMapDecoder) DecodeObject(pf *PortField) (result types.Pojo, err error) {
	result = make(types.Pojo)
	for k, v := range t.Values {
		result[k] = v
	}
	return
}

func (t TestMapDecoder) DecodeFile(pf *PortField) (result *multipart.FileHeader, err error) {
	panic("implement me")
}

func (t TestMapDecoder) DecodeFileArray(pf *PortField) (result []*multipart.FileHeader, err error) {
	panic("implement me")
}

func (t TestMapDecoder) DecodeAny(pf *PortField) (result types.Optional[any], err error) {
	panic("implement me")
}

func TestInputsPopulator_PopulateInputs_Primitives(t *testing.T) {
	type inputs struct {
		A *string    `test:"populator"`
		B MyTextType `test:"populator"`
	}

	pr := PortReflector{
		Direction: PortDirectionIn,
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive, FieldShapeContent),
			},
		},
		FieldTypeReflector: NewDefaultPortFieldTypeReflector(PortDirectionIn),
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
		decoder func() InputDecoder
		want    interface{}
		wantErr bool
	}{
		{
			name: "Primitive",
			values: map[string]string{
				"a": "123",
				"b": "456",
			},
			want: &inputs{
				A: types.NewStringPtr("123"),
				B: MyTextType("456"),
			},
		},
		{
			name: "ValidationFailure",
			values: map[string]string{
				"a": "123",
				"b": "error",
			},
			wantErr: true,
		},
		{
			name: "DecoderFailure",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodePrimitive", mock.Anything).
					Return(types.OptionalEmpty[string](), errors.New("decoder error"))
				return decoder
			},
			wantErr: true,
		},
		{
			name: "NonOptionalFailure",
			values: map[string]string{
				"a": "123",
			},
			wantErr: true,
		},
		{
			name: "OptionalOk",
			values: map[string]string{
				"b": "456",
			},
			want: &inputs{
				B: MyTextType("456"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoder InputDecoder = TestMapDecoder{Values: tt.values}
			if tt.decoder != nil {
				decoder = tt.decoder()
			}

			p := NewInputsPopulator(port, decoder)

			got, err := p.PopulateInputs()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("PopulateInputs() diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestInputsPopulator_PopulateInputs_Arrays(t *testing.T) {
	type inputs struct {
		A []MyTextType `test:"populator"`
		B []int        `test:"populator"`
		C []*string    `test:"populator"`
	}

	pr := PortReflector{
		Direction: PortDirectionIn,
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapeArray),
			},
		},
		FieldTypeReflector: NewDefaultPortFieldTypeReflector(PortDirectionIn),
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
		decoder func() InputDecoder
		want    interface{}
		wantErr bool
	}{
		{
			name: "Array",
			values: map[string]string{
				"a": "123,456",
				"b": "789,123",
				"c": "456,789",
			},
			want: &inputs{
				A: []MyTextType{
					"123",
					"456",
				},
				B: []int{
					789,
					123,
				},
				C: []*string{
					types.NewStringPtr("456"),
					types.NewStringPtr("789"),
				},
			},
		},
		{
			name: "ValidationFailure",
			values: map[string]string{
				"a": "123",
				"b": "error",
			},
			wantErr: true,
		},
		{
			name: "DecoderFailure",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodeArray", mock.Anything).
					Return([]string{}, errors.New("decoder error"))
				return decoder
			},
			wantErr: true,
		},
		{
			name: "Optional",
			values: map[string]string{
				"a": "123,456",
				"c": "456,789",
			},
			want: &inputs{
				A: []MyTextType{
					"123",
					"456",
				},
				B: nil,
				C: []*string{
					types.NewStringPtr("456"),
					types.NewStringPtr("789"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoder InputDecoder = TestMapDecoder{Values: tt.values}
			if tt.decoder != nil {
				decoder = tt.decoder()
			}

			p := NewInputsPopulator(port, decoder)

			got, err := p.PopulateInputs()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("PopulateInputs() diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestInputsPopulator_PopulateInputs_Objects(t *testing.T) {
	type subInputs struct {
		A int
		B string
		C *bool
	}

	type inputs struct {
		A *subInputs `test:"populator"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapeObject),
			},
		},
		FieldTypeReflector: NewDefaultPortFieldTypeReflector(PortDirectionIn),
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
		decoder func() InputDecoder
		want    interface{}
		wantErr bool
	}{
		{
			name: "Array",
			values: map[string]string{
				"a": "123",
				"b": "456",
				"c": "true",
			},
			want: &inputs{A: &subInputs{
				A: 123,
				B: "456",
				C: types.NewBoolPtr(true),
			}},
		},
		{
			name: "ValidationFailure",
			values: map[string]string{
				"a": "123",
				"b": "456",
				"c": "no",
			},
			wantErr: true,
		},
		{
			name: "DecoderFailure",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodeObject", mock.Anything).
					Return(nil, errors.New("decoder error"))
				return decoder
			},
			wantErr: true,
		},
		{
			name: "OptionalField",
			values: map[string]string{
				"a": "123",
				"b": "456",
			},
			want: &inputs{A: &subInputs{
				A: 123,
				B: "456",
			}},
		},
		{
			name: "OptionalObject",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodeObject", mock.Anything).
					Return(nil, nil)
				return decoder
			},
			want: &inputs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoder InputDecoder = TestMapDecoder{Values: tt.values}
			if tt.decoder != nil {
				decoder = tt.decoder()
			}

			p := NewInputsPopulator(port, decoder)

			got, err := p.PopulateInputs()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("PopulateInputs() diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestInputsPopulator_PopulateInputs_Content(t *testing.T) {
	type inputs struct {
		A Content  `test:"populator"`
		B *Content `test:"populator"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive, FieldShapeContent),
			},
		},
		FieldTypeReflector: NewDefaultPortFieldTypeReflector(PortDirectionIn),
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
		decoder func() InputDecoder
		options ContentOptions
		want    interface{}
		wantErr bool
	}{
		{
			name: "Content",
			values: map[string]string{
				"a": "123",
			},
			options: ContentOptions{
				MimeType: "text/plain",
				Encoding: nil,
			},
			want: &inputs{
				A: NewContentFromBytes(
					ContentOptions{
						MimeType: "text/plain",
						Encoding: nil,
					},
					[]byte("123")),
			},
		},
		{
			name: "DecoderFailure",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodeContent", mock.Anything).
					Return(Content{}, errors.New("decoder error"))
				return decoder
			},
			wantErr: true,
		},
		{
			name: "OptionalContent",
			decoder: func() InputDecoder {
				decoder := new(MockInputDecoder)
				decoder.
					On("DecodeContent", mock.Anything).
					Return(Content{}, nil)
				return decoder
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoder InputDecoder = TestMapDecoder{
				Values:          tt.values,
				ContentType:     tt.options.MimeType,
				ContentEncoding: tt.options.Encoding,
			}

			if tt.decoder != nil {
				decoder = tt.decoder()
			}

			p := NewInputsPopulator(port, decoder)

			got, err := p.PopulateInputs()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("PopulateInputs() diff\n%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}
