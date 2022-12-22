// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"github.com/emicklei/go-restful"
	"net/http"
)

const CacheIdExpectedHeader = "Idempotency-Key"

type CachedData struct {
	req  CachedRequest
	resp CachedResponse
}

type CachedRequest struct {
	Method     string
	RequestURI string
}

type CachedResponse struct {
	StatusCode int
	Data       []byte
	Header     http.Header
}

type RecordingHttpResponseWriter struct {
	w    http.ResponseWriter
	body bytes.Buffer
}

func (i *RecordingHttpResponseWriter) Write(buf []byte) (int, error) {
	j, err := i.body.Write(buf)
	if err != nil {
		return j, err
	}
	return i.w.Write(buf)
}

func (i *RecordingHttpResponseWriter) WriteHeader(statusCode int) {
	i.w.WriteHeader(statusCode)
}

func (i *RecordingHttpResponseWriter) Header() http.Header {
	return i.w.Header()
}

func RecordingCacheFilter(recordingCache lru.Cache) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		cacheId := req.Request.Header.Get(CacheIdExpectedHeader)
		if cacheId == "" { // no cache id then don't bother
			chain.ProcessFilter(req, resp)
			return
		}

		val, exists := recordingCache.Get(cacheId)
		if exists {
			cVal := val.(CachedData)

			if isValidCaching(req, cVal) { // quick checks
				serveFromCache(resp, cVal)

				logger.WithContext(req.Request.Context()).Debugf("served from cache: %s", cacheId)

				// the request has been processed
				return

			} // else when invalid just process the api as usual (no cache and no saving)

		} else {
			// wrapping the ResponseWriter so the data could be retrieved later
			rRespWriter := &RecordingHttpResponseWriter{
				w:    resp.ResponseWriter,
				body: bytes.Buffer{},
			}
			resp.ResponseWriter = rRespWriter
		}

		chain.ProcessFilter(req, resp)

		if !exists { // saving when data is not in cache
			saveToCache(req, resp, recordingCache, cacheId)
		}
	}
}

func saveToCache(req *restful.Request, resp *restful.Response, recordingCache lru.Cache, cacheId string) {
	respStatusCode := resp.StatusCode()
	if respStatusCode >= 200 && respStatusCode < 400 {
		recordingRW := resp.ResponseWriter.(*RecordingHttpResponseWriter)
		body := recordingRW.body.String()

		creq := CachedRequest{
			Method:     req.Request.Method,
			RequestURI: req.Request.RequestURI,
		}

		cresp := CachedResponse{
			StatusCode: respStatusCode,
			Data:       []byte(body),
			Header:     resp.Header(),
		}

		recordingCache.Set(cacheId, CachedData{
			req:  creq,
			resp: cresp,
		})
	}
}

func isValidCaching(req *restful.Request, cVal CachedData) bool {
	return cVal.req.Method == req.Request.Method && cVal.req.RequestURI == req.Request.RequestURI
}

func serveFromCache(resp *restful.Response, cVal CachedData) {
	resp.WriteHeader(cVal.resp.StatusCode)
	resp.Write(cVal.resp.Data)

	for k, v := range cVal.resp.Header {
		resp.Header()[k] = v
	}
}
