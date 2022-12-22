// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/refl"
	"reflect"
)

type PortValidatorFunction func(p interface{}) (err error)

type EndpointRequest struct {
	Port        *ops.Port
	Description string
	Parameters  []EndpointRequestParameter
	Validator   PortValidatorFunction
	Body        EndpointRequestBody
}

func (r EndpointRequest) WithParameter(p EndpointRequestParameter) EndpointRequest {
	for i, par := range r.Parameters {
		if par.Name == p.Name {
			r.Parameters[i] = p
			return r
		}
	}

	r.Parameters = append(r.Parameters, p)
	return r
}

func (r EndpointRequest) WithValidator(fn PortValidatorFunction) EndpointRequest {
	r.Validator = fn
	return r
}

func (r EndpointRequest) WithBody(body EndpointRequestBody) EndpointRequest {
	r.Body = body
	return r
}

func (r EndpointRequest) WithInputs(portStruct interface{}) EndpointRequest {
	result := r

	var portStructType reflect.Type
	if rt, ok := portStruct.(reflect.Type); ok {
		portStructType = rt
	} else {
		portStructType = refl.DeepIndirect(reflect.TypeOf(portStruct))
	}

	port, err := PortReflector{}.ReflectInputPort(portStructType)
	if err != nil {
		logger.WithError(err).Error("Failed to reflect request input port")
		return result
	}

	result.Port = port

	for _, field := range port.Fields {
		result = result.withPortField(field)
	}

	return result
}

func (r EndpointRequest) withPortField(field *ops.PortField) EndpointRequest {
	result := r
	if PortFieldIsBody(field) {
		result = result.WithBody(EndpointRequestBodyFromPortField(field))
	} else if PortFieldIsForm(field) {
		result.Body = result.Body.WithFormField(EndpointRequestBodyFormFieldFromPortField(field))
	} else {
		result = result.WithParameter(EndpointRequestParameterFromPortField(field))
	}
	return result
}

func (r EndpointRequest) PatchParameter(name string, fn func(p EndpointRequestParameter) EndpointRequestParameter) EndpointRequest {
	for i, p := range r.Parameters {
		if p.Name == name {
			r.Parameters[i] = fn(p)
			break
		}
	}
	return r
}

func (r EndpointRequest) Consumes() []string {
	return []string{r.Body.Mime}
}

func (r EndpointRequest) HasBody() bool {
	return r.Body.Mime != ""
}

func NewEndpointRequest() EndpointRequest {
	return EndpointRequest{}
}

type EndpointRequestBodyFormField struct {
	Name            string
	PortField       *ops.PortField
	Required        bool
	Deprecated      *bool
	AllowEmptyValue *bool
	Style           *string
	Explode         *bool
	AllowReserved   *bool
	Type            *string
	Format          *string
	Multi           *bool
	Enum            []interface{}
	Default         types.Optional[interface{}]
	Payload         types.Optional[interface{}]
	Example         types.Optional[interface{}]
}

func (f EndpointRequestBodyFormField) Parameter() EndpointRequestParameter {
	return EndpointRequestParameter{
		Name:      f.Name,
		Required:  types.NewBoolPtr(f.Required),
		PortField: f.PortField,
		Style:     types.NewStringPtr(""),
		In:        "form",
		Explode:   types.NewBoolPtr(true),
		Type:      f.Type,
		Format:    f.Format,
		Payload:   f.Payload,
		Example:   f.Example,
	}
}

func (f EndpointRequestBodyFormField) WithType(dataType string) EndpointRequestBodyFormField {
	f.Type = &dataType
	return f
}

func (f EndpointRequestBodyFormField) WithFormat(format string) EndpointRequestBodyFormField {
	f.Format = &format
	return f
}

func (f EndpointRequestBodyFormField) WithMulti(multi bool) EndpointRequestBodyFormField {
	f.Multi = &multi
	return f
}

func (f EndpointRequestBodyFormField) WithEnum(enum []interface{}) EndpointRequestBodyFormField {
	f.Enum = enum
	return f
}

func (f EndpointRequestBodyFormField) WithDefault(value interface{}) EndpointRequestBodyFormField {
	f.Default = types.OptionalOf(value)
	return f
}

func (f EndpointRequestBodyFormField) WithStyle(style string) EndpointRequestBodyFormField {
	f.Style = &style
	return f
}

func (f EndpointRequestBodyFormField) WithAllowEmptyValue(allowEmptyValue bool) EndpointRequestBodyFormField {
	f.AllowEmptyValue = &allowEmptyValue
	return f
}

func (f EndpointRequestBodyFormField) WithExplode(explode bool) EndpointRequestBodyFormField {
	f.Explode = &explode
	return f
}

func (f EndpointRequestBodyFormField) WithAllowReserved(allowReserved bool) EndpointRequestBodyFormField {
	f.AllowReserved = &allowReserved
	return f
}

