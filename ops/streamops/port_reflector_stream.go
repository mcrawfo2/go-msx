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

type PortReflector struct{}

func (r PortReflector) ReflectInputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupStreamChannel:   {Cardinality: types.CardinalityZeroToOne()},
			FieldGroupStreamMessageId: {Cardinality: types.CardinalityZeroToOne()},
			FieldGroupStreamHeader:    {Cardinality: types.CardinalityZeroToMany()},
			FieldGroupStreamBody:      {Cardinality: types.CardinalityOneToOne()},
		},
		FieldPostProcessor: r.PostProcessField,
	}

	return reflector.ReflectPortStruct(ops.PortTypeInput, st)
}

func (r PortReflector) ReflectOutputPort(st reflect.Type) (*ops.Port, error) {
	reflector := ops.PortReflector{
		FieldGroups: map[string]ops.FieldGroup{
			FieldGroupStreamChannel:   {Cardinality: types.CardinalityZeroToOne()},
			FieldGroupStreamMessageId: {Cardinality: types.CardinalityZeroToOne()},
			FieldGroupStreamHeader:    {Cardinality: types.CardinalityZeroToMany()},
			FieldGroupStreamBody:      {Cardinality: types.CardinalityOneToOne()},
		},
		FieldPostProcessor: r.PostProcessField,
	}

	return reflector.ReflectPortStruct(ops.PortTypeOutput, st)
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
