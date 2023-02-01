// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"reflect"
)

const (
	FieldShapePrimitive = "primitive" // Input, Output
	FieldShapeArray     = "array"     // Input, Output
	FieldShapeObject    = "object"    // Input, Output
	FieldShapeFile      = "file"      // Input
	FieldShapeFileArray = "fileArray" // Input
	FieldShapeContent   = "content"   // Input, Output
	FieldShapeUnknown   = "unknown"   // Ignored
	FieldShapeAny       = "any"       // Input, Output
)

var ErrIncorrectShape = errors.New("Port field has mismatched shape")

var TextUnmarshalerInstance types.TextUnmarshaler
var TextUnmarshalerType = reflect.TypeOf(&TextUnmarshalerInstance).Elem()
var TextMarshalerInstance types.TextMarshaler
var TextMarshalerType = reflect.TypeOf(&TextMarshalerInstance).Elem()
var MultipartFileHeaderInstance multipart.FileHeader
var MultipartFileHeaderType = reflect.TypeOf(MultipartFileHeaderInstance)
var MultipartFileHeaderPtrType = reflect.PtrTo(MultipartFileHeaderType)
var ByteSliceType = reflect.TypeOf([]byte{})
var RuneSliceType = reflect.TypeOf([]rune{})
var Base64BytesType = reflect.TypeOf(types.Base64Bytes{})
var IoReadCloserInstance io.ReadCloser
var IoReadCloserType = reflect.TypeOf(&IoReadCloserInstance).Elem()
var ContentInstance Content
var ContentType = reflect.TypeOf(&ContentInstance).Elem()
var JsonRawMessageInstance json.RawMessage
var JsonRawMessageType = reflect.TypeOf(&JsonRawMessageInstance).Elem()
var OptionalOfStringInstance types.Optional[string]
var OptionalOfStringType = reflect.TypeOf(&OptionalOfStringInstance).Elem()
var AnyInstance any
var AnyType = reflect.TypeOf(&AnyInstance).Elem()

type PortFieldTypeReflector interface {
	ReflectPortFieldType(reflect.Type) (*PortFieldType, error)
}

type PortFieldTypeReflectorFunc func(reflect.Type) (*PortFieldType, error)

func (f PortFieldTypeReflectorFunc) ReflectPortFieldType(t reflect.Type) (*PortFieldType, error) {
	return f(t)
}

type DefaultPortFieldTypeReflector struct {
	Direction         PortDirection
	OnReflectIndirect PortFieldTypeReflector
	OnReflectDirect   PortFieldTypeReflector
	Placeholders      map[reflect.Type]*PortFieldType
}

func NewDefaultPortFieldTypeReflector(direction PortDirection) DefaultPortFieldTypeReflector {
	return DefaultPortFieldTypeReflector{
		Direction:    direction,
		Placeholders: make(map[reflect.Type]*PortFieldType),
	}
}

func (r DefaultPortFieldTypeReflector) ReflectPortFieldType(t reflect.Type) (portFieldType *PortFieldType, err error) {
	if pft, ok := r.Placeholders[t]; ok {
		return pft, nil
	}

	pft := new(PortFieldType)
	r.Placeholders[t] = pft

	if t.Kind() == reflect.Ptr {
		portFieldType, err = r.reflectIndirect(t)
	} else {
		portFieldType, err = r.reflectDirect(t)
	}

	if err == nil {
		*pft = *portFieldType
		r.Placeholders[t] = pft
	}

	return
}

// reflectIndirect identifies types that are required to be pointers.
func (r DefaultPortFieldTypeReflector) reflectIndirect(t reflect.Type) (portFieldType *PortFieldType, err error) {
	if r.OnReflectIndirect != nil {
		portFieldType, err = r.OnReflectIndirect.ReflectPortFieldType(t)
		if err != nil || portFieldType.Shape != "" {
			return
		}
	}

	switch t {
	case MultipartFileHeaderPtrType:
		portFieldType = PortFieldTypeFromType(t, FieldShapeFile)
		return
	}

	t = t.Elem()
	portFieldType, err = r.ReflectPortFieldType(t)
	if err == nil {
		portFieldType = portFieldType.IncIndirections().SetOptional(true)
	}

	return
}