func EndpointRequestBodyFormFieldFromPortField(pf *ops.PortField) EndpointRequestBodyFormField {
	var result = EndpointRequestBodyFormField{
		Name:      pf.Peer,
		Required:  !pf.Optional,
		PortField: pf,
	}

	requestTagValue := pf.Tags()

	if err := refl.PopulateFieldsFromTags(&result, requestTagValue); err != nil {
		logger.WithError(err).Errorf("Failed to populate request body form field from tags for arg %q", pf.Name)
	}

	return result
}

func EndpointRequestBodyFormFieldFromSwaggerParam(sp SwaggerParam, options map[string]string) EndpointRequestBodyFormField {
	ff := EndpointRequestBodyFormField{
		Name:     sp.Name,
		Required: sp.Required,
		Type:     &sp.DataType,
	}

	if sp.DataFormat != "" {
		ff.Format = types.NewStringPtr(sp.DataFormat)
	}

	if sp.DataFormat != "" {
		ff = ff.WithFormat(sp.DataFormat)
	}

	if sp.DefaultValue != "" {
		ff = ff.WithDefault(sp.DefaultValue)
	}

	// CollectionFormat conversion
	if options["csv"] == "true" && sp.CollectionFormat == "csv" {
		switch sp.In {
		case FieldGroupHttpQuery:
			ff = ff.WithStyle("form").WithExplode(false)
		case FieldGroupHttpPath, FieldGroupHttpHeader:
			ff = ff.WithStyle("simple").WithExplode(false)
		}
	} else if options["multi"] == "true" && sp.CollectionFormat == "multi" {
		ff = ff.WithExplode(true)
	}

	if sp.AllowMultiple {
		// Wrap the swagger2 schema in an array
		ff = ff.WithMulti(true)
	}

	if len(sp.AllowableValues) > 0 {
		var enum []interface{}
		for key := range sp.AllowableValues {
			enum = append(enum, key)
		}
		ff = ff.WithEnum(enum)
	}

	return ff
}

type EndpointRequestBody struct {
	Description string
	Required    bool
	Mime        string
	Payload     types.Optional[interface{}]
	Example     types.Optional[interface{}]
	FormFields  []EndpointRequestBodyFormField
	ops.Documentors[EndpointRequestBody]
	PortField *ops.PortField
}

func (b EndpointRequestBody) WithDocumentor(doc ...ops.Documentor[EndpointRequestBody]) EndpointRequestBody {
	b.Documentors = b.Documentors.WithDocumentor(doc...)
	return b
}

func (b EndpointRequestBody) WithFormField(field EndpointRequestBodyFormField) EndpointRequestBody {
	b.Mime = MediaTypeMultipartForm
	b.Required = true
	b.FormFields = append(b.FormFields, field)
	return b
}

func (b EndpointRequestBody) HasFormField() bool {
	return len(b.FormFields) > 0
}

func EndpointRequestBodyFromPortField(pf *ops.PortField) EndpointRequestBody {
	var result = EndpointRequestBody{
		Required:  !pf.Optional,
		Mime:      MediaTypeJson,
		Payload:   types.OptionalOf(reflect.New(pf.Type.Type).Interface()),
		PortField: pf,
	}

	requestTagValue := pf.Tags()

	if err := refl.PopulateFieldsFromTags(&result, requestTagValue); err != nil {
		logger.WithError(err).Errorf("Failed to populate request body fields from tags for arg %q", pf.Name)
	}

	return result
}

// Type Precedence:
// - PortField: primary source (generated from port structure)
// - Payload: secondary source (custom provided by developer in endpoint definition)
// - Type/Format: tertiary source (from swagger 2.0)

type EndpointRequestParameter struct {
	Name            string // Required.
	In              string // Required.
	Description     *string
	Required        *bool
	Deprecated      *bool
	AllowEmptyValue *bool
	Style           *string
	Explode         *bool
	AllowReserved   *bool
	Example         types.Optional[interface{}]
	Reference       *string
	Format          *string
	Type            *string
	Multi           *bool
	Enum            []interface{}
	Default         types.Optional[interface{}]
	Payload         types.Optional[interface{}]
	PortField       *ops.PortField
	ops.Documentors[EndpointRequestParameter]
}

func (p EndpointRequestParameter) WithDocumentor(doc ...ops.Documentor[EndpointRequestParameter]) EndpointRequestParameter {
	p.Documentors = p.Documentors.WithDocumentor(doc...)
	return p
}

func (p EndpointRequestParameter) WithDescription(description string) EndpointRequestParameter {
	p.Description = &description
	return p
}

func (p EndpointRequestParameter) WithRequired(required bool) EndpointRequestParameter {
	p.Required = &required
	return p
}

func (p EndpointRequestParameter) WithDeprecated(deprecated bool) EndpointRequestParameter {
	p.Deprecated = &deprecated
	return p
}

func (p EndpointRequestParameter) WithStyle(style string) EndpointRequestParameter {
	p.Style = &style
	return p
}

func (p EndpointRequestParameter) WithAllowEmptyValue(allowEmptyValue bool) EndpointRequestParameter {
	p.AllowEmptyValue = &allowEmptyValue
	return p
}

