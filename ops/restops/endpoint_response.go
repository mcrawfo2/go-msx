// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/refl"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type EndpointResponse struct {
	Codes    EndpointResponseCodes
	Envelope bool
	Success  EndpointResponseContent
	Error    EndpointResponseContent
	Port     *ops.Port
	ops.Documentors[EndpointResponse]
}

func (r EndpointResponse) WithDocumentor(doc ...ops.Documentor[EndpointResponse]) EndpointResponse {
	r.Documentors = r.Documentors.WithDocumentor(doc...)
	return r
}

func (r EndpointResponse) WithEnvelope(envelope bool) EndpointResponse {
	r.Envelope = envelope
	if envelope {
		// Error is rendered as part of the envelope
		r.Error.Payload = types.OptionalEmpty[interface{}]()
	}
	return r
}

func (r EndpointResponse) WithResponseCodes(codes EndpointResponseCodes) EndpointResponse {
	r.Codes = codes
	return r
}

func (r EndpointResponse) WithSuccessPayload(body interface{}) EndpointResponse {
	if r.Success.Mime == "" {
		r.Success.Mime = MediaTypeJson
		r.Error.Mime = MediaTypeJson
	}
	r.Success = r.Success.WithPayload(body)
	return r
}

func (r EndpointResponse) WithSuccessHeader(name string, header EndpointResponseHeader) EndpointResponse {
	r.Success = r.Success.WithHeader(name, header)
	return r
}

func (r EndpointResponse) WithErrorPayload(body interface{}) EndpointResponse {
	r.Error = r.Error.WithPayload(body)
	return r
}

func (r EndpointResponse) WithErrorHeader(name string, header EndpointResponseHeader) EndpointResponse {
	r.Error = r.Error.WithHeader(name, header)
	return r
}

func (r EndpointResponse) WithMime(mime string) EndpointResponse {
	r.Success = r.Success.WithMime(mime)
	r.Error = r.Error.WithMime(mime)
	return r
}

func (r EndpointResponse) WithSuccess(success EndpointResponseContent) EndpointResponse {
	r.Success = success
	return r
}

func (r EndpointResponse) WithError(error EndpointResponseContent) EndpointResponse {
	r.Error = error
	return r
}

func (r EndpointResponse) WithHeader(name string, header EndpointResponseHeader) EndpointResponse {
	return r.
		WithSuccessHeader(name, header).
		WithErrorHeader(name, header)
}

func (r EndpointResponse) HasBody() bool {
	return r.Success.Mime != "" || r.Error.Mime != ""
}

func (r EndpointResponse) Produces() []string {
	var mimeTypes = make(types.StringSet)
	mimeTypes.AddAll(r.Success.Mime, r.Error.Mime)
	if _, ok := mimeTypes[""]; ok {
		delete(mimeTypes, "")
	}
	return mimeTypes.Values()
}

func (r EndpointResponse) withEnumCodes(enum string) EndpointResponse {
	vals := strings.Split(enum, ",")
	codes := EndpointResponseCodes{}
	for _, v := range vals {
		iv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.WithError(err).Errorf("Invalid response code %q", v)
			continue
		}
		if iv < 200 || iv > 599 {
			logger.Errorf("Invalid response code %d", iv)
			continue
		}
		if iv <= 399 {
			codes.Success = append(codes.Success, int(iv))
		} else {
			codes.Error = append(codes.Error, int(iv))
		}
	}
	return r.WithResponseCodes(codes)
}

func (r EndpointResponse) WithOutputs(portStruct interface{}) EndpointResponse {
	result := r

	var portStructType reflect.Type
	if rt, ok := portStruct.(reflect.Type); ok {
		portStructType = rt
	} else {
		portStructType = refl.DeepIndirect(reflect.TypeOf(portStruct))
	}

	port, err := PortReflector{}.ReflectOutputPort(portStructType)
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

func (r EndpointResponse) withPortField(field *ops.PortField) EndpointResponse {
	result := r

	if PortFieldIsBody(field) {
		if field.Options["envelope"] == "true" {
			result = result.WithEnvelope(true)
		}

		if field.Options["error"] == "true" {
			result.Error = result.Error.withPortFieldBody(field)
		} else {
			result.Success = result.Success.withPortFieldBody(field)
		}

	} else if PortFieldIsCode(field) {

		tags := field.Tags()
		if enumTag, ok := tags.Lookup("enum"); ok {
			result = result.withEnumCodes(enumTag)
		}

	} else if PortFieldIsHeader(field) {
		fieldHeader := EndpointResponseHeaderFromPortField(field)

		errorOnly, _ := field.BoolOption("error")
		if !errorOnly {
			result = result.WithSuccessHeader(field.Peer, fieldHeader)
		}

		successOnly, _ := field.BoolOption("success")
		if !successOnly {
			result = result.WithErrorHeader(field.Peer, fieldHeader)
		}
	} else if PortFieldIsPaging(field) {
		result.Success = result.Success.withPortFieldPaging(field)
	}

	return result
}

func NewEndpointResponse() EndpointResponse {
	return EndpointResponse{}
}

type EndpointResponseCodes struct {
	Success []int
	Error   []int
}

func (c EndpointResponseCodes) DefaultCode() int {
	if len(c.Success) == 0 {
		return http.StatusOK
	}
	return c.Success[0]
}

var (
	StandardErrorCodes = []int{
		http.StatusBadRequest,   // Some parameter issue
		http.StatusUnauthorized, // No tenant access
		http.StatusForbidden,    // No authentication token
		http.StatusNotFound,     // Target resource does not exist
		http.StatusConflict,     // State issue (eg. already exists)
	}
	ListResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized},
	}
	GetResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusNotFound},
	}
	CreateResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusCreated},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusConflict},
	}
	UpdateResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusNotFound},
	}
	DeleteResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   StandardErrorCodes,
	}
	AcceptResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusAccepted},
		Error:   StandardErrorCodes,
	}
	NoContentResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusNoContent},
		Error:   StandardErrorCodes,
	}
	UnknownVerbResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden},
	}
)

