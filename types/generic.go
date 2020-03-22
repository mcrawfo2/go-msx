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

	structName := GetTypeName(structType)
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
	return GetTypeName(instanceType)
}

func GetTypeName(instanceType reflect.Type) string {
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	if typeName, ok := parameterizedTypeNames[instanceType]; ok {
		return typeName
	}

	typeNamePrefix, typeNameSuffix, typeName := "", "", ""

	switch instanceType.Kind() {
	case reflect.Array, reflect.Slice:
		typeNamePrefix = "List«"
		typeNameSuffix = "»"
		instanceType = instanceType.Elem()
	case reflect.Map:
		typeNamePrefix = "Map«" + GetTypeName(instanceType.Key()) + ","
		typeNameSuffix = "»"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String:
		typeName = instanceType.Name()
	}

	if typeName == "" {
		instanceTypePackagePath := instanceType.PkgPath()
		instanceTypePackageName := ""
		if instanceTypePackagePath != "" {
			instanceTypePackageName = path.Base(instanceTypePackagePath)
			if instanceTypePackageName != "" {
				instanceTypePackageName += "."
			}
		}
		typeName = instanceTypePackageName + instanceType.Name()
	}

	return typeNamePrefix + typeName + typeNameSuffix
}
