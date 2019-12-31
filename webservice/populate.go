package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/emicklei/go-restful"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

const (
	requestTagSourceBody   = "body"
	requestTagSourceQuery  = "query"
	requestTagSourceHeader = "header"
	requestTagSourcePath   = "path"
)

func Populate(req *restful.Request, params interface{}) error {
	var paramsType = reflect.TypeOf(params)
	var paramsValue = reflect.ValueOf(params)

	if paramsType.Kind() != reflect.Ptr {
		return NewInternalError(errors.New("Parameters value not a pointer to struct"))
	}
	paramsType = paramsType.Elem()
	paramsValue = paramsValue.Elem()

	if paramsType.Kind() != reflect.Struct {
		return NewInternalError(errors.New("Parameters value not a pointer to struct"))
	}

	for i := 0; i < paramsType.NumField(); i++ {
		var structField = paramsType.Field(i)
		var fieldValue = paramsValue.FieldByName(structField.Name)

		if !fieldValue.CanSet() || !fieldValue.IsValid() {
			return NewBadRequestError(errors.Errorf("Cannot set field %s", structField.Name))
		}

		var fieldTag = structField.Tag
		var requestTag = fieldTag.Get("req")
		if requestTag == "" {
			continue
		}

		requestTagParts := append(strings.Split(requestTag, "="), "")[:2]
		switch requestTagParts[0] {
		case requestTagSourceBody:
			if err := populateBody(req, fieldValue); err != nil {
				return err
			}
		case requestTagSourceHeader:
			if err := populateHeader(req, fieldValue, requestTagParts[1], structField.Name); err != nil {
				return err
			}
		case requestTagSourcePath:
			if err := populatePath(req, fieldValue, requestTagParts[1], structField.Name); err != nil {
				return err
			}
		case requestTagSourceQuery:
			if err := populateQuery(req, fieldValue, requestTagParts[1], structField.Name); err != nil {
				return err
			}
		}

		fieldInterface := fieldValue.Interface()
		if fieldInterface == nil {
			continue
		}

		if paramsFieldValidatable, ok := fieldInterface.(validate.Validatable); ok {
			if err := validate.Validate(paramsFieldValidatable); err != nil {
				return NewBadRequestError(err)
			}
		}
	}

	return nil
}

func populateBody(req *restful.Request, fieldValue reflect.Value) error {
	var val = fieldValue.Addr().Interface()
	if err := req.ReadEntity(val); err != nil {
		return NewBadRequestError(err)
	}
	fieldValue.Set(reflect.ValueOf(val).Elem())
	return nil
}

func populateHeader(req *restful.Request, fieldValue reflect.Value, headerTag, fieldName string) error {
	var headerName string
	if headerTag != "" {
		headerName = strcase.ToKebab(headerTag)
	} else {
		headerName = strcase.ToKebab(fieldName)
	}

	headerValue := req.HeaderParameter(headerName)

	return populateValue(fieldValue, headerValue, fieldName)
}

func populatePath(req *restful.Request, fieldValue reflect.Value, pathTag, fieldName string) error {
	var pathName string
	if pathTag != "" {
		pathName = pathTag
	} else {
		pathName = strcase.ToLowerCamel(fieldName)
	}

	pathValue := req.PathParameter(pathName)

	return populateValue(fieldValue, pathValue, fieldName)
}

func populateQuery(req *restful.Request, fieldValue reflect.Value, queryTag, fieldName string) error {
	var queryName string
	if queryTag != "" {
		queryName = queryTag
	} else {
		queryName = strcase.ToLowerCamel(fieldName)
	}

	pathValue := req.QueryParameter(queryName)

	return populateValue(fieldValue, pathValue, fieldName)
}

func populateValue(fieldValue reflect.Value, value, fieldName string) error {
	if fieldValue.Kind() == reflect.String {
		fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		return nil
	}

	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.Elem().Kind() == reflect.String {
			ptrValue := &value
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil
		}
	}

	fieldType := fieldValue.Type()
	if fieldValue.IsNil() {
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(reflect.New(fieldType.Elem()))
		} else {
			fieldValue.Set(reflect.New(fieldType).Elem())
		}
	}

	fieldInterface := fieldValue.Interface()
	if fieldUnmarshaler, ok := fieldInterface.(types.TextUnmarshaler); ok {
		err := fieldUnmarshaler.UnmarshalText(value)
		if err != nil {
			return NewBadRequestError(err)
		} else {
			return nil
		}
	} else {
		fieldPointerInterface := fieldValue.Addr().Interface()
		if fieldUnmarshaler, ok := fieldPointerInterface.(types.TextUnmarshaler); ok {
			err := fieldUnmarshaler.UnmarshalText(value)
			if err != nil {
				return NewBadRequestError(err)
			} else {
				return nil
			}
		}
	}

	return NewInternalError(errors.Errorf("Cannot marshal string into field %s", fieldName))
}
