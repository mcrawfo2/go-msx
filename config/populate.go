// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	tagOptionIgnore   = "ignore"
	tagOptionDefault  = "default"
	tagOptionOptional = "optional"
	tagOptionKey      = "key"
)

type EnumUnpacker interface {
	Unpack(s string) error
}

type PopulatorSource interface {
	Value(key string) (Value, error)
	ValuesWithPrefix(prefix string) SnapshotValues
	ExpressionResolver
}

var populators sync.Map

func populator(t reflect.Type) (Populator, error) {
	if u, ok := populators.Load(t); ok {
		return u.(Populator), nil
	}

	u, err := newValuePopulator(t, nil)
	if err != nil {
		return nil, err
	}

	// Don't cache non-struct populator
	if (t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct) &&
		t.Kind() != reflect.Struct {
		return u, err
	}

	populators.LoadOrStore(t, u)
	return u, nil
}

func Populate(target interface{}, prefix string, source PopulatorSource) error {
	t := reflect.TypeOf(target)
	u, err := populator(t)
	if err != nil {
		return errors.Wrap(err, "Failed to create populator")
	}

	v := reflect.ValueOf(target)
	err = u.Populate(v, source, prefix)
	if err != nil {
		return errors.Wrapf(err, "Failed to populate target from %q", prefix)
	}

	if validatableTarget, ok := target.(validate.Validatable); ok {
		if err = validate.Validate(validatableTarget); err != nil {
			return err
		}
	}

	return nil
}

type Populator interface {
	Populate(target reflect.Value, source PopulatorSource, prefix string) error
}

type selfPopulator struct{}

func (u selfPopulator) Populate(target reflect.Value, source PopulatorSource, prefix string) error {
	other, ok := target.Interface().(Populator)
	if !ok {
		return errors.Wrap(ErrNoPopulator, prefix)
	}
	return other.Populate(target, source, prefix)
}

type structPopulator struct {
	StructType reflect.Type
	Fields     []Populator
}

