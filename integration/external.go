package integration

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/loginterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/traceinterceptor"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Endpoint struct {
	Name   string
	Method string
	Path   string
}

func (e Endpoint) Url(scheme, authority string, variables map[string]string, queryParameters url.Values) (string, error) {
	path := normalizePathTemplate(e.Path)

	subPathTemplate, err := template.New(e.Name).Parse(path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse e url template")
	}

	var pathBuffer bytes.Buffer
	pathBuffer.WriteString(scheme)
	pathBuffer.WriteString("://")
	pathBuffer.WriteString(authority)

	err = subPathTemplate.Execute(&pathBuffer, variables)
	if err != nil {
		return "", errors.Wrap(err, "Failed to fill url template")
	}

	if len(queryParameters) > 0 {
		pathBuffer.WriteString("?")
		pathBuffer.WriteString(queryParameters.Encode())
	}

	result := pathBuffer.String()
	return result, nil
}

// Normalize path variables to go template variable references
func normalizePathTemplate(path string) string {
	if strings.Contains(path, "/{") {
		pathParts := strings.Split(path, "/")
		for i := 0; i < len(pathParts); i++ {
			part := pathParts[i]
			if strings.HasPrefix(part, "{") && !strings.HasPrefix(part, "{{") {
				part = strings.TrimPrefix(part, "{")
				part = strings.TrimSuffix(part, "}")
				part = "{{." + part + "}}"
			}
			pathParts[i] = part
		}
		path = strings.Join(pathParts, "/")
	}
	return path
}

type ExternalService struct {
	ctx          context.Context
	scheme       string
	authority    string
	interceptors []httpclient.RequestInterceptor
}

func (v *ExternalService) AddInterceptor(interceptor httpclient.RequestInterceptor) {
	v.interceptors = append(v.interceptors, interceptor)
}

func (v *ExternalService) Request(endpoint Endpoint, uriVariables map[string]string, queryVariables url.Values, headers http.Header, body []byte) (req *http.Request, err error) {
	fullUrl, err := endpoint.Url(v.scheme, v.authority, uriVariables, queryVariables)
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		buf = bytes.NewBuffer(body)
	} else {
		buf = http.NoBody
	}

	req, err = http.NewRequest(endpoint.Method, fullUrl, buf)
	if err != nil {
		return nil, err
	}

	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	req = req.WithContext(v.ctx)

	return req, nil
}

func (v *ExternalService) Do(req *http.Request, responseBody interface{}) (*http.Response, []byte, error) {
	factory := httpclient.FactoryFromContext(v.ctx)
	if factory == nil {
		return nil, nil, errors.New("Failed to retrieve http client factory from context")
	}

	httpClientDo := factory.NewHttpClient().Do
	for _, interceptor := range v.interceptors {
		httpClientDo = interceptor(httpClientDo)
	}

	resp, err := httpClientDo(req)
	if err != nil {
		return resp, nil, errors.Wrap(err, "Failed to execute request")
	}

	var responseBodyBytes []byte
	if responseBody != nil {
		defer resp.Body.Close()
		if responseBodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
			return resp, nil, errors.Wrap(err, "Failed to read response body")
		}
	}

	if resp.StatusCode > 399 {
		return nil, responseBodyBytes, StatusError{
			Code: resp.StatusCode,
			Body: string(responseBodyBytes),
			Err:  errors.Errorf("Response code %d", resp.StatusCode),
		}
	}

	if len(responseBodyBytes) > 0 {
		if err = json.Unmarshal(responseBodyBytes, responseBody); err != nil {
			return resp, responseBodyBytes, errors.Wrap(err, "Failed to unmarshal body")
		}
	}

	return resp, responseBodyBytes, nil
}

func NewExternalService(ctx context.Context, scheme, authority string) *ExternalService {
	return &ExternalService{
		ctx: ctx,
		scheme:    scheme,
		authority: authority,
		interceptors: []httpclient.RequestInterceptor{
			loginterceptor.NewInterceptor,
			traceinterceptor.NewInterceptor,
		},
	}
}
