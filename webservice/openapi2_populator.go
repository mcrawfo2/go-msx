// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/emicklei/go-restful"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	requestTag  = "req"
	sanitizeTag = "san"

	requestTagSourceBody   = "body"
	requestTagSourceQuery  = "query"
	requestTagSourceHeader = "header"
	requestTagSourcePath   = "path"
	requestTagSourceForm   = "form"
)

func Populate(req *restful.Request, params interface{}) (err error) {
	routeParams, err := GetRouteParams(req.Request.Context(), params)
	if err != nil {
		return
	}

	var paramsValue = reflect.ValueOf(params).Elem()
	return routeParams.Populate(req, paramsValue)
}

type RouteParam struct {
	Field           reflect.StructField
	Source          string
	Name            string
	Options         map[string]string
	SanitizeOptions sanitize.Options
	Parameter       restful.ParameterData
}

func (r RouteParam) Populate(req *restful.Request, paramsValue reflect.Value) error {
	var fieldValue = paramsValue.FieldByName(r.Field.Name)

	if !fieldValue.CanSet() || !fieldValue.IsValid() {
		return NewBadRequestError(errors.Errorf("Cannot set field %s", r.Field.Name))
	}

	switch r.Source {
	case requestTagSourceBody:
		if err := r.populateBody(req, fieldValue); err != nil {
			return err
		}
	case requestTagSourceHeader:
		if err := r.populateHeader(req, fieldValue); err != nil {
			return err
		}
	case requestTagSourcePath:
		if err := r.populatePath(req, fieldValue); err != nil {
			return err
		}
	case requestTagSourceQuery:
		if err := r.populateQuery(req, fieldValue); err != nil {
			return err
		}
	case requestTagSourceForm:
		if err := r.populateForm(req, fieldValue); err != nil {
			return err
		}
	}

	fieldInterface := fieldValue.Interface()
	if fieldInterface == nil || (fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()) {
		return nil
	}

	if fieldValue.CanAddr() {
		// In case Validate method is declared with a pointer receiver
		fieldInterface = fieldValue.Addr().Interface()
	}

	if paramsFieldValidatable, ok := fieldInterface.(validate.Validatable); ok {
		if err := validate.Validate(paramsFieldValidatable); err != nil {
			err = errors.Wrap(err, r.Field.Name)
			return NewBadRequestError(err)
		}
	}

	return nil
}

func (r RouteParam) populateBody(req *restful.Request, fieldValue reflect.Value) error {
	var val = fieldValue.Addr().Interface()
	if err := req.ReadEntity(val); err != nil {
		if err == io.EOF {
			if r.Options["optional"] == "true" {
				return nil
			}
			err = errors.Wrap(err, "Missing required body")
		}
		return NewBadRequestError(err)
	}

	if r.Options["san"] == "true" {
		if err := sanitize.Input(val, r.SanitizeOptions); err != nil {
			return err
		}
	}

	fieldValue.Set(reflect.ValueOf(val).Elem())
	return nil
}

func (r RouteParam) populateHeader(req *restful.Request, fieldValue reflect.Value) error {
	headerValues := req.Request.Header.Values(r.Name)
	if len(headerValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr && r.Options["optional"] != "true" {
			return errors.Errorf("Missing non-optional header %q", r.Name)
		}
		return nil
	}

	headerValue := headerValues[0]
	return r.populateScalar(fieldValue, headerValue)
}

func (r RouteParam) populatePath(req *restful.Request, fieldValue reflect.Value) error {
	pathValue := req.PathParameter(r.Name)
	return r.populateScalar(fieldValue, pathValue)
}

func (r RouteParam) populateQuery(req *restful.Request, fieldValue reflect.Value) error {
	queryValues, ok := req.Request.URL.Query()[r.Name]
	if !ok || len(queryValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr && r.Options["optional"] != "true" {
			return errors.Errorf("Missing non-optional query parameter %q", r.Name)
		}
		return nil
	}

	return r.populateField(fieldValue, queryValues)
}