func (u structPopulator) Populate(v reflect.Value, source PopulatorSource, prefix string) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			subsource := source.ValuesWithPrefix(prefix)
			if subsource.Empty() {
				return nil
			}

			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	for _, field := range u.Fields {
		err := field.Populate(v, source, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func newStructPopulator(t reflect.Type) (Populator, error) {
	u := structPopulator{
		StructType: t,
	}

	for i := 0; i < u.StructType.NumField(); i++ {
		fieldUnmarshaler, err := newStructFieldPopulator(u.StructType.Field(i))
		if err != nil {
			return nil, err
		}

		if fieldUnmarshaler != nil {
			u.Fields = append(u.Fields, fieldUnmarshaler)
		}

	}

	return u, nil
}

type structFieldPopulator struct {
	FieldType      reflect.StructField
	NormalizedName string
	IsOptional     bool
	Value          Populator
}

func (u structFieldPopulator) Populate(v reflect.Value, source PopulatorSource, prefix string) error {
	fieldValue := v.FieldByName(u.FieldType.Name)
	key := PrefixWithName(prefix, u.NormalizedName)
	err := u.Value.Populate(fieldValue, source, key)
	if errors.Is(err, ErrNotFound) && u.IsOptional {
		err = nil
	} else if err != nil {
		err = errors.Wrapf(err, "Failed to populate value %q", PrefixWithName(prefix, u.FieldType.Name))
	}
	return err
}

func newStructFieldPopulator(fieldType reflect.StructField) (Populator, error) {
	var err error

	fieldTag := fieldType.Tag.Get("config")
	fieldOptions := parseStructTag(fieldTag)

	if fieldOptions[tagOptionIgnore] != "" {
		return nil, nil
	}

	var defaultValue expression
	if defaultValueString, ok := fieldOptions[tagOptionDefault]; ok {
		defaultValue, err = parseExpression(defaultValueString)
		if err != nil {
			return nil, err
		}
	}

	valueUnmarshaler, err := newValuePopulator(fieldType.Type, defaultValue)
	if err != nil {
		return nil, err
	}

	key := fieldType.Name
	if fieldOptions[tagOptionKey] != "" {
		key = fieldOptions[tagOptionKey]
	}
	key = NormalizeKey(key)

	optional := false
	if fieldOptions[tagOptionOptional] != "" {
		optional = true
	}

	return structFieldPopulator{
		FieldType:      fieldType,
		NormalizedName: key,
		IsOptional:     optional,
		Value:          valueUnmarshaler,
	}, nil
}

func parseStructTag(tag string) map[string]string {
	if tag == "" {
		return nil
	}

	var options = make(map[string]string)
	for j, option := range strings.Split(tag, ",") {
		optionParts := strings.SplitN(option, "=", 2)

		switch j {
		case 0:
			if option == "-" {
				// No source, completely ignore
				options[tagOptionIgnore] = ""
				return options
			} else if option == tagOptionOptional {
				options[tagOptionOptional] = "true"
				break
			} else if len(optionParts) == 1 {
				options[tagOptionKey] = optionParts[0]
				break
			}

			fallthrough
		default:
			var value = "true"
			if len(optionParts) == 2 {
				value = optionParts[1]
			}
			options[optionParts[0]] = value
		}
	}

	return options
}

type scalarPopulator struct {
	DefaultValue expression
	Setter       func(value Value, target reflect.Value) error
}

func setInt(value Value, target reflect.Value) error {
	if unpacker, ok := target.Interface().(EnumUnpacker); ok {
		return unpacker.Unpack(string(value))
	} else if unpacker, ok = target.Addr().Interface().(EnumUnpacker); ok {
		return unpacker.Unpack(string(value))
	}

	typedValue, err := value.Int()
	if err != nil {
		return errors.Wrap(err, "Failed to parse int")
	}
	target.SetInt(typedValue)
	return nil
}

func setFloat(value Value, target reflect.Value) error {
	typedValue, err := value.Float()
	if err != nil {
		return errors.Wrap(err, "Failed to parse float")
	}
	target.SetFloat(typedValue)
	return nil
}

func setBool(value Value, target reflect.Value) error {
	typedValue, err := value.Bool()
	if err != nil {
		return errors.Wrap(err, "Failed to parse bool")
	}
	target.SetBool(typedValue)
	return nil
}

func setString(value Value, target reflect.Value) error {
	typedValue := value.String()
	target.SetString(typedValue)
	return nil
}

func setUint(value Value, target reflect.Value) error {
	if unpacker, ok := target.Interface().(EnumUnpacker); ok {
		return unpacker.Unpack(string(value))
	} else if unpacker, ok = target.Addr().Interface().(EnumUnpacker); ok {
		return unpacker.Unpack(string(value))
	}

	typedValue, err := value.Uint()
	if err != nil {
		return errors.Wrap(err, "Failed to parse unsigned int")
	}
	target.SetUint(typedValue)
	return nil
}

func setDuration(value Value, target reflect.Value) error {
	typedValue, err := value.Duration()
	if err != nil {
		return errors.Wrap(err, "Failed to parse duration")
	}
	target.Set(reflect.ValueOf(typedValue).Convert(target.Type()))
	return nil
}

func (u scalarPopulator) Populate(v reflect.Value, source PopulatorSource, key string) error {
	if !v.CanSet() {
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return errors.Wrap(ErrValueCannotBeSet, key)
		}
	}

	value, err := source.Value(key)
	if errors.Is(err, ErrNotFound) {
		if u.DefaultValue == nil {
			return errors.Wrap(err, key)
		}

		vs, err := u.DefaultValue.Resolve(source)
		if err != nil {
			return err
		}

		value = Value(vs)
	}

	if v.Kind() == reflect.Ptr {
		if v.CanSet() && v.IsNil() {
			vp := reflect.New(v.Type().Elem())
			v.Set(vp)
			v = vp.Elem()
		} else {
			v = v.Elem()
		}
	}

	return u.Setter(value, v)
}

type mapPopulator struct {
	Value Populator
}

func (u mapPopulator) Populate(v reflect.Value, source PopulatorSource, prefix string) error {
	valueType := v.Type().Elem()
	isPtr := valueType.Kind() == reflect.Ptr
	if isPtr {
		valueType = valueType.Elem()
	}

	if v.IsNil() {
		keyType := v.Type().Key()
		valType := v.Type().Elem()
		mapType := reflect.MapOf(keyType, valType)
		v.Set(reflect.MakeMapWithSize(mapType, 0))
	}

	for _, childNodeName := range source.ValuesWithPrefix(prefix).Entries().ChildNodeNames(prefix) {
		key := reflect.ValueOf(childNodeName.Name)
		value := reflect.New(valueType)

		err := u.Value.Populate(value, source, childNodeName.NormalizedName)
		if err != nil {
			return errors.Wrapf(err, "Failed to populate value %q", PrefixWithName(prefix, childNodeName.Name))
		}

		if isPtr {
			v.SetMapIndex(key, value)
		} else if !value.IsNil() {
			v.SetMapIndex(key, value.Elem())
		}
	}

	return nil
}

type slicePopulator struct {
	Value        Populator
	DefaultValue expression
}

func (u slicePopulator) Populate(v reflect.Value, source PopulatorSource, prefix string) error {
	// Slice type
	valueType := v.Type()
	valueTypeIsPtr := valueType.Kind() == reflect.Ptr
	if valueTypeIsPtr {
		valueType = valueType.Elem()
		v = v.Elem()
	}

	// Element type
	elemType := valueType.Elem()
	elemTypeIsPtr := elemType.Kind() == reflect.Ptr
	if elemTypeIsPtr {
		elemType = elemType.Elem()
	}

	childNodeNames := source.ValuesWithPrefix(prefix).Entries().ChildNodeNames(prefix)
	if len(childNodeNames) == 0 {
		if resolvedEntry, err := source.ResolveByName(prefix); err == nil {
			entries := resolvedEntry.ResolvedValue.StringSlice(",")
			return u.populateFromEntries(v, entries)
		} else if !errors.Is(err, ErrNotFound) {
			return err
		}

		if u.DefaultValue != nil {
			defaultValue, err := u.DefaultValue.Resolve(source)
			if err != nil {
				return err
			}

			entries := Value(defaultValue).StringSlice(";")
			return u.populateFromEntries(v, entries)
		}

		return nil
	}

	sort.Slice(childNodeNames, func(i, j int) bool {
		return childNodeNames[i].Index < childNodeNames[j].Index
	})

	sliceType := reflect.SliceOf(valueType.Elem())
	v.Set(reflect.MakeSlice(sliceType, len(childNodeNames), len(childNodeNames)))

	for i, childNodeName := range childNodeNames {
		childValue := reflect.New(elemType)

		err := u.Value.Populate(childValue, source, childNodeName.NormalizedName)
		if err != nil {
			return errors.Wrapf(err, "Failed to populate value %q", PrefixWithIndex(prefix, childNodeName.Index))
		}

		if elemTypeIsPtr {
			v.Index(i).Set(childValue)
		} else if !childValue.IsNil() {
			v.Index(i).Set(childValue.Elem())
		}
	}

	return nil
}

func (u slicePopulator) populateFromEntries(v reflect.Value, entries []string) error {
	// Slice type
	valueType := v.Type()
	valueTypeIsPtr := valueType.Kind() == reflect.Ptr
	if valueTypeIsPtr {
		valueType = valueType.Elem()
		v = v.Elem()
	}

	// Element type
	elemType := valueType.Elem()
	elemTypeIsPtr := elemType.Kind() == reflect.Ptr
	if elemTypeIsPtr {
		elemType = elemType.Elem()
	}

	entriesSource := customUnmarshalerSource{
		valuer: sliceValues{
			values: entries,
		},
	}

	sliceType := reflect.SliceOf(valueType.Elem())
	v.Set(reflect.MakeSlice(sliceType, len(entries), len(entries)))

	for i := range entries {
		childValue := reflect.New(elemType)

		err := u.Value.Populate(childValue, entriesSource, strconv.Itoa(i))
		if err != nil {
			return err
		}

		if elemTypeIsPtr {
			v.Index(i).Set(childValue)
		} else if !childValue.IsNil() {
			v.Index(i).Set(childValue.Elem())
		}
	}

	return nil
}

var populatorType = reflect.TypeOf((*Populator)(nil)).Elem()
var durationType = reflect.TypeOf(time.Duration(0))

func newValuePopulator(valueType reflect.Type, defaultValue expression) (Populator, error) {
	// Self-populators
	if valueType.Implements(populatorType) {
		return selfPopulator{}, nil
	}

	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	// Concrete type populators
	switch valueType {
	case durationType:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setDuration}, nil
	}

	// Kinds
	switch valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setInt}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setUint}, nil
	case reflect.String:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setString}, nil
	case reflect.Bool:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setBool}, nil
	case reflect.Float64, reflect.Float32:
		return scalarPopulator{DefaultValue: defaultValue, Setter: setFloat}, nil
	case reflect.Struct:
		return newStructPopulator(valueType)
	case reflect.Map:
		valueUnmarshaler, err := newValuePopulator(valueType.Elem(), nil)
		if err != nil {
			return nil, err
		}

		return mapPopulator{
			Value: valueUnmarshaler,
		}, nil
	case reflect.Slice:
		valueUnmarshaler, err := newValuePopulator(valueType.Elem(), nil)
		if err != nil {
			return nil, err
		}

		return slicePopulator{
			Value:        valueUnmarshaler,
			DefaultValue: defaultValue,
		}, nil
	}

	return nil, errors.Wrapf(ErrNoPopulator, "%q", valueType)
}

type valuer interface {
	Value(key string) (Value, error)
}

type sliceValues struct {
	values []string
}

func (s sliceValues) Value(key string) (Value, error) {
	idx, _ := strconv.Atoi(key)
	return Value(s.values[idx]), nil
}

type customUnmarshalerSource struct {
	valuer
}

func (s customUnmarshalerSource) ValuesWithPrefix(_ string) SnapshotValues {
	return emptySnapshotValues
}

func (s customUnmarshalerSource) ResolveByName(name string) (ResolvedEntry, error) {
	v, err := s.Value(name)
	if err != nil {
		return ResolvedEntry{}, errors.Wrapf(err, "Failed to resolve %q", name)
	}

	return ResolvedEntry{
		ResolvedValue: v,
	}, nil
}
