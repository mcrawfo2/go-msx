// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockResponseDataSink struct {
	code    types.Optional[int]
	headers http.Header
	body    []byte
}

func (m *MockResponseDataSink) UnsetHeader(name string) {
	m.headers.Del(name)
}

func (m *MockResponseDataSink) SetHeader(name string, value string) {
	m.headers.Set(name, value)
}

func (m *MockResponseDataSink) AddHeader(name string, value string) {
	m.headers.Add(name, value)
}

func (m *MockResponseDataSink) SetStatus(code int) {
	m.code = types.OptionalOf(code)
}

func (m *MockResponseDataSink) WriteBody(src io.ReadCloser) error {
	if src == nil {
		return nil
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	m.body = data
	return nil
}

func (m *MockResponseDataSink) WriteBodyEntity(entity interface{}) error {
	contentType := m.headers.Get(HeaderContentType)
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.Wrap(err, "Failed to parse media type")
	}

	contentOptions := ops.NewContentOptions(mediaType)
	bodyBuffer := types.CloseableByteBuffer{Buffer: bytes.NewBuffer(m.body)}
	return contentOptions.WriteEntity(bodyBuffer, entity)
}

func TestRestfulResponseDataSink_SetHeader(t *testing.T) {
	resp := &httptest.ResponseRecorder{}
	r := &RestfulResponseDataSink{
		Response: restful.NewResponse(resp),
	}

	r.SetHeader("My-Header", "abc")

	assert.Equal(t,
		resp.Header(),
		http.Header{
			"My-Header": {"abc"},
		})
}

func TestRestfulResponseDataSink_UnsetHeader(t *testing.T) {
	resp := &httptest.ResponseRecorder{
		HeaderMap: http.Header{
			"My-Header": {"abc"},
		},
	}
	r := &RestfulResponseDataSink{
		Response: restful.NewResponse(resp),
	}

	r.UnsetHeader("My-Header")

	assert.Equal(t,
		resp.Header(),
		http.Header{})
}

func TestRestfulResponseDataSink_AddHeader(t *testing.T) {
	resp := &httptest.ResponseRecorder{
		HeaderMap: http.Header{
			"My-Header": {"abc"},
		},
	}
	r := &RestfulResponseDataSink{
		Response: restful.NewResponse(resp),
	}

	r.AddHeader("My-Header", "123")

	assert.Equal(t,
		resp.Header(),
		http.Header{
			"My-Header": {"abc", "123"},
		})
}

func TestRestfulResponseDataSink_SetStatus(t *testing.T) {
	resp := &httptest.ResponseRecorder{}
	r := &RestfulResponseDataSink{
		Response: restful.NewResponse(resp),
	}

	r.SetStatus(http.StatusResetContent)

	assert.Equal(t, http.StatusResetContent, r.Status)
}

func TestRestfulResponseDataSink_WriteBody(t *testing.T) {
	resp := &httptest.ResponseRecorder{
		Body: new(bytes.Buffer),
	}
	r := &RestfulResponseDataSink{
		Status:   http.StatusOK,
		Response: restful.NewResponse(resp),
	}

	responseBodyString := `{"a":123}`
	err := r.WriteBody(io.NopCloser(bytes.NewBufferString(responseBodyString)))
	assert.NoError(t, err)

	result := resp.Result()
	assert.NotNil(t, result)

	resultBody, err := ioutil.ReadAll(result.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, responseBodyString, string(resultBody))
}

func TestRestfulResponseDataSink_WriteBodyEntity(t *testing.T) {
	resp := &httptest.ResponseRecorder{
		Body: new(bytes.Buffer),
	}

	r := &RestfulResponseDataSink{
		Response: restful.NewResponse(resp),
	}
	r.Response.AddHeader(HeaderContentType, ContentTypeJson)
	r.Response.SetRequestAccepts(MediaTypeJson)

	responseBodyEntity := map[string]interface{}{
		"a": 123,
	}
	responseBodyBytes, err := json.Marshal(responseBodyEntity)
	assert.NoError(t, err)
	responseBodyBytes = append(responseBodyBytes, '\n')

	err = r.WriteBodyEntity(responseBodyEntity)
	assert.NoError(t, err)

	result := resp.Result()
	assert.NotNil(t, result)

	resultBody, err := ioutil.ReadAll(result.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, responseBodyBytes, resultBody)
}
