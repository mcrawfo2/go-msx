// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"strconv"
	"strings"
)

type PortFields []*PortField

func (f PortFields) All(predicates ...PortFieldPredicate) PortFields {
	var result PortFields
	for _, v := range f {
		match := true
		for _, predicate := range predicates {
			if match = match && predicate(v); !match {
				break
			}
		}
		if match {
			result = append(result, v)
		}
	}
	return result
}

func (f PortFields) First(predicates ...PortFieldPredicate) *PortField {
	for _, v := range f {
		match := true
		for _, predicate := range predicates {
			if match = match && predicate(v); !match {
				break
			}
		}
		if match {
			return v
		}
	}
	return nil
}

type PortFieldPredicate func(p *PortField) bool

func PortFieldHasGroup(group string) PortFieldPredicate {
	return func(p *PortField) bool {
		return p.Group == group
	}
}

func PortFieldHasName(name string) PortFieldPredicate {
	return func(p *PortField) bool {
		return p.Name == name
	}
}

func PortFieldHasPeer(peer string) PortFieldPredicate {
	return func(p *PortField) bool {
		return p.Peer == peer
	}
}

type PortFieldElementType struct {
	Indices  []int
	Optional bool
	PortFieldType
}

func (t PortFieldElementType) WithIndirections(n int) PortFieldElementType {
	t.PortFieldType.Indirections = n
	return t
}

type PortFieldType struct {
	Shape        string // Primitive, Array, Object, File, FileArray, Reader, Unknown
	Type         reflect.Type
	Indirections int
	HandlerType  reflect.Type
	Items        *PortFieldElementType  // array/slice elements
	Keys         *PortFieldElementType  // map keys
	Values       *PortFieldElementType  // map values
	Fields       []PortFieldElementType // struct fields
}

func PortFieldTypeFromType(t reflect.Type, shape string) PortFieldType {
	return PortFieldType{
		Shape:        shape,
		Type:         t,
		Indirections: 0,
		HandlerType:  t,
	}
}

func (p *PortFieldType) IncIndirections() *PortFieldType {
	p.Indirections++
	return p
}

func (p *PortFieldType) WithHandlerType(t reflect.Type) *PortFieldType {
	p.HandlerType = t
	return p
}

type PortField struct {
	Name     string
	Indices  []int
	Peer     string
	Group    string
	Optional bool
	PortType string
	Type     PortFieldType
	Options  map[string]string
	Baggage  map[interface{}]interface{}
}

func (p *PortField) SanitizeOptions() sanitize.Options {
	return sanitize.NewOptions(p.Options["san"])
}

func (p *PortField) WithOptional(optional bool) *PortField {
	p.Optional = optional
	return p
}

func (p *PortField) BoolOption(optionName string) (bool, bool) {
	value, ok := p.Options[optionName]
	if !ok {
		return false, false
	}
	return value == "true", true
}

func (p *PortField) WithBoolOptionDefault(optionName string, value bool) *PortField {
	return p.WithOptionDefault(optionName, strconv.FormatBool(value))
}

func (p *PortField) WithOptionDefault(optionName string, value string) *PortField {
	_, exists := p.Options[optionName]
	if !exists {
		p.Options[optionName] = value
	}
	return p
}

func (p *PortField) WithOption(optionName string, value string) *PortField {
	p.Options[optionName] = value
	return p
}

func (p *PortField) WithBaggageItem(key interface{}, value interface{}) *PortField {
	p.Baggage[key] = value
	return p
}

func (p *PortField) Enum() []interface{} {
	fieldVal := reflect.New(p.Type.Type).Interface()
	if e, isEnumer := fieldVal.(jsonschema.Enum); isEnumer {
		return e.Enum()
	}

	eval, ok := p.Options["enum"]
	if ok && len(eval) > 0 {
		values := strings.Split(eval, ",")
		var result []interface{}
		for _, v := range values {
			result = append(result, v)
		}
		return result
	}

	return nil
}

func (p *PortField) ExpectShape(shape string) error {
	if p.Type.Shape != shape {
		return errors.Wrapf(ErrIncorrectShape,
			"Field %q: Expected %q but got %q",
			p.Name,
			shape,
			p.Type.Shape)
	}
	return nil
}

func (p *PortField) optionToShapedValue(optionName string) (value interface{}) {
	d, ok := p.Options[optionName]
	if !ok {
		return nil
	}

	switch p.Type.Shape {
	case FieldShapePrimitive:
		value = d
	case FieldShapeArray:
		value = strings.Split(d, ",")
	case FieldShapeObject, FieldShapeAny:
		_ = json.Unmarshal([]byte(d), &value)
		return
	}

	return

}

func (p *PortField) Default() interface{} {
	return p.optionToShapedValue("default")
}

func (p *PortField) Const() interface{} {
	return p.optionToShapedValue("const")
}

func (p *PortField) Tags() reflect.StructTag {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%s:%q ", p.PortType, fmt.Sprintf("%s=%s", p.Group, p.Peer)))
	for k, v := range p.Options {
		sb.WriteString(fmt.Sprintf("%s:%q ", k, fmt.Sprintf("%v", v)))
	}

	if p.Optional {
		sb.WriteString(`optional="true"`)
	} else {
		sb.WriteString(`optional="false"`)
	}

	return reflect.StructTag(sb.String())
}

func NewPortField(name, peer, group string, optional bool, portType string, typ PortFieldType, indices []int) *PortField {
	return &PortField{
		Name:     name,
		Indices:  indices,
		Peer:     peer,
		Group:    group,
		Optional: optional,
		PortType: portType,
		Type:     typ,
		Options:  make(map[string]string),
		Baggage:  make(map[interface{}]interface{}),
	}
}
