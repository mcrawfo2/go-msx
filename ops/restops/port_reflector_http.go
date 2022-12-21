// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/iancoleman/strcase"
	"net/textproto"
	"reflect"
)

const (
	FieldGroupHttpMethod = "method"
	FieldGroupHttpHeader = "header"
	FieldGroupHttpCookie = "cookie"
	FieldGroupHttpPath   = "path"
	FieldGroupHttpQuery  = "query"
	FieldGroupHttpForm   = "form"
	FieldGroupHttpBody   = "body"
	FieldGroupHttpPaging = "paging"
	FieldGroupHttpCode   = "code"
)

type PortReflector struct{}

func (r PortReflector) ReflectInputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		Direction: ops.PortDirectionIn,
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupHttpMethod: {
				Cardinality: types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
				),
			},
			FieldGroupHttpHeader: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
					ops.FieldShapeArray,
					ops.FieldShapeObject,
				),
			},
			FieldGroupHttpCookie: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
				),
			},
			FieldGroupHttpPath: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
				),
			},
			FieldGroupHttpQuery: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
					ops.FieldShapeArray,
					ops.FieldShapeObject,
				),
			},
			FieldGroupHttpForm: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
					ops.FieldShapeArray,
					ops.FieldShapeObject,
					ops.FieldShapeFile,
					ops.FieldShapeFileArray,
				),
			},
			FieldGroupHttpBody: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapeContent,
				),
			},
		},
		FieldPostProcessor: r.postProcessField,
	}

	return reflector.ReflectPortStruct(PortTypeRequest, st)
}

func (r PortReflector) ReflectOutputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		Direction: ops.PortDirectionOut,
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupHttpCode: {
				Cardinality: types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
				),
			},
			FieldGroupHttpHeader: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapePrimitive,
					ops.FieldShapeArray,
					ops.FieldShapeObject,
				),
			},
			FieldGroupHttpBody: {
				Cardinality: types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapeContent,
				),
			},
			FieldGroupHttpPaging: {
				Cardinality: types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(
					ops.FieldShapeObject,
				),
			},
		},
		FieldPostProcessor: r.postProcessField,
	}

	return reflector.ReflectPortStruct(PortTypeResponse, st)
}

func (r PortReflector) postProcessField(pf *ops.PortField, sf reflect.StructField) {
	if PortReflectorPostProcessField != nil {
		PortReflectorPostProcessField(pf, sf)
	}

	if pf.Group == FieldGroupHttpPath {
		pf.Optional = false
		pf.Options["required"] = "true"
		delete(pf.Options, "optional")
	}

	if pf.Group == FieldGroupHttpHeader {
		pf.Peer = textproto.CanonicalMIMEHeaderKey(strcase.ToKebab(pf.Name))
	}

	// Input and Output message body are processed as Content (encoded/marshalled by Content-Type/Content-Encoding)
	if pf.Group == FieldGroupHttpBody && pf.Type.Shape != ops.FieldShapeUnknown {
		pf.Type.Shape = ops.FieldShapeContent
	}
}

var PortReflectorPostProcessField ops.PortFieldPostProcessorFunc
