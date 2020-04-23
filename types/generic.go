package types

import (
	"fmt"
	"path"
	"reflect"
)

var parameterizedTypeNames = make(map[reflect.Type]string)

func NewParameterizedStruct(structType reflect.Type, payloadField string, payload interface{}) reflect.Type {
	var structFields []reflect.StructField
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		if structField.Name == payloadField {
			if payload == nil {
				continue
			} else {
				structField.Type = reflect.TypeOf(payload)
			}
		}
		structFields = append(structFields, structField)
	}

	structName := GetTypeName(structType, false)
	payloadTypeName := GetInstanceTypeName(payload)

	parameterizedStructType := reflect.StructOf(structFields)
	parameterizedTypeNames[parameterizedStructType] = fmt.Sprintf("%s«%s»", structName, payloadTypeName)
	return parameterizedStructType
}

func GetInstanceTypeName(instance interface{}) string {
	if instance == nil {
		return "Void"
	}

	instanceType := reflect.TypeOf(instance)
	return GetTypeName(instanceType, true)
}

func GetTypeName(instanceType reflect.Type, root bool) string {
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	if typeName, ok := parameterizedTypeNames[instanceType]; ok {
		return typeName
	}

	typeNamePrefix, typeNameSuffix, typeName := "", "", ""

	switch instanceType.Kind() {
	case reflect.Array, reflect.Slice:
		if root {
			return GetTypeName(instanceType.Elem(), false)
		} else {
			typeNamePrefix = "List«"
			typeNameSuffix = "»"
			typeName = GetTypeName(instanceType.Elem(), false)
			return typeNamePrefix + typeName + typeNameSuffix
		}

	case reflect.Map:
		typeNamePrefix = "Map«" + GetTypeName(instanceType.Key(), false) + ","
		typeNameSuffix = "»"
		typeName = GetTypeName(instanceType.Elem(), false)
		return typeNamePrefix + typeName + typeNameSuffix
	}

	switch instanceType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String:
		return instanceType.Name()

	default:
		instanceTypePackagePath := instanceType.PkgPath()
		instanceTypePackageName := ""
		if instanceTypePackagePath != "" {
			instanceTypePackageName = path.Base(instanceTypePackagePath)
			if instanceTypePackageName != "" {
				instanceTypePackageName += "."
			}
		} else {
			return ""
		}
		return instanceTypePackageName + instanceType.Name()
	}
}
