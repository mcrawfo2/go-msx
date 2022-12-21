// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponseDataSink struct {
	Response *http.Response
}

func (r *HttpResponseDataSink) UnsetHeader(name string) {
	r.Response.Header.Del(name)
}

func (r *HttpResponseDataSink) SetHeader(name string, value string) {
	r.Response.Header.Set(name, value)
}

func (r *HttpResponseDataSink) AddHeader(name string, value string) {
	r.Response.Header.Add(name, value)
}

func (r *HttpResponseDataSink) SetStatus(code int) {
	r.Response.StatusCode = code
	r.Response.Status = http.StatusText(code)
}

func (r *HttpResponseDataSink) WriteBody(src io.ReadCloser) error {
	//r.Response.WriteHeader(r.Status)
	var buffer = new(bytes.Buffer)
	_, err := io.Copy(buffer, src)
	if err == nil {
		r.Response.Body = io.NopCloser(buffer)
	}
	return err
}

func (r *HttpResponseDataSink) WriteBodyEntity(entity interface{}) error {
	bodyBytes, err := json.Marshal(entity)
	if err == nil {
		r.Response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return err
}

func NewHttpResponseDataSink(resp *http.Response) *HttpResponseDataSink {
	return &HttpResponseDataSink{
		Response: resp,
	}
}
