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

var ErrInvalidCardinality = errors.New("Field group occurs an invalid number of times in port struct")
var ErrInvalidShape = errors.New("Field group does not allow field of specified shape in port struct")

type FieldGroup struct {
	Cardinality   types.CardinalityRange
	AllowedShapes types.StringSet
	Content       bool
}

type PortDirection bool

const (
	PortDirectionIn  = PortDirection(true)
	PortDirectionOut = PortDirection(false)
)

type PortFieldPostProcessorFunc func(*PortField, reflect.StructField)

type PortReflector struct {
	Direction          PortDirection
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

	if err = r.validateFieldGroups(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortReflector) validateFieldGroups(p *Port) error {
	for name, fieldGroup := range r.FieldGroups {
		fields := p.Fields.All(PortFieldHasGroup(name))
		if err := r.validateFieldGroupCardinality(name, fieldGroup, len(fields)); err != nil {
			return err
		}
		if err := r.validateFieldGroupShapes(name, fieldGroup, fields); err != nil {
			return err
		}
	}
	return nil
}

func (r PortReflector) validateFieldGroupShapes(fieldGroupName string, fieldGroup FieldGroup, fields PortFields) error {
	for _, field := range fields {
		if !fieldGroup.AllowedShapes.Contains(field.Type.Shape) {
			return r.invalidShape(fieldGroupName, fieldGroup, field)
		}
	}
	return nil
}

func (r PortReflector) validateFieldGroupCardinality(name string, fieldGroup FieldGroup, count int) error {
	switch fieldGroup.Cardinality.Min {
	case types.CardinalityZero:
	case types.CardinalityOne:
		if count < 1 {
			return r.invalidCardinality(name, fieldGroup, count)
		}
	}

	switch fieldGroup.Cardinality.Max {
	case types.CardinalityZero:
		if count > 0 {
			return r.invalidCardinality(name, fieldGroup, count)
		}
	case types.CardinalityOne:
		if count > 1 {
			return r.invalidCardinality(name, fieldGroup, count)
		}
	}

	return nil
}

func (r PortReflector) invalidCardinality(name string, fieldGroup FieldGroup, count int) error {
	return errors.Wrapf(
		ErrInvalidCardinality,
		"Group %q Min %d Max %d Found %d",
		name,
		fieldGroup.Cardinality.Min,
		fieldGroup.Cardinality.Max,
		count)
}

func (r PortReflector) invalidShape(name string, group FieldGroup, field *PortField) error {
	return errors.Wrapf(
		ErrInvalidShape,
		"Group %q Allowed %q Found %q in field %q",
		name,
		strings.Join(group.AllowedShapes.Values(), ","),
		field.Type.Shape,
		field.Name)
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

type FieldVisitor struct {
	Indices []int
}

func (v *FieldVisitor) incrementIndex() {
	lastIndex := v.Indices[len(v.Indices)-1]
	v.Indices[len(v.Indices)-1] = lastIndex + 1
}

func (v *FieldVisitor) EnterAnonymousStructField(_ reflect.StructField) {
	// push
	v.Indices = append(v.Indices, 0)
}

func (v *FieldVisitor) ExitAnonymousStructField(_ reflect.StructField) {
	// pop
	v.Indices = v.Indices[:len(v.Indices)-1]
	v.incrementIndex()
}

type PortReflectorFieldVisitor struct {
	Port      *Port
	Reflector PortReflector
	*FieldVisitor
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
		shapeReflector = DefaultPortFieldTypeReflector{
			Direction: v.Reflector.Direction,
		}
	}
	fieldType, optional := shapeReflector.ReflectPortFieldType(field.Type)

	indices := append([]int{}, v.Indices...)
	f := NewPortField(name, peer, group, optional, v.Port.Type, fieldType, indices)

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

func newPortReflectorFieldVisitor(r PortReflector, port *Port) *PortReflectorFieldVisitor {
	return &PortReflectorFieldVisitor{
		Port:         port,
		Reflector:    r,
		FieldVisitor: newFieldVisitor(),
	}
}

func newFieldVisitor() *FieldVisitor {
	return &FieldVisitor{
		Indices: []int{0},
	}
}