func (r RouteParam) populateForm(req *restful.Request, fieldValue reflect.Value) error {
	contentType := req.Request.Header.Get("Content-Type")
	baseContentType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	switch baseContentType {
	case MIME_MULTIPART_FORM:
		return r.populateMultipartForm(req, fieldValue)
	case MIME_APPLICATION_FORM:
		return errors.Errorf("Content-Type %q currently unsupported", baseContentType)
	default:
		if r.Options["file"] == "true" {
			return r.populateRestFormFile(req, fieldValue)
		} else {
			return r.populateQuery(req, fieldValue)
		}
	}
}

func (r RouteParam) populateRestFormFile(req *restful.Request, fieldValue reflect.Value) error {
	// Populate a Multipart Form with one file containing the request body
	body := bytes.NewBuffer(make([]byte, 0, 32768))
	w := multipart.NewWriter(body)

	contentType := w.FormDataContentType()
	boundary := strings.TrimPrefix(contentType, "multipart/form-data; boundary=")
	boundary = strings.Trim(boundary, "\"")

	// Add file from the body
	part, err := w.CreateFormFile(r.Name, "body")
	if req.Request.Body != nil {
		var bodyBytes []byte
		bodyBytes, err = ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return err
		}
		_, err = part.Write(bodyBytes)
		if err != nil {
			return err
		}
	}
	err = w.Close()
	if err != nil {
		return err
	}

	reader := multipart.NewReader(body, boundary)
	multipartForm, err := reader.ReadForm(32 << 20)
	if err != nil {
		return err
	}

	formFiles, ok := multipartForm.File[r.Name]
	if ok {
		return r.populateFile(fieldValue, formFiles[0])
	} else {
		err = errors.New("Failed to retrieve multipart form file")
	}
	return err
}

func (r RouteParam) populateMultipartForm(req *restful.Request, fieldValue reflect.Value) error {
	if req.Request.PostForm == nil {
		err := req.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			return err
		}
	}

	formFiles, ok := req.Request.MultipartForm.File[r.Name]
	if ok {
		return r.populateFile(fieldValue, formFiles[0])
	}

	formValues, ok := req.Request.MultipartForm.Value[r.Name]
	if !ok || len(formValues) == 0 {
		// Attempt a lookup in the query values
		queryValues := req.Request.URL.Query()
		formValues, ok = queryValues[r.Name]
	}
	if !ok || len(formValues) == 0 {
		if fieldValue.Kind() != reflect.Ptr && r.Options["optional"] != "true" {
			return errors.Errorf("Missing non-optional form field %q", r.Name)
		}
		return nil
	}

	formValue := formValues[0]
	return r.populateScalar(fieldValue, formValue)
}

func (r RouteParam) populateFile(fieldValue reflect.Value, header *multipart.FileHeader) error {
	if fieldValue.Kind() != reflect.Ptr {
		return errors.Errorf("Cannot marshal multi-part file header into field %q", r.Name)
	}

	if fieldValue.Type() != reflect.TypeOf(header) {
		return errors.Errorf("Cannot marshal multi-part file header into field %q", r.Name)
	}

	fieldValue.Set(reflect.ValueOf(header))
	return nil
}

func (r RouteParam) populateField(fieldValue reflect.Value, values []string) (err error) {
	// Handle pointers
	elemKind := fieldValue.Kind()
	if elemKind == reflect.Ptr {
		elemKind = fieldValue.Type().Elem().Kind()
	}

	if r.Options["multi"] == "true" {
		switch elemKind {
		case reflect.Slice:
			return r.populateSlice(fieldValue, values)
		default:
			return errors.Errorf("Cannot unmarshal multiple values into field %s", r.Name)
		}
	} else if r.Options["csv"] == "true" {
		switch elemKind {
		case reflect.Slice:
			csvValues := strings.Split(values[0], ",")
			return r.populateSlice(fieldValue, csvValues)
		default:
			return errors.Errorf("Cannot unmarshal multiple values into field %s", r.Name)
		}
	} else {
		return r.populateScalar(fieldValue, values[0])
	}

}

