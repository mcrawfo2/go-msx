package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
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
	return NewContentFromBytes(
		ContentOptions{
			MimeType: t.ContentType,
			Encoding: t.ContentEncoding,
		},
		[]byte(t.Values[pf.Peer])), nil
}

func TestInputsPopulator_PopulateInputs_Primitives(t *testing.T) {
	type inputs struct {
		A *string    `test:"populator"`
		B MyTextType `test:"populator"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive, FieldShapeContent),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewInputsPopulator(
				port,
				TestMapDecoder{Values: tt.values},
			)

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
		A Content `test:"populator"`
	}

	pr := PortReflector{
		FieldGroups: map[string]FieldGroup{
			FieldGroupPopulator: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(FieldShapePrimitive, FieldShapeContent),
			},
		},
		FieldTypeReflector: DefaultPortFieldTypeReflector{},
	}

	port, err := pr.ReflectPortStruct(PortTypeTest, reflect.TypeOf(inputs{}))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to reflect port struct", err.Error())
	}

	tests := []struct {
		name    string
		values  map[string]string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewInputsPopulator(
				port,
				TestMapDecoder{
					Values:          tt.values,
					ContentType:     tt.options.MimeType,
					ContentEncoding: tt.options.Encoding,
				},
			)

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
