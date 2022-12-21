// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"io"
	"net/http"
)

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
	Response *restful.Response
}

func (r *RestfulResponseDataSink) SetHeader(name string, value string) {
	r.Response.Header().Set(name, value)
}

func (r *RestfulResponseDataSink) AddHeader(name string, value string) {
	r.Response.Header().Add(name, value)
}

func (r *RestfulResponseDataSink) UnsetHeader(name string) {
	r.Response.Header().Del(name)
}

func (r *RestfulResponseDataSink) SetStatus(code int) {
	r.Status = code
}

func (r *RestfulResponseDataSink) WriteBody(src io.ReadCloser) (err error) {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}

	if src != nil {
		defer src.Close()
	}

	r.Response.WriteHeader(r.Status)
	if src != nil {
		_, err = io.Copy(r.Response, src)
	}
	return err
}

func (r *RestfulResponseDataSink) WriteBodyEntity(entity interface{}) error {
	contentOptions := ops.ContentOptions{
		MimeType: r.Response.Header().Get(HeaderContentType),
	}

	buffer := types.CloseableByteBuffer{
		Buffer: new(bytes.Buffer),
	}

	err := contentOptions.WriteEntity(buffer, entity)
	if err != nil {
		return err
	}

	return r.WriteBody(buffer)
}

func NewRestfulResponseDataSink(resp *restful.Response) *RestfulResponseDataSink {
	return &RestfulResponseDataSink{
		Response: resp,
	}
}
