// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
)

const (
	FieldGroupStreamChannel   = "channel"
	FieldGroupStreamMessageId = "messageId"
	FieldGroupStreamHeader    = "header"
	FieldGroupStreamBody      = "body"
)

const (
	PortTypeInput  = "in"
	PortTypeOutput = "out"
)

type PortReflector struct{}

func (r PortReflector) ReflectInputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupStreamChannel: {
				Cardinality:   types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamMessageId: {
				Cardinality:   types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamHeader: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamBody: {
				Cardinality:   types.CardinalityOneToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapeContent),
			},
		},
		FieldPostProcessor: r.PostProcessField,
	}

	return reflector.ReflectPortStruct(PortTypeInput, st)
}

func (r PortReflector) ReflectOutputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupStreamChannel: {
				Cardinality:   types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamMessageId: {
				Cardinality:   types.CardinalityZeroToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamHeader: {
				Cardinality:   types.CardinalityZeroToMany(),
				AllowedShapes: types.NewStringSet(ops.FieldShapePrimitive),
			},
			FieldGroupStreamBody: {
				Cardinality:   types.CardinalityOneToOne(),
				AllowedShapes: types.NewStringSet(ops.FieldShapeContent),
			},
		},
		FieldPostProcessor: r.PostProcessField,
	}

	return reflector.ReflectPortStruct(PortTypeOutput, st)
}

func (r PortReflector) PostProcessField(pf *ops.PortField, sf reflect.StructField) {
	if PortReflectorPostProcessField != nil {
		PortReflectorPostProcessField(pf, sf)
	}

	// Input and Output message body are processed as Content (encoded/marshalled by Content-Type/Content-Encoding)
	if pf.Group == FieldGroupStreamBody && pf.Type.Shape != ops.FieldShapeUnknown {
		pf.Type.Shape = ops.FieldShapeContent
	}
}

var PortReflectorPostProcessField ops.PortFieldPostProcessorFunc