func DefaultResponseCodes(method string) EndpointResponseCodes {
	switch method {
	case http.MethodGet:
		return GetResponseCodes
	case http.MethodPost:
		return CreateResponseCodes
	case http.MethodPut:
		return UpdateResponseCodes
	case http.MethodDelete:
		return DeleteResponseCodes
	default:
		return UnknownVerbResponseCodes
	}
}

type EndpointResponseContent struct {
	Mime    string
	Headers map[string]EndpointResponseHeader
	Paging  types.Optional[interface{}]
	Payload types.Optional[interface{}]
	Example types.Optional[interface{}]
}

func (c EndpointResponseContent) WithHeader(name string, header EndpointResponseHeader) EndpointResponseContent {
	if c.Headers == nil {
		c.Headers = make(map[string]EndpointResponseHeader)
	}
	c.Headers[name] = header
	return c
}

func (c EndpointResponseContent) WithPayload(payload interface{}) EndpointResponseContent {
	if payload == nil {
		c.Payload = types.OptionalEmpty[interface{}]()
	} else {
		c.Payload = types.OptionalOf(payload)
		if c.Mime == "" || c.Mime == "*/*" {
			return c.WithMime(MediaTypeJson)
		}
	}

	return c
}

func (c EndpointResponseContent) WithPaging(paging interface{}) EndpointResponseContent {
	if paging == nil {
		c.Paging = types.OptionalEmpty[interface{}]()
	} else {
		c.Paging = types.OptionalOf(paging)
	}
	return c
}

func (c EndpointResponseContent) WithExample(example interface{}) EndpointResponseContent {
	if example == nil {
		c.Example = types.OptionalEmpty[interface{}]()
	} else {
		c.Example = types.OptionalOf(example)
	}
	return c
}

func (c EndpointResponseContent) WithMime(mime string) EndpointResponseContent {
	c.Mime = mime
	return c
}

func (c EndpointResponseContent) withPortFieldBody(field *ops.PortField) EndpointResponseContent {
	payload := types.Instantiate(field.Type.Type)
	content := c.WithPayload(payload)

	mime := MediaTypeJson
	if mimeOverride := field.Options["mime"]; mimeOverride != "" {
		mime = mimeOverride
	}
	content = content.WithMime(mime)

	return content
}

func (c EndpointResponseContent) withPortFieldPaging(field *ops.PortField) EndpointResponseContent {
	payload := types.Instantiate(field.Type.Type)
	content := c.WithPaging(payload)

	mime := MediaTypeJson
	if mimeOverride := field.Options["mime"]; mimeOverride != "" {
		mime = mimeOverride
	}
	content = content.WithMime(mime)

	return content
}

type EndpointResponseHeader struct {
	Description     *string
	Required        *bool
	Deprecated      *bool
	AllowEmptyValue *bool
	Explode         *bool
	AllowReserved   *bool
	Payload         types.Optional[interface{}]
	Example         types.Optional[interface{}]
	Reference       *string
	PortField       *ops.PortField
}

func (p EndpointResponseHeader) WithDescription(description string) EndpointResponseHeader {
	p.Description = &description
	return p
}

func (p EndpointResponseHeader) WithRequired(required bool) EndpointResponseHeader {
	p.Required = &required
	return p
}

func (p EndpointResponseHeader) WithDeprecated(deprecated bool) EndpointResponseHeader {
	p.Deprecated = &deprecated
	return p
}

func (p EndpointResponseHeader) WithAllowEmptyValue(allowEmptyValue bool) EndpointResponseHeader {
	p.AllowEmptyValue = &allowEmptyValue
	return p
}

func (p EndpointResponseHeader) WithExplode(explode bool) EndpointResponseHeader {
	p.Explode = &explode
	return p
}

func (p EndpointResponseHeader) WithAllowReserved(allowReserved bool) EndpointResponseHeader {
	p.AllowReserved = &allowReserved
	return p
}

func (p EndpointResponseHeader) WithExample(example interface{}) EndpointResponseHeader {
	if example != nil {
		p.Example = types.OptionalOf(example)
	} else {
		p.Example = types.OptionalEmpty[interface{}]()
	}
	return p
}

func (p EndpointResponseHeader) WithPayload(payload interface{}) EndpointResponseHeader {
	if payload != nil {
		p.Payload = types.OptionalOf(payload)
	} else {
		p.Payload = types.OptionalEmpty[interface{}]()
	}
	return p
}

func (p EndpointResponseHeader) WithReference(reference string) EndpointResponseHeader {
	p.Reference = &reference
	return p
}

func (p EndpointResponseHeader) WithPortField(pf *ops.PortField) EndpointResponseHeader {
	p.PortField = pf
	return p
}

func NewEndpointResponseHeader() EndpointResponseHeader {
	return EndpointResponseHeader{}
}

func EndpointResponseHeaderFromPortField(pf *ops.PortField) EndpointResponseHeader {
	header := NewEndpointResponseHeader().
		WithRequired(!pf.Optional).
		WithPortField(pf)

	tagValue := pf.Tags()

	if err := refl.PopulateFieldsFromTags(&header, tagValue); err != nil {
		logger.WithError(err).Error("Failed to populate header from struct tags")
	}

	if example, ok := pf.Options[exampleTag]; ok {
		header = header.WithExample(example)
	}

	if header.Explode == nil {
		header = header.WithExplode(false)
	}

	return header
}
