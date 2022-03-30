package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"encoding/json"
	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/refl"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	PortTypeRequest  = StructTagRequest
	PortTypeResponse = StructTagResponse

	FieldShapePrimitive = "primitive" // Input, Output
	FieldShapeArray     = "array"     // Input, Output
	FieldShapeObject    = "object"    // Input, Output
	FieldShapeFile      = "file"      // Input
	FieldShapeFileArray = "filearray" // Input
	FieldShapeReader    = "reader"    // Input (Body), Output (Body)
	FieldShapeUnknown   = "unknown"   // Ignored
)

// EndpointPort defines the inputs/outputs for an Entry Point
type EndpointPort struct {
	Type   reflect.Type
	Fields EndpointPortFields
}

func (p EndpointPort) NewStruct() interface{} {
	return reflect.New(p.Type).Interface()
}

func (p EndpointPort) Validate(inputs interface{}) error {
	if inputs == nil {
		return nil
	}

	isv := reflect.ValueOf(inputs)
	for isv.Kind() == reflect.Ptr {
		isv = isv.Elem()
	}

	errs := make(ValidationErrors)
	for _, portField := range p.Fields {
		ifv := isv.FieldByIndex(portField.Field.Index)
		if err := validate.ValidateValue(ifv); err != nil {
			errs[portField.Name] = NewValidationFailure("").Apply(err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

type EndpointPortFields []EndpointPortField

func (f EndpointPortFields) First(pred func(p EndpointPortField) bool) (EndpointPortField, bool) {
	for _, p := range f {
		if pred(p) {
			return p, true
		}
	}
	return EndpointPortField{}, false
}

func (f EndpointPortFields) All(pred func(p EndpointPortField) bool) []EndpointPortField {
	var results []EndpointPortField
	for _, p := range f {
		if pred(p) {
			results = append(results, p)
		}
	}
	return results
}

func (f EndpointPortFields) Headers() EndpointPortFields {
	return f.All(func(p EndpointPortField) bool { return p.IsHeader() && !p.IsError() })
}

func (f EndpointPortFields) ErrorHeaders() EndpointPortFields {
	return f.All(func(p EndpointPortField) bool { return p.IsHeader() && p.IsError() })
}

func (f EndpointPortFields) Code() (EndpointPortField, bool) {
	return f.First(func(f EndpointPortField) bool { return f.IsCode() })
}

func (f EndpointPortFields) Body() (EndpointPortField, bool) {
	return f.First(func(f EndpointPortField) bool { return f.IsBody() && !f.IsError() })
}

func (f EndpointPortFields) ErrorBody() (EndpointPortField, bool) {
	return f.First(func(f EndpointPortField) bool { return f.IsBody() && f.IsError() })
}

// EndpointPortField defines a single input/output field for an Entry point
type EndpointPortField struct {
	Name            string
	Field           reflect.StructField
	Options         map[string]string
	SanitizeOptions sanitize.Options
	Optional        bool
	In              string
	Shape           string
	SchemaOrRef     *openapi3.SchemaOrRef
}

func (r EndpointPortField) FieldValue(portStruct interface{}) reflect.Value {
	fieldIndex := r.Field.Index
	ov := reflect.ValueOf(portStruct)
	return ov.FieldByIndex(fieldIndex)
}

func (r EndpointPortField) PrimitiveValue(portStruct interface{}) types.OptionalString {
	fv := r.FieldValue(portStruct)
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return types.OptionalString{}
		}
		fv = fv.Elem()
	}

	result, err := cast.ToStringE(fv.Interface())
	if err != nil {
		logger.WithError(err).Errorf("Could not coerce %T to primitive", fv.Interface())
		return types.OptionalString{}
	} else {
		return types.OptionalString{
			Value: &result,
		}
	}
}

func (r EndpointPortField) ArrayValue(portStruct interface{}) []string {
	fv := r.FieldValue(portStruct)
	result, err := cast.ToStringSliceE(fv.Interface())
	if err != nil {
		logger.WithError(err).Errorf("Could not coerce %T to array", fv.Interface())
		return nil
	} else {
		return result
	}
}

func (r EndpointPortField) ObjectValue(portStruct interface{}) types.Pojo {
	fv := r.FieldValue(portStruct)

	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}

	var result types.Pojo
	var err error

	if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Map {
		data, _ := json.Marshal(fv.Interface())
		_ = json.Unmarshal(data, &result)
	} else {
		result, err = cast.ToStringMapE(fv.Interface())
	}

	if err != nil {
		logger.WithError(err).Errorf("Could not coerce %T to object", fv.Interface())
		return nil
	} else {
		return result
	}
}

func (r EndpointPortField) BoolOption(name string) bool {
	return r.Options[name] == "true"
}

func (r EndpointPortField) Schema() (schema *openapi3.Schema) {
	if r.SchemaOrRef.SchemaReference != nil {
		refName := SchemaRefName(r.SchemaOrRef.SchemaReference)

		var ok bool
		if schema, ok = DeepLookupSchema(refName); !ok {
			return nil
		}

	} else if r.SchemaOrRef.Schema != nil {
		schema = r.SchemaOrRef.Schema
	} else {
		return nil
	}

	return schema
}

func (r EndpointPortField) DefaultValue() *interface{} {
	schema := r.Schema()
	if schema == nil {
		return nil
	}

	if schema.Default != nil && schema.Type != nil {
		castValue := r.CastWide(*schema.Type, *schema.Default)
		return &castValue
	}

	return nil
}

func (r EndpointPortField) CastWide(schemaType openapi3.SchemaType, val interface{}) interface{} {
	switch schemaType {
	case openapi3.SchemaTypeBoolean:
		return cast.ToBool(val)
	case openapi3.SchemaTypeString:
		return cast.ToString(val)
	case openapi3.SchemaTypeInteger:
		return cast.ToInt(val)
	case openapi3.SchemaTypeNumber:
		return cast.ToFloat64(val)
	default:
		return nil
	}
}

func (r EndpointPortField) IsBody() bool {
	return r.In == "body"
}

func (r EndpointPortField) IsHeader() bool {
	return r.In == "header"
}

func (r EndpointPortField) IsCode() bool {
	return r.In == "code"
}

func (r EndpointPortField) IsForm() bool {
	return r.In == "form"
}

func (r EndpointPortField) IsError() bool {
	return r.BoolOption("error")
}

func (r EndpointPortField) Tags() reflect.StructTag {
	tagBuilder := strings.Builder{}
	for key, value := range r.Options {
		if key == "optional" {
			key = "required"
			if value == "true" {
				value = "false"
			} else {
				value = "true"
			}
		}

		tagBuilder.WriteString(key)
		tagBuilder.WriteString(":\"")
		tagBuilder.WriteString(value)
		tagBuilder.WriteString("\" ")
	}

	tags, err := structtag.Parse(string(r.Field.Tag))
	if err != nil {
		logger.WithError(err).Errorf("Failed to parse struct tag for field %q", r.Field.Name)
		return ""
	}
	for _, tag := range tags.Tags() {
		if tag.Key == PortTypeRequest || tag.Key == PortTypeResponse {
			continue
		}

		tagBuilder.WriteRune(' ')
		tagBuilder.WriteString(tag.String())
	}

	return reflect.StructTag(tagBuilder.String())
}

func (r EndpointPortField) Parameter() EndpointRequestParameter {
	parameter := NewEndpointRequestParameter(r.Name, r.In).
		WithRequired(!r.Optional)

	tagValue := r.Tags()

	if err := refl.PopulateFieldsFromTags(&parameter, tagValue); err != nil {
		logger.WithError(err).Errorf("Failed to populate parameter fields from `req` tag for port field %q", parameter.Name)
	}

	if example, ok := r.Options["example"]; ok {
		parameter = parameter.WithExample(example)
	}

	if parameter.Style == nil {
		switch r.In {
		case "header":
			parameter = parameter.WithStyle("simple")
		case "path":
			parameter = parameter.WithStyle("simple")
		case "query":
			parameter = parameter.WithStyle("form")
		case "form":
			parameter = parameter.WithStyle("form")
		case "cookie":
			parameter = parameter.WithStyle("form")
		}
	}

	if parameter.Explode == nil {
		switch r.In {
		case "form":
			parameter = parameter.WithExplode(true)
		default:
			parameter = parameter.WithExplode(false)
		}
	}

	if r.SchemaOrRef != nil {
		parameter = parameter.WithSchema(*r.SchemaOrRef)
	}

	parameter.PortField = &r

	return parameter
}

func (r EndpointPortField) RequestBody() EndpointRequestBody {
	var result = EndpointRequestBody{
		Required: r.Field.Type.Kind() != reflect.Ptr,
		Mime:     MIME_JSON,
	}

	requestTagValue := r.Tags()

	if err := refl.PopulateFieldsFromTags(&result, requestTagValue); err != nil {
		logger.WithError(err).Errorf("Failed to populate request body fields from tags for arg %q", r.Name)
	}

	bodyInstance := types.Instantiate(r.Field.Type)
	schemaOrRef, err := Reflect(bodyInstance)
	if err != nil {
		logger.WithError(err).Errorf("Failed to reflect body schema arg %q", r.Name)
		return result
	}

	result.Schema = schemaOrRef
	return result
}

func (r EndpointPortField) RequestBodyFormField() EndpointRequestBodyFormField {
	return EndpointRequestBodyFormField{
		Name:     r.Name,
		Required: !r.Optional,
		Schema:   r.SchemaOrRef,
	}
}

func (r EndpointPortField) ResponseHeader() EndpointResponseHeader {
	header := NewEndpointResponseHeader().
		WithRequired(!r.Optional)

	tagValue := r.Tags()

	if err := refl.PopulateFieldsFromTags(&header, tagValue); err != nil {
		logger.WithError(err).Error("Failed to populate header from struct tags")
	}

	if example, ok := r.Options["example"]; ok {
		header = header.WithExample(example)
	}

	if header.Explode == nil {
		header = header.WithExplode(false)
	}

	if r.SchemaOrRef != nil {
		header = header.WithSchema(*r.SchemaOrRef)
	}

	header.PortField = &r

	return header
}

// loadEnum loads enum from interface or field tag: json array or comma-separated string.
func loadEnum(fieldTag reflect.StructTag, fieldVal interface{}) []interface{} {
	var items []interface{}

	if e, isEnumer := fieldVal.(jsonschema.Enum); isEnumer {
		items = e.Enum()
	}

	if enumTag := fieldTag.Get("enum"); enumTag != "" {
		var e []interface{}

		err := json.Unmarshal([]byte(enumTag), &e)
		if err != nil {
			es := strings.Split(enumTag, ",")
			e = make([]interface{}, len(es))

			for i, s := range es {
				e[i] = s
			}
		}

		items = e
	}

	return items
}

func NewPortField(portType string, field reflect.StructField) (EndpointPortField, bool) {
	var optional bool

	// Skip fields with invalid tags
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		logger.WithError(err).Errorf("Failed to parse field %q tag", field.Name)
		return EndpointPortField{}, false
	}

	// Skip fields with no port tag
	portTag, err := tags.Get(portType)
	if err != nil {
		return EndpointPortField{}, false
	}

	var result = EndpointPortField{
		Field:   field,
		Options: make(map[string]string),
	}
	optional, result.Shape = getPortFieldShape(field.Type)

	// Skip ignored fields
	if portTag.Name == "-" {
		return EndpointPortField{}, false
	}

	// Parse the options
	sourceParts := strings.SplitN(portTag.Name, "=", 2)
	result.In = sourceParts[0]
	if len(sourceParts) == 2 {
		result.Name = sourceParts[1]
	} else {
		switch result.In {
		case "header":
			result.Name = strcase.ToKebab(field.Name)
		case "path":
			result.Name = strcase.ToLowerCamel(field.Name)
			result.Optional = false
		case "query":
			result.Name = strcase.ToLowerCamel(field.Name)
		case "form":
			result.Name = strcase.ToLowerCamel(field.Name)
		case "cookie":
			result.Name = strcase.ToLowerCamel(field.Name)
		case "body":
			result.Name = "body"
		case "code":
			result.Name = "code"
		}
	}

	for _, option := range portTag.Options {
		optionParts := strings.SplitN(option, "=", 2)

		var value = "true"
		if len(optionParts) == 2 {
			value = optionParts[1]
		}

		if optionParts[0] == "optional" {
			optional = value == "true"
			optionParts[0] = "required"
			if len(optionParts) == 2 {
				optionParts[1] = strconv.FormatBool(!optional)
			} else {
				optionParts = append(optionParts, strconv.FormatBool(!optional))
			}
			value = optionParts[1]
		} else if optionParts[0] == "required" {
			optional = value != "true"
		}

		result.Options[optionParts[0]] = value
	}

	sanTag := field.Tag.Get(sanitizeTag)
	result.SanitizeOptions = sanitize.NewOptions(sanTag)
	if sanTag != "" {
		if _, exists := result.Options["san"]; !exists {
			result.Options["san"] = "true"
		}
	}

	requiredTag := field.Tag.Get("required")
	if requiredTag != "" {
		optional = requiredTag != "true"
	}

	optionalTag := field.Tag.Get("optional")
	if optionalTag != "" {
		optional = optionalTag == "true"
	}

	nullableTag := field.Tag.Get("nullable")
	if nullableTag != "" {
		optional = nullableTag == "true"
	}

	result.Optional = optional

	fieldInstance := types.Instantiate(result.Field.Type)

	schemaOrRef, err := Reflect(fieldInstance)
	if err != nil {
		logger.WithError(err).Errorf("Failed to calculate schema for field %q", result.Name)
	} else {
		result.SchemaOrRef = schemaOrRef
	}

	if schemaOrRef.Schema != nil {
		if err = populateSchemaFromTags(fieldInstance, field.Tag, schemaOrRef.Schema); err != nil {
			logger.WithError(err).Errorf("Failed to populate openapi schema fields from tags for port field %q", field.Name)
		}

	} else {
		tagSchema := new(openapi3.Schema)

		if err = populateSchemaFromTags(fieldInstance, field.Tag, tagSchema); err != nil {
			logger.WithError(err).Errorf("Failed to populate openapi schema fields from tags for port field %q", field.Name)
		}

		if schemaOrRef != nil {
			allOf := CombineSchemas(*schemaOrRef, NewSchemaOrRef(tagSchema))
			schemaOrRef = &allOf
		} else {
			schemaOrRef = NewSchemaOrRefPtr(tagSchema)
		}

		result.SchemaOrRef = schemaOrRef
	}

	if result.BoolOption("envelope") && result.Shape == FieldShapeReader {
		logger.WithError(err).Errorf("Enveloped reader not supported for field %q", result.Name)
		result.Options["envelope"] = "false"
	}

	return result, true
}