// reflectDirect identifies types that are not pointers.
func (r DefaultPortFieldTypeReflector) reflectDirect(t reflect.Type) (portFieldType *PortFieldType, err error) {
	if r.OnReflectDirect != nil {
		portFieldType, err = r.OnReflectDirect.ReflectPortFieldType(t)
		if err != nil || portFieldType.Shape != "" {
			return
		}
	}

	// Concrete Types
	switch t {
	case Base64BytesType:
		return PortFieldTypeFromType(t, FieldShapeFile), nil
	case ByteSliceType, JsonRawMessageType:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		portFieldType.WithHandlerType(ByteSliceType)
		return
	case RuneSliceType:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	case OptionalOfStringType:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive).SetOptional(true)
		return
	case ContentType:
		portFieldType = PortFieldTypeFromType(t, FieldShapeContent)
		return
	case IoReadCloserType:
		portFieldType = PortFieldTypeFromType(t, FieldShapeContent)
		return
	}

	// Interfaces
	pt := reflect.PtrTo(t)
	switch {
	case r.Direction == PortDirectionIn && pt.Implements(TextUnmarshalerType):
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		portFieldType.WithHandlerType(TextUnmarshalerType)
		return
	case r.Direction == PortDirectionOut && pt.Implements(TextMarshalerType):
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		portFieldType.WithHandlerType(TextMarshalerType)
		return
	case pt.Implements(IoReadCloserType):
		portFieldType = PortFieldTypeFromType(t, FieldShapeContent)
		portFieldType.WithHandlerType(IoReadCloserType)
		return
	}

	// Kinds
	switch t.Kind() {
	case reflect.Slice:
		te := t.Elem()

		if te == MultipartFileHeaderPtrType || te == Base64BytesType {
			portFieldType = PortFieldTypeFromType(t, FieldShapeFileArray).SetOptional(true)
			return
		}

		portFieldType = PortFieldTypeFromType(t, FieldShapeArray)

		elemPortFieldType, err := r.ReflectPortFieldType(te)
		if err != nil {
			return nil, err
		}

		portFieldType.Items = &PortFieldElementType{
			PortFieldType: elemPortFieldType,
		}

		portFieldType = portFieldType.SetOptional(true)
		return portFieldType, nil

	case reflect.Map:
		portFieldType = PortFieldTypeFromType(t, FieldShapeObject)

		keyPortFieldType, err := r.ReflectPortFieldType(t.Key())
		if err != nil {
			return nil, err
		}
		portFieldType.Keys = &PortFieldElementType{
			PortFieldType: keyPortFieldType,
		}

		valuePortFieldType, err := r.ReflectPortFieldType(t.Elem())
		if err != nil {
			return nil, err
		}
		portFieldType.Values = &PortFieldElementType{
			PortFieldType: valuePortFieldType,
		}
		portFieldType = portFieldType.SetOptional(true)
		return portFieldType, nil

	case reflect.Struct:
		portFieldType = PortFieldTypeFromType(t, FieldShapeObject)
		visitor := newPortFieldTypeReflectorFieldVisitor(r)
		err = WalkStruct(t, visitor)
		if err != nil {
			return nil, err
		}

		portFieldType.Fields = visitor.Fields
		return
	case reflect.Float64, reflect.Float32:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	case reflect.String:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	case reflect.Bool:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		return
	}

	if t == AnyType {
		portFieldType = PortFieldTypeFromType(t, FieldShapeAny).SetOptional(true)
		return
	}

	err = errors.Wrapf(ErrInvalidShape, "%+v", t)
	return
}

type PortFieldTypeReflectorFieldVisitor struct {
	Reflector PortFieldTypeReflector
	Fields    []PortFieldElementType
	*FieldVisitor
}

func (v *PortFieldTypeReflectorFieldVisitor) VisitField(f reflect.StructField) error {
	pft, err := v.Reflector.ReflectPortFieldType(f.Type)
	if err != nil {
		return err
	}

	pfet := PortFieldElementType{
		Peer:          strcase.ToLowerCamel(f.Name),
		Indices:       f.Index,
		PortFieldType: pft,
	}

	v.Fields = append(v.Fields, pfet)

	v.incrementIndex()
	return nil
}

func newPortFieldTypeReflectorFieldVisitor(r PortFieldTypeReflector) *PortFieldTypeReflectorFieldVisitor {
	return &PortFieldTypeReflectorFieldVisitor{
		Reflector:    r,
		Fields:       []PortFieldElementType{},
		FieldVisitor: newFieldVisitor(),
	}
}
