// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cache

import (
	"bytes"
	"net/http"
)

type RecordingHttpResponseWriter struct {
	ResponseWriter http.ResponseWriter
	Body           bytes.Buffer
}

func (i *RecordingHttpResponseWriter) Write(buf []byte) (int, error) {
	j, err := i.Body.Write(buf)
	if err != nil {
		return j, err
	}
	return i.ResponseWriter.Write(buf)
}

func (i *RecordingHttpResponseWriter) WriteHeader(statusCode int) {
	i.ResponseWriter.WriteHeader(statusCode)
}

func (i *RecordingHttpResponseWriter) Header() http.Header {
	return i.ResponseWriter.Header()
}
