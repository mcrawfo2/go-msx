// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type MockRequestDataSource struct {
	cookies        []*http.Cookie
	headers        http.Header
	form           url.Values
	multipartForm  *multipart.Form
	query          url.Values
	pathParameters map[string]string
	body           []byte
}

func (m MockRequestDataSource) Body() ([]byte, error) {
	return m.body, nil
}

func (m MockRequestDataSource) Cookies() []*http.Cookie {
	return m.cookies
}

func (m MockRequestDataSource) Headers() http.Header {
	return m.headers
}

func (m MockRequestDataSource) Form() (url.Values, *multipart.Form, error) {
	return m.form, m.multipartForm, nil
}

func (m MockRequestDataSource) Query() url.Values {
	return m.query
}

func (m MockRequestDataSource) PathParameters() map[string]string {
	return m.pathParameters
}

func (m MockRequestDataSource) ReadEntity(e interface{}) (err error) {
	return json.Unmarshal(m.body, e)
}

func (m MockRequestDataSource) BodyContentOptions(defaultContentType string, defaultContentEncoding string) ops.ContentOptions {
	contentType := types.NewOptionalStringFromString(m.headers.Get(HeaderContentType)).NilIfEmpty().OrElse(defaultContentType)
	contentEncoding := types.NewOptionalStringFromString(m.headers.Get(HeaderContentEncoding)).NilIfEmpty().OrElse(defaultContentEncoding)

	options := ops.NewContentOptions(contentType)
	if contentEncoding != "" {
		options.WithEncoding(strings.Split(contentEncoding, "")...)
	}
	return options
}

func TestRestfulRequestDataSource_PathParameters(t *testing.T) {
	s := RestfulRequestDataSource{Request: new(restful.Request)}
	got := s.PathParameters()
	assert.Nil(t, got)
}

func TestRestfulRequestDataSource_BodyContentOptions(t *testing.T) {
	tests := []struct {
		name                   string
		request                *http.Request
		defaultContentType     string
		defaultContentEncoding string
		want                   ops.ContentOptions
	}{
		{
			name: "FullySpecified",
			request: &http.Request{
				Header: map[string][]string{
					"Content-Type":     {ContentTypeJson},
					"Content-Encoding": {ContentEncodingGzip},
				},
			},
			want: ops.ContentOptions{
				MimeType: ContentTypeJson,
				Encoding: []string{ContentEncodingGzip},
			},
		},
		{
			name: "DefaultEncoding",
			request: &http.Request{
				Header: map[string][]string{
					"Content-Type": {ContentTypeJson},
				},
			},
			defaultContentEncoding: ContentEncodingGzip,
			want: ops.ContentOptions{
				MimeType: ContentTypeJson,
				Encoding: []string{ContentEncodingGzip},
			},
		},
		{
			name: "NoEncoding",
			request: &http.Request{
				Header: map[string][]string{
					"Content-Type": {ContentTypeJson},
				},
			},
			want: ops.ContentOptions{
				MimeType: ContentTypeJson,
				Encoding: nil,
			},
		},
		{
			name: "NoContentType",
			request: &http.Request{
				Header: map[string][]string{},
			},
			defaultContentType: ContentTypeJson,
			want: ops.ContentOptions{
				MimeType: ContentTypeJson,
				Encoding: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestfulRequestDataSource{
				Request: restful.NewRequest(tt.request),
			}
			got := r.BodyContentOptions(tt.defaultContentType, tt.defaultContentEncoding)
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestRestfulRequestDataSource_DecodedBody(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		bodyData string
		wantData string
		wantErr  bool
	}{
		{
			name: "Unencoded",
			request: &http.Request{
				Header: map[string][]string{
					HeaderContentType: {MediaTypeTextPlain},
				},
			},
			bodyData: `abc123`,
			wantData: `abc123`,
		},
		{
			name: "Base64",
			request: &http.Request{
				Header: map[string][]string{
					HeaderContentType:     {MediaTypeTextPlain},
					HeaderContentEncoding: {ContentEncodingBase64},
				},
			},
			bodyData: `YWJjMTIz`,
			wantData: `abc123`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestfulRequestDataSource{
				Request:  restful.NewRequest(tt.request),
				BodyData: []byte(tt.bodyData),
			}
			got, err := r.DecodedBody()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				gotData, err := io.ReadAll(got)
				assert.NoError(t, err)
				assert.Equal(t,
					tt.wantData,
					string(gotData),
					testhelpers.Diff(tt.wantData, string(gotData)))
			}
		})
	}
}

func TestRestfulRequestDataSource_Form(t *testing.T) {
	r := &RestfulRequestDataSource{
		Request: restful.NewRequest(&http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				RawQuery: `a=1&b=2`,
			},
			Header: http.Header{
				HeaderContentType: {`multipart/form-data; boundary="boundary"`},
			},
			Body: io.NopCloser(bytes.NewBufferString(
				`--boundary
Content-Disposition: form-data; name="field1"

value1
--boundary
Content-Disposition: form-data; name="field2"; filename="example.txt"

value2
--boundary--`,
				)),
		}),
	}

	wantValues := url.Values{
		"a": {"1"},
		"b": {"2"},
		"field1": {"value1"},
	}

	wantForm := &multipart.Form{
		Value: url.Values{
			"field1": {"value1"},
		},
		File: map[string][]*multipart.FileHeader{
			"field2": {
				{
					Filename: "example.txt",
					Header: textproto.MIMEHeader{
						HeaderContentDisposition: {`form-data; name="field2"; filename="example.txt"`},
					},
					Size: 6,
				},
			},
		},
	}



	gotValues, gotForm, err := r.Form()

	assert.NoError(t, err)

	assert.True(t,
		reflect.DeepEqual(wantValues, gotValues),
		testhelpers.Diff(wantValues, gotValues))

	assert.Equal(t,
		wantForm.File["field2"][0].Filename,
		gotForm.File["field2"][0].Filename)
	wantForm.File["field2"][0] = gotForm.File["field2"][0]

	assert.True(t,
		reflect.DeepEqual(wantForm, gotForm),
		testhelpers.Diff(wantForm, gotForm))
}

