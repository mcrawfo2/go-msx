// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	encodingGzip    = "gzip"
	encodingDeflate = "deflate"
)

var ErrUnknownFieldSource = errors.New("Unknown field source")
var ErrUnsupportedStyle = errors.New("Unsupported style")
var ErrMissingRequiredFile = errors.New("Missing required file")

type RequestDataSource interface {
	Cookies() []*http.Cookie
	Headers() http.Header
	Form() (url.Values, *multipart.Form, error)
	Query() url.Values
	PathParameters() map[string]string
	Body() ([]byte, error)
	BodyContentOptions(string, string) ops.ContentOptions
}

type RestfulRequestDataSource struct {
	Request  *restful.Request
	BodyData []byte
}

func (r *RestfulRequestDataSource) PathParameters() map[string]string {
	return r.Request.PathParameters()
}

func (r *RestfulRequestDataSource) BodyContentOptions(defaultContentType string, defaultContentEncoding string) ops.ContentOptions {
	request := r.Request.Request

	contentType := request.Header.Get(HeaderContentType)
	if contentType == "" {
		contentType = defaultContentType
	}

	contentOptions := ops.NewContentOptions(contentType)

	contentEncoding := request.Header.Get(HeaderContentEncoding)
	if contentEncoding == "" {
		contentEncoding = defaultContentEncoding
	}
	if contentEncoding != "" {
		contentOptions.WithEncoding(strings.Split(contentEncoding, ",")...)
	}
	return contentOptions
}

func (r *RestfulRequestDataSource) DecodedBody() (io.ReadCloser, error) {
	contentOptions := r.BodyContentOptions(ContentTypeJson, ContentEncodingNone)

	bodyBytes, err := r.Body()
	if err != nil {
		return nil, err
	}

	if bodyBytes == nil {
		return nil, nil
	}

	content := ops.NewContentFromBytes(contentOptions, bodyBytes)
	return content.Reader()
}

func (r *RestfulRequestDataSource) Form() (url.Values, *multipart.Form, error) {
	if r.Request.Request.MultipartForm == nil {
		bodyData, err := r.DecodedBody()
		if err != nil {
			return nil, nil, err
		}

		// Attach the decoded body reader
		r.Request.Request.Body = bodyData

		err = r.Request.Request.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			return nil, nil, err
		}
	}

	return r.Request.Request.Form, r.Request.Request.MultipartForm, nil
}

func (r *RestfulRequestDataSource) Query() url.Values {
	return r.Request.Request.URL.Query()
}

func (r *RestfulRequestDataSource) Headers() http.Header {
	return r.Request.Request.Header
}

func (r *RestfulRequestDataSource) Cookies() []*http.Cookie {
	return r.Request.Request.Cookies()
}

func (r *RestfulRequestDataSource) Body() (result []byte, err error) {
	if r.BodyData != nil {
		return r.BodyData, nil
	}

	if r.Request.Request.Body == nil {
		return nil, nil
	}

	// Store the body
	bodyReader := r.Request.Request.Body
	defer bodyReader.Close()

	r.BodyData, err = ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	return r.BodyData, nil
}

func NewRestfulRequestDataSource(req *restful.Request) RequestDataSource {
	return &RestfulRequestDataSource{
		Request: req,
	}
}
