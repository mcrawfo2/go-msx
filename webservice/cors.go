package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
)

const (
	MIME_XML              = restful.MIME_XML  // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_JSON             = restful.MIME_JSON // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_TEXT_PLAIN       = "text/plain"      // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_APPLICATION_FORM = "application/x-www-form-urlencoded"
	MIME_MULTIPART_FORM   = "multipart/form-data"

	headerNameAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	headerNameAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	headerNameAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	headerNameAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	headerNameAuthorization                 = "Authorization"
	headerNameAccessToken                   = "access_token"
	headerNameCacheControl                  = "Cache-Control"
	headerNameCurrency                      = "currency"
	headerNameIfModifiedSince               = "If-Modified-Since"
	headerNameAccept                        = restful.HEADER_Accept
	headerNameLocale                        = "locale"
	headerNamePragma                        = "pragma"
	headerNameContentType                   = restful.HEADER_ContentType
	headerNameContentLength                 = "Content-Length"
	headerNameContentEncoding               = restful.HEADER_ContentEncoding
	headerNameAcceptEncoding                = restful.HEADER_AcceptEncoding
	headerNameXCsrfToken                    = "X-CSRF-Token"
	headerNameApiKey                        = "api_key"
	headerNameXRequestedWith                = "x-requested-with"
	headerNameAccessControlRequestMethod    = restful.HEADER_AccessControlRequestMethod

	// Deprecated
	HEADER_ContentEncoding = "Content-Encoding"
	// Deprecated
	HEADER_AccessControlRequestHeaders = "Access-Control-Request-Headers"
	// Deprecated
	HEADER_AccessControlAllowMethods = "Access-Control-Allow-Methods"
	// Deprecated
	HEADER_AccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// Deprecated
	HEADER_AccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// Deprecated
	HEADER_Authorization = "Authorization"
	// Deprecated
	HEADER_Accept = restful.HEADER_Accept
	// Deprecated
	HEADER_ContentType = restful.HEADER_ContentType
	// Deprecated
	HEADER_ContentLength = "Content-Length"
	// Deprecated
	HEADER_AcceptEncoding = restful.HEADER_AcceptEncoding
	// Deprecated
	HEADER_XCsrfToken = "X-CSRF-Token"
	// Deprecated
	HEADER_ApiKey = "api_key"
	// Deprecated
	HEADER_XRequestedWith = "x-requested-with"
)

func ActivateCors(container *restful.Container) {
	cors := newCors(container)
	filter := corsFilter(container, cors)
	container.Filter(filter)
}

func newCors(container *restful.Container) restful.CrossOriginResourceSharing {
	return restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{
			headerNameAuthorization,
			headerNameAccessToken,
			headerNameCacheControl,
			headerNameCurrency,
			headerNameIfModifiedSince,
			headerNameLocale,
			headerNamePragma,
			headerNameContentEncoding,
			headerNameContentType,
			headerNameContentLength,
			headerNameAcceptEncoding,
			headerNameXCsrfToken,
			headerNameApiKey,
			headerNameXRequestedWith,
			headerNameAccept,
			headerNameAccessControlAllowOrigin,
			headerNameAccessControlAllowHeaders,
			headerNameAccessControlAllowCredentials,
		},
		AllowedDomains: []string{"^.*$"},
		Container:      container,
	}
}

func corsFilter(container *restful.Container, cors restful.CrossOriginResourceSharing) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		if req.Request.Header.Get("Origin") != "" {
			resp.AddHeader(headerNameAccessControlAllowOrigin, "*")
		}

		if "OPTIONS" != req.Request.Method {
			chain.ProcessFilter(req, resp)
			return
		}

		if req.Request.Header.Get(headerNameAccessControlRequestMethod) == "" {
			optionsHandler(container, req, resp)
			return
		}

		var corsRecorder = httptest.NewRecorder()
		var corsResponse = restful.NewResponse(corsRecorder)
		cors.Filter(req, corsResponse, chain)

		if corsRecorder.Code == http.StatusNotFound {
			http.NotFound(resp, req.Request)
			return
		}

		// Copy all of the headers
		for k, values := range corsRecorder.Header() {
			for i := 0; i < len(values); i++ {
				if i == 0 {
					resp.Header().Set(k, values[i])
				} else {
					resp.Header().Add(k, values[i])
				}
			}
		}

		// Rewrite the allow header to include OPTIONS
		allowedMethodsResult := corsResponse.Header().Get(headerNameAccessControlAllowMethods)
		allowHeaderValue := getAllowHeader(allowedMethodsResult)
		resp.Header().Set("Allow", allowHeaderValue)

		allowedHeader := strings.Join(cors.AllowedHeaders, ",")

		// Override some headers
		resp.Header().Set("Vary", "Origin")
		resp.Header().Add("Vary", "Access-Control-Request-Method")
		resp.Header().Add("Vary", "Access-Control-Request-Headers")
		resp.Header().Set(headerNameAccessControlAllowMethods, "PATCH,POST,GET,PUT,DELETE,HEAD,OPTIONS,TRACE")
		resp.Header().Set(headerNameAccessControlAllowHeaders, allowedHeader)
		resp.Header().Set(headerNameContentEncoding, "application/json")

		resp.WriteHeader(http.StatusOK)
	}
}

// optionsHandler handles non-CORS OPTIONS responses
func optionsHandler(container *restful.Container, req *restful.Request, resp *restful.Response) {
	// Standard options request
	allowedMethods := findAllowedMethods(container, req.Request.URL.Path)
	if len(allowedMethods) > 0 {
		resp.AddHeader("Allow", strings.Join(allowedMethods, ","))
		resp.AddHeader("Cache-Control", "max-age=604800")
		resp.WriteHeader(204)
	} else {
		resp.WriteHeader(404)
	}
}

// getAllowHeader adds OPTIONS to the specified comma-separated list of methods
func getAllowHeader(allowedMethodsResult string) string {
	allowedMethods := make(types.StringSet)
	allowedMethods.AddAll(strings.Split(allowedMethodsResult, ",")...)
	allowedMethods.Add("OPTIONS")
	allowedMethodsValues := allowedMethods.Values()
	sort.Strings(allowedMethodsValues)
	return strings.Join(allowedMethodsValues, ",")
}

// findAllowedMethods returns the list of HTTP methods for which routes are defined in the container.
// OPTIONS will be added if any matching routes are defined.
func findAllowedMethods(c *restful.Container, requested string) []string {
	var results = make(types.StringSet)
	requestedTokens := tokenizePath(requested)
	for _, ws := range c.RegisteredWebServices() {
		for _, r := range ws.Routes() {
			availableTokens := tokenizePath(r.Path)
			if matchesPath(requestedTokens, availableTokens) {
				results.Add(r.Method)
			}
		}
	}
	if len(results) > 0 {
		results.Add("OPTIONS")
	}
	resultValues := results.Values()
	sort.Strings(resultValues)
	return resultValues
}

func matchesPath(requested, available []string) bool {
	if len(requested) != len(available) {
		return false
	}

	for i, a := range available {
		if a[0] != '{' && requested[i] != a {
			return false
		}
	}

	return true
}

func tokenizePath(path string) []string {
	if "/" == path {
		return nil
	}
	return strings.Split(strings.Trim(path, "/"), "/")
}
