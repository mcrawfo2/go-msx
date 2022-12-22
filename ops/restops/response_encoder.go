// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/spf13/cast"
	"io"
	"sort"
	"strings"
)

type ResponseEncoder interface {
	EncodeHeaderPrimitive(name string, values types.Optional[string], style string, explode bool) (err error)
	EncodeHeaderArray(name string, values []string, style string, explode bool) (err error)
	EncodeHeaderObject(name string, value types.Pojo, style string, explode bool) (err error)
	EncodeMime(mime string) (err error)
	EncodeCode(code int) (err error)
	EncodeBody(body interface{}) error
}

type EndpointResponseEncoder struct {
	Sink ResponseDataSink
}

func (o EndpointResponseEncoder) EncodeHeaderPrimitive(name string, value types.Optional[string], _ string, _ bool) (err error) {
	o.Sink.UnsetHeader(name)
	if value.IsPresent() {
		o.Sink.SetHeader(name, value.Value())
	}
	return nil
}

func (o EndpointResponseEncoder) EncodeHeaderArray(name string, values []string, _ string, _ bool) (err error) {
	o.Sink.UnsetHeader(name)
	if len(values) == 0 {
		return nil
	}

	o.Sink.AddHeader(name, strings.Join(values, ","))
	return nil
}

type keyResult struct {
	Key    string
	Result string
}

func (o EndpointResponseEncoder) EncodeHeaderObject(name string, value types.Pojo, _ string, explode bool) (err error) {
	o.Sink.UnsetHeader(name)
	if len(value) == 0 {
		return nil
	}

	var keys []string

	// Always uses "simple" style
	fieldSep := ","
	kvSep := ","
	if explode {
		kvSep = "="
	}

	for k := range value {
		keys = append(keys, k)
	}

	// sort results by key for stable output
	sort.Strings(keys)

	var result strings.Builder
	for i, k := range keys {
		if i > 0 {
			result.WriteString(fieldSep)
		}
		result.WriteString(k)
		result.WriteString(kvSep)
		result.WriteString(cast.ToString(value[k]))
	}

	o.Sink.AddHeader(name, result.String())

	return nil
}

func (o EndpointResponseEncoder) EncodeCode(code int) (err error) {
	o.Sink.SetStatus(code)
	return nil
}

func (o EndpointResponseEncoder) EncodeMime(mime string) (err error) {
	switch mime {
	case MediaTypeJson:
		mime = ContentTypeJson
	case MediaTypeXml:
		mime = ContentTypeXml
	}
	o.Sink.AddHeader(HeaderContentType, mime)
	return nil
}

func (o EndpointResponseEncoder) EncodeBody(body interface{}) (err error) {
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
	case nil:
		err = o.encodeNoBody()
	default:
		err = o.encodeBodyEntity(body)
	}
	return
}

func (o EndpointResponseEncoder) encodeBodyBytes(body []byte) (err error) {
	reader := io.NopCloser(bytes.NewReader(body))
	return o.Sink.WriteBody(reader)
}

func (o EndpointResponseEncoder) encodeBodyReader(reader io.ReadCloser) (err error) {
	return o.Sink.WriteBody(reader)
}

func (o EndpointResponseEncoder) encodeBodyEntity(body interface{}) (err error) {
	return o.Sink.WriteBodyEntity(body)
}

func (o EndpointResponseEncoder) encodeBodyString(body string) (err error) {
	reader := io.NopCloser(strings.NewReader(body))
	return o.Sink.WriteBody(reader)
}

func (o EndpointResponseEncoder) encodeBodyTextMarshaler(body types.TextMarshaler) (err error) {
	bodyString, err := body.MarshalText()
	if err != nil {
		return err
	}

	return o.encodeBodyString(bodyString)
}

func (o EndpointResponseEncoder) encodeNoBody() error {
	return o.Sink.WriteBody(nil)
}

func NewResponseEncoder(sink ResponseDataSink) ResponseEncoder {
	return &EndpointResponseEncoder{
		Sink: sink,
	}
}
