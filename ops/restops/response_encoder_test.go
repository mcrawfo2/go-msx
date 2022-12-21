// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEndpointResponseEncoder_EncodeBody_Entity(t *testing.T) {
	resp := new(http.Response)

	encoder := EndpointResponseEncoder{
		Sink: NewHttpResponseDataSink(resp),
	}

	err := encoder.EncodeBody(types.Pojo{"abc": 123})
	assert.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"abc":123}`, string(bodyBytes))
}

func TestEndpointResponseEncoder_EncodeBody_ReadCloser(t *testing.T) {
	resp := new(http.Response)

	encoder := EndpointResponseEncoder{
		Sink: NewHttpResponseDataSink(resp),
	}

	wantBodyBytes, err := json.Marshal(types.Pojo{"abc": 123})
	assert.NoError(t, err)

	err = encoder.EncodeBody(io.NopCloser(bytes.NewBuffer(wantBodyBytes)))
	assert.NoError(t, err)

	gotBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, wantBodyBytes, gotBodyBytes)
}

func TestEndpointResponseEncoder_EncodeBody_ByteSlice(t *testing.T) {
	resp := new(http.Response)

	encoder := EndpointResponseEncoder{
		Sink: NewHttpResponseDataSink(resp),
	}

	wantBodyBytes, err := json.Marshal(types.Pojo{"abc": 123})
	assert.NoError(t, err)

	err = encoder.EncodeBody(wantBodyBytes)
	assert.NoError(t, err)

	gotBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, wantBodyBytes, gotBodyBytes)
}

func TestEndpointResponseEncoder_EncodeBody_String(t *testing.T) {
	resp := new(http.Response)

	encoder := EndpointResponseEncoder{
		Sink: NewHttpResponseDataSink(resp),
	}

	wantBodyBytes, err := json.Marshal(types.Pojo{"abc": 123})
	assert.NoError(t, err)

	err = encoder.EncodeBody(string(wantBodyBytes))
	assert.NoError(t, err)

	gotBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, wantBodyBytes, gotBodyBytes)
}

func TestEndpointResponseEncoder_EncodeBody_TextMarshaler(t *testing.T) {
	resp := new(http.Response)

	encoder := EndpointResponseEncoder{
		Sink: NewHttpResponseDataSink(resp),
	}

	bodyContent := types.MustNewUUID()
	wantBodyBytes := []byte(bodyContent.String())

	err := encoder.EncodeBody(bodyContent)
	assert.NoError(t, err)

	gotBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, wantBodyBytes, gotBodyBytes)
}
