package mocks

import (
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type RoundTripFunc func(request *http.Request) (response *http.Response, e error)

type HttpRoute struct {
	Method  string
	URL     string
	Tripper RoundTripFunc
}

type HttpRouteBuilder struct {
	Url        string
	Method     string
	StatusCode int
	Status     string
	BodyPath   string
	Headers    map[string]string
}

func (d HttpRouteBuilder) Build() HttpRoute {
	responseHeaders := func() http.Header {
		result := make(http.Header)
		for k, v := range d.Headers {
			result.Set(k, v)
		}
		return result
	}

	return HttpRoute{
		URL:    d.Url,
		Method: strings.ToUpper(d.Method),
		Tripper: func(request *http.Request) (response *http.Response, e error) {
			var body io.ReadCloser
			if d.BodyPath != "" {
				var err error
				body, err = os.Open(d.BodyPath)
				if err != nil {
					panic(err.Error())
				}
			}

			return &http.Response{
				Status:     d.Status,
				StatusCode: d.StatusCode,
				Body:       body,
				Header:     responseHeaders(),
			}, nil
		},
	}
}

type HttpClientFactory struct {
	routes []HttpRoute
}

func (f *HttpClientFactory) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	method := req.Method
	for _, mock := range f.routes {
		if url == mock.URL && method == mock.Method {
			return mock.Tripper(req)
		}
	}
	return nil, errors.Errorf("No mock for %q %q", method, url)
}

func (f *HttpClientFactory) AddRouteDefinition(definition HttpRouteBuilder) {
	f.AddRoute(definition.Build())
}

func (f *HttpClientFactory) AddRoute(route HttpRoute) {
	f.routes = append(f.routes, route)
}

func (f *HttpClientFactory) NewHttpClient() *http.Client {
	return &http.Client{Transport: f}
}

func NewMockHttpClientFactory() *HttpClientFactory {
	return &HttpClientFactory{}
}