func (r RouteParam) populateScalar(fieldValue reflect.Value, value string) (err error) {
	errorWrapper := func(err error) error {
		return errors.Wrapf(err, "Cannot marshal string %q into field %q", value, r.Name)
	}

	defer func() {
		if v := recover(); v != nil {
			err = errors.Errorf("Cannot marshal string %q into field %q: %s", value, r.Name, v)
		}
	}()

	// Scalars
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		if r.Options["san"] == "true" {
			if err = sanitize.Input(fieldValue.Addr().Interface(), r.SanitizeOptions); err != nil {
				return err
			}
		}
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
	switch fieldValue.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
		if fieldValue.IsNil() {
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			} else {
				fieldValue.Set(reflect.New(fieldType).Elem())
			}
		}
	}

	// Pointers to scalars
	if fieldValue.Kind() == reflect.Ptr {
		switch fieldValue.Elem().Kind() {
		case reflect.String:
			if r.Options["san"] == "true" {
				if err = sanitize.Input(fieldValue.Interface(), r.SanitizeOptions); err != nil {
					return err
				}
			}
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

		case reflect.Float32:
			floatValue, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := float32(floatValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := float64(floatValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int8:
			intValue, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int8(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int16:
			intValue, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int16(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int32:
			intValue, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int32(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int64(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint8:
			uintValue, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint8(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint16:
			uintValue, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint16(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint32:
			uintValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint32(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint64:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint64(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil
		}
	}

	// Unmarshaler
	fieldInterface := fieldValue.Interface()
	if fieldUnmarshaler, ok := fieldInterface.(types.TextUnmarshaler); ok {
		err := fieldUnmarshaler.UnmarshalText(value)
		if err != nil {
			err = errors.Wrap(err, r.Name)
			return NewBadRequestError(err)
		} else {
			return nil
		}
	} else {
		fieldPointerInterface := fieldValue.Addr().Interface()
		if fieldUnmarshaler, ok := fieldPointerInterface.(types.TextUnmarshaler); ok {
			err := fieldUnmarshaler.UnmarshalText(value)
			if err != nil {
				err = errors.Wrap(err, r.Name)
				return NewBadRequestError(err)
			} else {
				return nil
			}
		}
	}

	return NewInternalError(errors.Errorf("Cannot marshal string %q into field %q", value, r.Name))
}

func (r RouteParam) populateSlice(fieldValue reflect.Value, values []string) error {
	sliceType := fieldValue.Type()
	isPtr := sliceType.Kind() == reflect.Ptr
	if isPtr {
		sliceType = sliceType.Elem()
	}

	sliceValue := reflect.MakeSlice(sliceType, len(values), len(values))

	for i, queryValue := range values {
		err := r.populateScalar(sliceValue.Index(i), queryValue)
		if err != nil {
			return err
		}
	}

	if isPtr {
		x := reflect.New(sliceType)
		x.Elem().Set(sliceValue)
		fieldValue.Set(x)
	} else {
		fieldValue.Set(sliceValue)
	}

	return nil
}

func NewRouteParam(ctx context.Context, route *restful.Route, field reflect.StructField) *RouteParam {
	r := new(RouteParam)
	r.Field = field
	r.Options = make(map[string]string)

	tag := field.Tag.Get(requestTag)

	if tag == "" {
		return nil
	}

	for j, option := range strings.Split(tag, ",") {
		optionParts := strings.SplitN(option, "=", 2)

		switch j {
		case 0:
			if option == "-" {
				// No source
				continue
			}

			r.Source = optionParts[0]
			if len(optionParts) == 2 {
				r.Name = optionParts[1]
			} else {
				switch r.Source {
				case requestTagSourceHeader:
					r.Name = strcase.ToKebab(field.Name)
				case requestTagSourcePath:
					r.Name = strcase.ToLowerCamel(field.Name)
				case requestTagSourceQuery:
					r.Name = strcase.ToLowerCamel(field.Name)
				case requestTagSourceForm:
					r.Name = strcase.ToLowerCamel(field.Name)
				case requestTagSourceBody:
					r.Name = "body"
				}
			}
		default:
			var value = "true"
			if len(optionParts) == 2 {
				value = optionParts[1]
			}
			r.Options[optionParts[0]] = value
		}
	}

	parameterFound := false
	for _, parameter := range route.ParameterDocs {
		parameterData := parameter.Data()
		if parameterData.Name == r.Name {
			r.Parameter = parameterData
			parameterFound = true
		}
	}

	if !parameterFound {
		r.Parameter = restful.ParameterData{
			Name: strcase.ToLowerCamel(r.Name),
		}

		switch r.Source {
		case requestTagSourceHeader:
			r.Parameter.Kind = restful.HeaderParameterKind
		case requestTagSourcePath:
			r.Parameter.Kind = restful.PathParameterKind
		case requestTagSourceQuery:
			r.Parameter.Kind = restful.QueryParameterKind
		case requestTagSourceBody:
			r.Parameter.Kind = restful.BodyParameterKind
		case requestTagSourceForm:
			r.Parameter.Kind = restful.FormParameterKind
		default:
			logger.
				WithContext(ctx).
				WithError(errors.Errorf("Unknown parameter source: %q", r.Source)).
				Warnf("Defining dynamic parameter %q", r.Parameter.Name)
		}
	}

	tag = field.Tag.Get(sanitizeTag)
	r.SanitizeOptions = sanitize.NewOptions(tag)
	if tag != "" {
		if _, exists := r.Options["san"]; !exists {
			r.Options["san"] = "true"
		}
	}

	return r
}

type RouteParams struct {
	Type   reflect.Type
	Fields []*RouteParam
}

func (r RouteParams) Populate(req *restful.Request, paramsValue reflect.Value) error {
	for _, routeParam := range r.Fields {
		err := routeParam.Populate(req, paramsValue)
		if err != nil {
			return err
		}
	}

	return nil
}

var routeParamsMtx sync.Mutex
var routeParamsIndex = make(map[reflect.Type]*RouteParams)

func GetRouteParams(ctx context.Context, params interface{}) (*RouteParams, error) {
	routeParamsMtx.Lock()
	defer routeParamsMtx.Unlock()

	route := RouteFromContext(ctx)
	if route == nil {
		return nil, NewInternalError(errors.New("Route not set in context"))
	}

	paramsType, err := getParamsType(params)
	if err != nil {
		return nil, err
	}

	if result, ok := routeParamsIndex[paramsType]; ok {
		return result, nil
	}

	result, err := generateRouteParams(ctx, route, paramsType)
	if err != nil {
		return nil, err
	}
	routeParamsIndex[paramsType] = result
	return result, nil
}

func getParamsType(params interface{}) (reflect.Type, error) {
	var paramsType = reflect.TypeOf(params)

	if paramsType.Kind() != reflect.Ptr {
		return nil, NewInternalError(errors.New("Parameters value not a pointer to struct"))
	}
	paramsType = paramsType.Elem()

	if paramsType.Kind() != reflect.Struct {
		return nil, NewInternalError(errors.New("Parameters value not a pointer to struct"))
	}

	return paramsType, nil
}

func generateRouteParams(ctx context.Context, route *restful.Route, paramsType reflect.Type) (*RouteParams, error) {
	routeParams := &RouteParams{
		Type:   paramsType,
		Fields: nil,
	}

	for i := 0; i < paramsType.NumField(); i++ {
		routeParam := NewRouteParam(ctx, route, paramsType.Field(i))
		if routeParam == nil {
			continue
		}
		routeParams.Fields = append(routeParams.Fields, routeParam)
	}

	return routeParams, nil
}
