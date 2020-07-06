package httpclient

import "net/http"

type Factory interface {
	NewHttpClient() *http.Client
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type DoFunc func(req *http.Request) (*http.Response, error)

func (d DoFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return d(request)
}

type RequestInterceptor func(fn DoFunc) DoFunc
