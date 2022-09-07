// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
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
)

var ErrIncorrectShape = errors.New("Port field has mismatched shape")

var TextUnmarshalerInstance types.TextUnmarshaler
var TextUnmarshalerType = reflect.TypeOf(&TextUnmarshalerInstance).Elem()
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

type PortFieldTypeReflector interface {
	ReflectPortFieldType(reflect.Type) (PortFieldType, bool)
}

type DefaultPortFieldTypeReflector struct{}

func (r DefaultPortFieldTypeReflector) ReflectPortFieldType(t reflect.Type) (PortFieldType, bool) {
	if t.Kind() == reflect.Ptr {
		return r.reflectIndirect(t)
	}

	return r.reflectDirect(t)
}

// reflectIndirect identifies types that are required to be pointers.
func (r DefaultPortFieldTypeReflector) reflectIndirect(t reflect.Type) (portFieldType PortFieldType, optional bool) {
	switch t {
	case MultipartFileHeaderPtrType:
		portFieldType = PortFieldTypeFromType(t, FieldShapeFile)
		return
	}

	t = t.Elem()
	portFieldType, _ = r.ReflectPortFieldType(t)
	portFieldType.IncIndirections()
	optional = true
	return
}

// reflectDirect identifies types that are not pointers.
func (r DefaultPortFieldTypeReflector) reflectDirect(t reflect.Type) (portFieldType PortFieldType, optional bool) {
	// Concrete Types
	switch t {
	case Base64BytesType:
		return PortFieldTypeFromType(t, FieldShapeFile), false
	case ByteSliceType, JsonRawMessageType:
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		portFieldType.WithHandlerType(ByteSliceType)
		return portFieldType, false
	case RuneSliceType:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case OptionalOfStringType:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case ContentType:
		return PortFieldTypeFromType(t, FieldShapeContent), false
	case IoReadCloserType:
		return PortFieldTypeFromType(t, FieldShapeContent), false
	}

	// Interfaces
	pt := reflect.PtrTo(t)
	switch {
	case pt.Implements(TextUnmarshalerType):
		portFieldType = PortFieldTypeFromType(t, FieldShapePrimitive)
		portFieldType.WithHandlerType(TextUnmarshalerType)
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
			return PortFieldTypeFromType(t, FieldShapeFileArray), true
		}
		return PortFieldTypeFromType(t, FieldShapeArray), true
	case reflect.Map:
		return PortFieldTypeFromType(t, FieldShapeObject), true
	case reflect.Struct:
		return PortFieldTypeFromType(t, FieldShapeObject), false
	case reflect.Float64, reflect.Float32:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case reflect.String:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	case reflect.Bool:
		return PortFieldTypeFromType(t, FieldShapePrimitive), false
	}

	logger.Warnf("Cannot determine field shape '%+v'", t)
	return PortFieldTypeFromType(t, FieldShapeUnknown), true
}