func populateSchemaFromTags(fieldInstance interface{}, tag reflect.StructTag, schema *openapi3.Schema) error {
	if err := refl.PopulateFieldsFromTags(schema, tag); err != nil {
		return err
	}

	if items := loadEnum(tag, fieldInstance); items != nil {
		schema.Enum = items
	}

	return nil
}

var typesTimeType = reflect.TypeOf(types.Time{})
var typesUuidType = reflect.TypeOf(types.UUID{})
var textUnmarshalerInstance types.TextUnmarshaler
var textUnmarshalerType = reflect.TypeOf(&textUnmarshalerInstance).Elem()
var multipartFileHeaderInstance multipart.FileHeader
var multipartFileHeaderType = reflect.TypeOf(multipartFileHeaderInstance)
var byteType = reflect.TypeOf(byte(0))
var runeType = reflect.TypeOf(rune(0))
var byteSliceType = reflect.TypeOf([]byte{})
var b64BytesType = reflect.TypeOf(types.Base64Bytes{})
var ioReadCloserInstance io.ReadCloser
var ioReadCloserType = reflect.TypeOf(&ioReadCloserInstance).Elem()

func getPortFieldShape(fieldType reflect.Type) (optional bool, shape string) {
	isPtr := false
	if fieldType.Kind() == reflect.Ptr {
		optional = true
		isPtr = true
		fieldType = fieldType.Elem()
	} else if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Map {
		optional = true
	}

	fieldPtrType := reflect.PtrTo(fieldType)

	// Custom type handling
	if fieldType == multipartFileHeaderType {
		return false, FieldShapeFile
	} else if fieldType == byteSliceType || fieldType == typesUuidType || fieldType == b64BytesType {
		return isPtr, FieldShapePrimitive
	} else if fieldType.Implements(textUnmarshalerType) || fieldPtrType.Implements(textUnmarshalerType) {
		return optional, FieldShapePrimitive
	} else if fieldType.Implements(ioReadCloserType) || fieldPtrType.Implements(ioReadCloserType) {
		return optional, FieldShapeReader
	}

	// General type handling
	switch fieldType.Kind() {
	case reflect.Slice, reflect.Array:
		elemType := fieldType.Elem()
		if elemType == multipartFileHeaderType {
			return false, FieldShapeFileArray
		} else if elemType == byteType || elemType == runeType {
			shape = FieldShapePrimitive
		} else {
			shape = FieldShapeArray
		}
	case reflect.Map, reflect.Struct:
		shape = FieldShapeObject
	case reflect.Float64, reflect.Float32:
		shape = FieldShapePrimitive
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		shape = FieldShapePrimitive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		shape = FieldShapePrimitive
	case reflect.String:
		shape = FieldShapePrimitive
	case reflect.Bool:
		shape = FieldShapePrimitive
	default:
		logger.Warnf("Cannot marshal field type %t", fieldType)
		shape = FieldShapeUnknown
	}

	return optional, shape
}

