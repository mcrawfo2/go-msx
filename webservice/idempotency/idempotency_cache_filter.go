// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idempotency

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
	"net/http"
)

const (
	CacheIdExpectedHeader            = "Idempotency-Key"
	CacheProviderRedis               = "redis"
	CacheProviderInMemory            = "in-memory"
	ConfigRootIdempotencyKey         = "server.idempotency-key"
	ConfigRootIdempotencyKeyRedis    = ConfigRootIdempotencyKey + "." + CacheProviderRedis
	ConfigRootIdempotencyKeyInMemory = ConfigRootIdempotencyKey + "." + CacheProviderInMemory
	IdempotencyDocsLink              = "https://cto-github.cisco.com/NFV-BU/go-msx/blob/main/idempotency.md"
	ValidateIdempotencyRequire       = "REQUIRE"
	ValidateIdempotencyRecommend     = "RECOMMEND"
)

func IdempotencyCacheFilter(idempotencyCache lru.ContextCache) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		cacheId := req.Request.Header.Get(CacheIdExpectedHeader)
		if cacheId == "" { // no cache id then don't bother
			chain.ProcessFilter(req, resp)
			return
		}

		val, exists, err := idempotencyCache.Get(req.Request.Context(), cacheId)
		if err != nil {
			logger.WithContext(req.Request.Context()).Error(err)
			chain.ProcessFilter(req, resp)
			return
		}

		if exists {
			cVal := val.(CachedWebData)

			if isValidCaching(req, cVal) { // quick checks for same method and api
				serveFromCache(resp, cVal)

				logger.WithContext(req.Request.Context()).Debugf("served from cache: %s", cacheId)

				return

			} else {
				// key exists but is not for the same request. Likely reusing
				resp.WriteHeader(http.StatusUnprocessableEntity)
				resp.Write([]byte(`{
					"type": "https://cto-github.cisco.com/NFV-BU/go-msx/blob/main/idempotency.md",
					"title": "Idempotency-Key is already used",
					"detail": "This operation is idempotent and it requires correct usage of Idempotency Key. Idempotency Key MUST not be reused across different payloads of this operation."
				}`))

				return
			}

		} else {
			// wrapping the ResponseWriter so the data could be retrieved later
			rRespWriter := &RecordingHttpResponseWriter{
				ResponseWriter: resp.ResponseWriter,
				Body:           bytes.Buffer{},
			}
			resp.ResponseWriter = rRespWriter
		}

		chain.ProcessFilter(req, resp)

		if !exists { // saving when data is not in cache
			saveToCache(req, resp, idempotencyCache, cacheId)
		}
	}
}

func saveToCache(req *restful.Request, resp *restful.Response, idempotencyCache lru.ContextCache, cacheId string) {
	respStatusCode := resp.StatusCode()
	if respStatusCode >= 200 && respStatusCode < 400 {
		recordingRW := resp.ResponseWriter.(*RecordingHttpResponseWriter)
		body := recordingRW.Body.String()

		creq := CachedRequest{
			Method:     req.Request.Method,
			RequestURI: req.Request.RequestURI,
		}

		cresp := CachedResponse{
			StatusCode: respStatusCode,
			Data:       []byte(body),
			Header:     resp.Header(),
		}

		err := idempotencyCache.Set(req.Request.Context(), cacheId, CachedWebData{
			Req:  creq,
			Resp: cresp,
		})
		if err != nil {
			logger.WithContext(req.Request.Context()).Error(err)
		}
	}
}

func isValidCaching(req *restful.Request, cVal CachedWebData) bool {
	return cVal.Req.Method == req.Request.Method && cVal.Req.RequestURI == req.Request.RequestURI
}

func serveFromCache(resp *restful.Response, cVal CachedWebData) {
	resp.WriteHeader(cVal.Resp.StatusCode)
	resp.Write(cVal.Resp.Data)

	for k, v := range cVal.Resp.Header {
		resp.Header()[k] = v
	}
}

func ApplyIdempotencyKeyFilter(ctx context.Context) (err error) {
	logger.Info("Applying idempotency key filter")

	server := webservice.WebServerFromContext(ctx)
	if server == nil {
		// Server disabled
		return
	}

	cfg := config.MustFromContext(ctx)

	enabled, err := cfg.BoolOr(ConfigRootIdempotencyKey+".enabled", false)
	if err != nil {
		return err
	}
	if !enabled {
		return
	}

	idempotencyCache, err := IdempotencyCacheProviderFactory(ctx)
	if err != nil {
		return err
	}

	server.AddFilter(IdempotencyCacheFilter(idempotencyCache))

	return
}

func IdempotencyCacheProviderFactory(ctx context.Context) (lru.ContextCache, error) {
	cfg := config.MustFromContext(ctx)

	cacheProvider, err := cfg.StringOr(ConfigRootIdempotencyKey+".cache-provider", CacheProviderRedis)
	if err != nil {
		return nil, err
	}

	configRoot := config.PrefixWithName(ConfigRootIdempotencyKey, cacheProvider)

	return lru.NewContextCache(ctx, cacheProvider, configRoot)
}

func ValidateIdempotency(mode string) restfulcontext.RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
			cacheId := req.Request.Header.Get(CacheIdExpectedHeader)
			if cacheId == "" {
				if mode == ValidateIdempotencyRequire {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{
						"type": "` + IdempotencyDocsLink + `",
						"title": "Idempotency-Key is missing",
						"detail": "This operation is idempotent and it requires correct usage of Idempotency Key."
					}`))

					return

				} else if mode == ValidateIdempotencyRecommend {
					resp.Header()["Link"] = []string{IdempotencyDocsLink}
					resp.Header()["go-msx-idempotency"] = []string{"This operation is idempotent and it requires correct usage of Idempotency Key."}
				}
			}

			chain.ProcessFilter(req, resp)
		})
	}
}
