package webservice

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/spf13/cast"
	"io"
	"strings"
)

type RequestDescriber interface {
	Parameters() map[string]interface{}
	Path() string
}

type RestfulRequestDescriber struct {
	Request *restful.Request
}

func (r RestfulRequestDescriber) Path() string {
	return r.Request.Request.URL.Path
}

func (r RestfulRequestDescriber) Parameters() map[string]interface{} {
	// TODO: Return parameters for envelope
	return make(map[string]interface{})
}

type EndpointRequestDescriber struct {
	Endpoint Endpoint
}

func (e EndpointRequestDescriber) Parameters() map[string]interface{} {
	// TODO: Return parameters for example envelope
	return make(map[string]interface{})
}

func (e EndpointRequestDescriber) Path() string {
	return e.Endpoint.Path
}

type ResponseDataSink interface {
	UnsetHeader(name string)
	SetHeader(name string, value string)
	AddHeader(name string, value string)
	SetStatus(code int)
	WriteBody(src io.ReadCloser) error
	WriteBodyEntity(entity interface{}) error
}

type RestfulResponseDataSink struct {
	Status   int
	Mime     string
	Response *restful.Response
}

func (r RestfulResponseDataSink) UnsetHeader(name string) {
	r.Response.Header().Del(name)
}

func (r RestfulResponseDataSink) SetHeader(name string, value string) {
	r.Response.Header().Set(name, value)
}

func (r RestfulResponseDataSink) AddHeader(name string, value string) {
	r.Response.Header().Add(name, value)
}

func (r *RestfulResponseDataSink) SetStatus(code int) {
	r.Status = code
}

func (r RestfulResponseDataSink) WriteBody(src io.ReadCloser) error {
	r.Response.WriteHeader(r.Status)
	_, err := io.Copy(r.Response, src)
	return err
}

func (r RestfulResponseDataSink) WriteBodyEntity(entity interface{}) error {
	return r.Response.WriteHeaderAndEntity(r.Status, entity)
}

type ResponseEncoder interface {
	EncodeHeaderPrimitive(name string, values types.OptionalString, style string, explode bool) (err error)
	EncodeHeaderArray(name string, values []string, style string, explode bool) (err error)
	EncodeHeaderObject(name string, value types.Pojo, style string, explode bool) (err error)
	EncodeMime(mime string) (err error)
	EncodeCode(code int) (err error)
	EncodeBody(body interface{}) error
}

type OpenApiResponseEncoder struct {
	Sink ResponseDataSink
}

func (o OpenApiResponseEncoder) EncodeHeaderPrimitive(name string, value types.OptionalString, _ string, _ bool) (err error) {
	o.Sink.UnsetHeader(name)
	if value.IsPresent() {
		o.Sink.SetHeader(name, cast.ToString(value))
	}
	return nil
}

func (o OpenApiResponseEncoder) EncodeHeaderArray(name string, values []string, _ string, _ bool) (err error) {
	o.Sink.UnsetHeader(name)
	if len(values) == 0 {
		return nil
	}

	o.Sink.AddHeader(name, strings.Join(values, ","))
	return nil
}

func (o OpenApiResponseEncoder) EncodeHeaderObject(name string, value types.Pojo, _ string, explode bool) (err error) {
	o.Sink.UnsetHeader(name)
	if len(value) == 0 {
		return nil
	}

	var result strings.Builder

	fieldSep := ","
	kvSep := ","
	if explode {
		kvSep = "="
	}

	for k, v := range value {
		if result.Len() > 0 {
			result.WriteString(fieldSep)
		}
		stringValue := cast.ToString(v)
		result.WriteString(k)
		result.WriteString(kvSep)
		result.WriteString(stringValue)
	}

	o.Sink.AddHeader(name, result.String())

	return nil
}

func (o OpenApiResponseEncoder) EncodeCode(code int) (err error) {
	o.Sink.SetStatus(code)
	return nil
}

func (o OpenApiResponseEncoder) EncodeMime(mime string) (err error) {
	if mime == MIME_JSON {
		mime = MIME_JSON_CHARSET
	}
	o.Sink.AddHeader(headerNameContentType, mime)
	return nil
}

func (o OpenApiResponseEncoder) EncodeBody(body interface{}) (err error) {
	// Encode body
	switch typedBody := body.(type) {
	case io.ReadCloser:
		err = o.encodeBodyReader(typedBody)
	case []byte:
		err = o.encodeBodyBytes(typedBody)
	case string:
		err = o.encodeBodyString(typedBody)
	case types.TextMarshaler:
		err = o.encodeBodyTextMarshaler(typedBody)
	default:
		err = o.encodeBodyEntity(body)
	}
	return
}

func (o OpenApiResponseEncoder) encodeBodyBytes(body []byte) (err error) {
	reader := io.NopCloser(bytes.NewReader(body))
	return o.Sink.WriteBody(reader)
}

func (o OpenApiResponseEncoder) encodeBodyReader(reader io.ReadCloser) (err error) {
	return o.Sink.WriteBody(reader)
}

func (o OpenApiResponseEncoder) encodeBodyEntity(body interface{}) (err error) {
	return o.Sink.WriteBodyEntity(body)
}

func (o OpenApiResponseEncoder) encodeBodyString(body string) (err error) {
	reader := io.NopCloser(strings.NewReader(body))
	return o.Sink.WriteBody(reader)
}

func (o OpenApiResponseEncoder) encodeBodyTextMarshaler(body types.TextMarshaler) (err error) {
	bodyString, err := body.MarshalText()
	if err != nil {
		return err
	}

	return o.encodeBodyString(bodyString)
}
