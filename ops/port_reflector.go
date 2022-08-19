// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

type FieldGroup struct {
	Cardinality   types.CardinalityRange
	AllowedShapes types.StringSet
	Content       bool
}

type PortFieldPostProcessorFunc func(*PortField, reflect.StructField)

type PortReflector struct {
	FieldGroups        map[string]FieldGroup
	FieldPostProcessor PortFieldPostProcessorFunc
	FieldTypeReflector PortFieldTypeReflector
}

func (r PortReflector) ReflectPortStruct(typ string, st reflect.Type) (*Port, error) {
	p, err := NewPort(typ, st)
	if err != nil {
		return nil, err
	}

	v := newPortReflectorFieldVisitor(r, p)

	if err = WalkStruct(p.StructType, v); err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortReflector) fieldGroupNames() string {
	var results []string
	for k := range r.FieldGroups {
		results = append(results, k)
	}
	return strings.Join(results, ",")
}

func (r PortReflector) cardinality(groupName string) types.CardinalityRange {
	group, ok := r.FieldGroups[groupName]
	if !ok {
		return types.CardinalityNone()
	}
	return group.Cardinality
}

type PortReflectorFieldVisitor struct {
	Port      *Port
	Reflector PortReflector
	Indices   []int
}

func (v *PortReflectorFieldVisitor) reflectPortField(field reflect.StructField) (*PortField, error) {
	// Skip fields with invalid tags
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse field %q tag: %q", field.Name, string(field.Tag))
	}

	// Skip fields with no port tag
	portTag, ok := LookupTag(tags, v.Port.Type)
	if !ok {
		return nil, nil
	}

	// Skip ignored fields
	if portTag.Name == "-" {
		return nil, nil
	}

	// Parse the name
	sourceParts := strings.SplitN(portTag.Name, "=", 2)
	group := sourceParts[0]
	if cardinality := v.Reflector.cardinality(group); cardinality.Max == types.CardinalityZero {
		return nil, errors.Errorf(
			"Invalid group %q, exepected one of %q",
			group,
			v.Reflector.fieldGroupNames())
	}
	name := field.Name
	peer := ""
	if len(sourceParts) == 2 {
		peer = sourceParts[1]
	}

	// Parse shape
	var shapeReflector = v.Reflector.FieldTypeReflector
	if shapeReflector == nil {
		shapeReflector = DefaultPortFieldTypeReflector{}
	}
	fieldType, optional := shapeReflector.ReflectPortFieldType(field.Type)

	indices := append([]int{}, v.Indices...)
	f := NewPortField(name, peer, group, optional, fieldType, indices)

	// Apply all tag to the options set
	v.reflectPrimaryTag(f, portTag)
	if requiredTag, _ := tags.Get("required"); requiredTag != nil && requiredTag.Name != "" {
		v.reflectSupplementalTag(f, requiredTag)
	}
	if optionalTag, _ := tags.Get("optional"); optionalTag != nil && optionalTag.Name != "" {
		v.reflectSupplementalTag(f, optionalTag)
	}
	for _, tag := range tags.Tags() {
		if tag.Key != v.Port.Type && tag.Key != "required" && tag.Key != "optional" {
			v.reflectSupplementalTag(f, tag)
		}
	}

	// Allow overriding everything above
	if v.Reflector.FieldPostProcessor != nil {
		v.Reflector.FieldPostProcessor(f, field)
	}

	// Override the peer name if still unset
	if f.Peer == "" {
		f.Peer = strcase.ToLowerCamel(f.Name)
	}

	return f, nil
}

func (v PortReflectorFieldVisitor) reflectPrimaryTag(p *PortField, tag *structtag.Tag) *PortField {
	for _, option := range tag.Options {
		optionParts := strings.SplitN(option, "=", 2)

		var name = optionParts[0]
		var value = "true"
		if len(optionParts) == 2 {
			value = optionParts[1]
		}

		if optionParts[0] == "optional" {
			p.WithOptional(value == "true")
			name = "required"
			value = strconv.FormatBool(!p.Optional)
		} else if optionParts[0] == "required" {
			p.WithOptional(value != "true")
		}

		p.WithOption(name, value)
	}

	return p
}

func (v PortReflectorFieldVisitor) reflectSupplementalTag(p *PortField, tag *structtag.Tag) *PortField {
	switch tag.Key {
	case "required":
		p.WithOptional(tag.Value() != "true")
	case "optional":
		p.WithOptional(tag.Value() == "true")
	case "san":
		p.WithBoolOptionDefault("san", true)
	}

	p.WithOption(tag.Key, tag.Value())

	return p
}

func (v *PortReflectorFieldVisitor) incrementIndex() {
	lastIndex := v.Indices[len(v.Indices)-1]
	v.Indices[len(v.Indices)-1] = lastIndex + 1
}

func (v *PortReflectorFieldVisitor) VisitField(f reflect.StructField) error {
	pf, err := v.reflectPortField(f)
	if err != nil {
		return err
	} else if pf != nil {
		v.Port = v.Port.WithField(pf)
	}
	v.incrementIndex()
	return nil
}

func (v *PortReflectorFieldVisitor) EnterAnonymousStructField(_ reflect.StructField) {
	// push
	v.Indices = append(v.Indices, 0)
}

func (v *PortReflectorFieldVisitor) ExitAnonymousStructField(_ reflect.StructField) {
	// pop
	v.Indices = v.Indices[:len(v.Indices)-1]
	v.incrementIndex()
}

func newPortReflectorFieldVisitor(r PortReflector, port *Port) *PortReflectorFieldVisitor {
	return &PortReflectorFieldVisitor{
		Port:      port,
		Reflector: r,
		Indices:   []int{0},
	}
}