func TestRestfulRequestDataSource_Query(t *testing.T) {
	r := &RestfulRequestDataSource{
		Request: restful.NewRequest(&http.Request{
			URL: &url.URL{
				RawQuery: "a=1&b=2",
			},
		}),
	}

	want := url.Values{
		"a": {"1"},
		"b": {"2"},
	}

	values := r.Query()
	assert.True(t,
		reflect.DeepEqual(want, values),
		testhelpers.Diff(want, values))
}

func TestRestfulRequestDataSource_Headers(t *testing.T) {
	r := &RestfulRequestDataSource{
		Request: restful.NewRequest(&http.Request{
			Header: http.Header{
				"a": {"1"},
				"b": {"2"},
			},
		}),
	}

	want := http.Header{
		"a": {"1"},
		"b": {"2"},
	}

	values := r.Headers()
	assert.True(t,
		reflect.DeepEqual(want, values),
		testhelpers.Diff(want, values))
}

func TestRestfulRequestDataSource_Cookies(t *testing.T) {
	r := &RestfulRequestDataSource{
		Request: restful.NewRequest(&http.Request{
			Header: http.Header{
				"Cookie": {
					"SESSION=abc",
					"JSESSION=123",
				},
			},
		}),
	}

	want := []*http.Cookie{
		{
			Name:  "SESSION",
			Value: "abc",
		},
		{
			Name:  "JSESSION",
			Value: "123",
		},
	}

	values := r.Cookies()

	assert.True(t,
		reflect.DeepEqual(want, values),
		testhelpers.Diff(want, values))
}

func TestRestfulRequestDataSource_Body(t *testing.T) {
	r := &RestfulRequestDataSource{
		Request: restful.NewRequest(&http.Request{
			Body: io.NopCloser(bytes.NewBufferString("abc123")),
		}),
	}

	want := []byte("abc123")

	value, err := r.Body()
	assert.NoError(t, err)
	assert.Equal(t, want, value)
}