var portsIndex = make(map[reflect.Type]EndpointPort)
var portsMtx sync.Mutex

func parsePort(portType string, requestArgsType reflect.Type) (EndpointPort, error) {
	requestArgs := EndpointPort{
		Type: requestArgsType,
	}

	for i := 0; i < requestArgsType.NumField(); i++ {
		requestArgField := requestArgsType.Field(i)
		requestArgument, ok := NewPortField(portType, requestArgField)
		if ok {
			requestArgs.Fields = append(requestArgs.Fields, requestArgument)
		}
	}

	return requestArgs, nil
}

func NewEndpointPort(portType string, portStruct interface{}) (result EndpointPort, err error) {
	portsMtx.Lock()
	defer portsMtx.Unlock()

	if portStruct == nil {
		return EndpointPort{
			Type:   reflect.TypeOf(struct{}{}),
			Fields: nil,
		}, nil
	}

	portStructType := reflect.TypeOf(portStruct)
	if portStructType.Kind() != reflect.Struct {
		return result, errors.New("Port value not a struct")
	}

	var ok bool
	if result, ok = portsIndex[portStructType]; ok {
		return
	}

	result, err = parsePort(portType, portStructType)
	if err != nil {
		return
	}

	portsIndex[portStructType] = result

	return result, nil
}