func (p EndpointRequestParameter) WithExplode(explode bool) EndpointRequestParameter {
	p.Explode = &explode
	return p
}

func (p EndpointRequestParameter) WithAllowReserved(allowReserved bool) EndpointRequestParameter {
	p.AllowReserved = &allowReserved
	return p
}

func (p EndpointRequestParameter) WithExample(example interface{}) EndpointRequestParameter {
	if example != nil {
		p.Example = types.OptionalOf(example)
	} else {
		p.Example = types.OptionalEmpty[interface{}]()
	}
	return p
}

func (p EndpointRequestParameter) WithReference(reference string) EndpointRequestParameter {
	p.Reference = &reference
	return p
}

func (p EndpointRequestParameter) WithType(dataType string) EndpointRequestParameter {
	p.Type = &dataType
	return p
}

func (p EndpointRequestParameter) WithFormat(format string) EndpointRequestParameter {
	p.Format = &format
	return p
}

func (p EndpointRequestParameter) WithMulti(multi bool) EndpointRequestParameter {
	p.Multi = &multi
	return p
}

func (p EndpointRequestParameter) WithEnum(enum []interface{}) EndpointRequestParameter {
	p.Enum = enum
	return p
}

func (p EndpointRequestParameter) WithDefault(value interface{}) EndpointRequestParameter {
	p.Default = types.OptionalOf(value)
	return p
}

func (p EndpointRequestParameter) WithPayload(payload interface{}) EndpointRequestParameter {
	if payload == nil {
		p.Payload = types.OptionalEmpty[interface{}]()
	} else {
		p.Payload = types.OptionalOf(payload)
	}
	return p
}

func NewEndpointRequestParameter(name, in string) EndpointRequestParameter {
	return EndpointRequestParameter{
		Name: name,
		In:   in,
	}
}

func EndpointRequestParameterFromPortField(pf *ops.PortField) EndpointRequestParameter {
	result := NewEndpointRequestParameter(pf.Peer, pf.Group)
	result.PortField = pf

	if pf.Group != FieldGroupHttpPath && !pf.Optional {
		result.Required = types.NewBoolPtr(!pf.Optional)
	}

	requestTagValue := pf.Tags()

	if err := refl.PopulateFieldsFromTags(&result, requestTagValue); err != nil {
		logger.WithError(err).Errorf("Failed to populate request parameter fields from tags for arg %q", pf.Name)
	}

	if example, ok := pf.Options[exampleTag]; ok {
		result = result.WithExample(example)
	}

	if result.Style == nil {
		switch pf.Group {
		case FieldGroupHttpHeader:
			result = result.WithStyle("simple")
		case FieldGroupHttpPath:
			result = result.WithStyle("simple")
		case FieldGroupHttpQuery:
			result = result.WithStyle("form")
		case FieldGroupHttpForm:
			result = result.WithStyle("form")
		case FieldGroupHttpCookie:
			result = result.WithStyle("form")
		}
	}

	if result.Explode == nil {
		switch pf.Group {
		case FieldGroupHttpForm:
			result = result.WithExplode(true)
		default:
			result = result.WithExplode(false)
		}
	}

	return result
}

func EndpointRequestParameterFromSwaggerParam(paramData SwaggerParam, options map[string]string) EndpointRequestParameter {
	p := NewEndpointRequestParameter(paramData.Name, paramData.In).
		WithDescription(paramData.Description).
		WithRequired(paramData.Required).
		WithType(paramData.DataType)

	if paramData.DataFormat != "" {
		p = p.WithFormat(paramData.DataFormat)
	}

	if paramData.DefaultValue != "" {
		p = p.WithDefault(paramData.DefaultValue)
	}

	// CollectionFormat conversion
	if options["csv"] == "true" && paramData.CollectionFormat == "csv" {
		switch paramData.In {
		case FieldGroupHttpQuery:
			p = p.WithStyle("form").WithExplode(false)
		case FieldGroupHttpPath, FieldGroupHttpHeader:
			p = p.WithStyle("simple").WithExplode(false)
		}
	} else if options["multi"] == "true" && paramData.CollectionFormat == "multi" {
		p = p.WithExplode(true)
	}

	if paramData.AllowMultiple {
		// Wrap the swagger2 schema in an array
		p = p.WithMulti(true)
	}

	if len(paramData.AllowableValues) > 0 {
		var enum []interface{}
		for key := range paramData.AllowableValues {
			enum = append(enum, key)
		}
		p = p.WithEnum(enum)
	}

	return p
}

const (
	ParameterInPath   = "path"
	ParameterInQuery  = "query"
	ParameterInHeader = "header"
	ParameterInCookie = "cookie"
)

func PathParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, ParameterInPath).
		WithDescription(description).
		WithRequired(true)
}

func QueryParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, ParameterInQuery).
		WithDescription(description)
}

func HeaderParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, ParameterInHeader).
		WithDescription(description)
}

func CookieParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, ParameterInCookie).
		WithDescription(description)
}
