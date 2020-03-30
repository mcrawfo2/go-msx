package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/emicklei/go-restful"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
)

const (
	requestTagSourceBody   = "body"
	requestTagSourceQuery  = "query"
	requestTagSourceHeader = "header"
	requestTagSourcePath   = "path"
	requestTagSourceForm   = "form"
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
		case requestTagSourceForm:
			if err := populateForm(req, fieldValue, requestTagParts[1], structField.Name); err != nil {
				return err
			}
		}

		fieldInterface := fieldValue.Interface()
		if fieldInterface == nil || (fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()) {
			continue
		}

		if fieldValue.CanAddr() {
			// In case Validate method is declared with a pointer receiver
			fieldInterface = fieldValue.Addr().Interface()
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

	headerValues, ok := req.Request.Header[headerName]
	if !ok || len(headerValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr {
			return errors.Errorf("Missing non-optional header %q", headerName)
		}
		return nil
	}

	headerValue := headerValues[0]
	return populateScalar(fieldValue, headerValue, fieldName)
}

func populatePath(req *restful.Request, fieldValue reflect.Value, pathTag, fieldName string) error {
	var pathName string
	if pathTag != "" {
		pathName = pathTag
	} else {
		pathName = strcase.ToLowerCamel(fieldName)
	}

	pathValue := req.PathParameter(pathName)

	return populateScalar(fieldValue, pathValue, fieldName)
}

func populateQuery(req *restful.Request, fieldValue reflect.Value, queryTag, fieldName string) error {
	var queryName string
	if queryTag != "" {
		queryName = queryTag
	} else {
		queryName = strcase.ToLowerCamel(fieldName)
	}

	queryValues, ok := req.Request.URL.Query()[queryName]
	if !ok || len(queryValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr {
			return errors.Errorf("Missing non-optional query parameter %q", fieldName)
		}
		return nil
	}

	queryValue := queryValues[0]
	return populateScalar(fieldValue, queryValue, fieldName)
}

func populateForm(req *restful.Request, fieldValue reflect.Value, formTag, fieldName string) error {
	var formName string
	if formTag != "" {
		formName = formTag
	} else {
		formName = strcase.ToLowerCamel(fieldName)
	}

	if req.Request.PostForm == nil {
		err := req.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			return err
		}
	}

	formFiles, ok := req.Request.MultipartForm.File[formName]
	if ok {
		return populateFile(fieldValue, formFiles[0], fieldName)
	}

	formValues, ok := req.Request.MultipartForm.Value[formName]
	if !ok || len(formValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr {
			return errors.Errorf("Missing non-optional form field %q", fieldName)
		}
		return nil
	}

	formValue := formValues[0]
	return populateScalar(fieldValue, formValue, fieldName)
}

func populateFile(fieldValue reflect.Value, header *multipart.FileHeader, fieldName string) error {
	if fieldValue.Kind() != reflect.Ptr {
		return errors.Errorf("Cannot marshal multi-part file header into field %q", fieldName)
	}

	if fieldValue.Type() != reflect.TypeOf(header) {
		return errors.Errorf("Cannot marshal multi-part file header into field %q", fieldName)
	}

	fieldValue.Set(reflect.ValueOf(header))
	return nil
}

func populateScalar(fieldValue reflect.Value, value, fieldName string) (err error) {
	errorWrapper := func(err error) error {
		return errors.Wrapf(err, "Cannot marshal string %q into field %q", value, fieldName)
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("Cannot marshal string %q into field %q", value, fieldName)
		}
	}()

	// Scalars
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		return nil

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(boolValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(floatValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(intValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(uintValue).Convert(fieldValue.Type()))
		return nil
	}

	// Nil Pointers
	fieldType := fieldValue.Type()
	if fieldValue.IsNil() {
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(reflect.New(fieldType.Elem()))
		} else {
			fieldValue.Set(reflect.New(fieldType).Elem())
		}
	}

	// Pointers to scalars
	if fieldValue.Kind() == reflect.Ptr {
		switch fieldValue.Elem().Kind() {
		case reflect.String:
			ptrValue := &value
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return errorWrapper(err)
			}
			ptrValue := &boolValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errorWrapper(err)
			}
			ptrValue := &floatValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			ptrValue := &intValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			ptrValue := &uintValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil
		}
	}

	// Unmarshaler
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

	return NewInternalError(errors.Errorf("Cannot marshal string %q into field %q", value, fieldName))
}
