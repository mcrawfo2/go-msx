package webservice

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

const (
	MIME_XML              = restful.MIME_XML  // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_JSON             = restful.MIME_JSON // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_TEXT_PLAIN       = "text/plain"      // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_APPLICATION_FORM = "application/x-www-form-urlencoded"
	MIME_MULTIPART_FORM   = "multipart/form-data"

	HEADER_ContentEncoding             = "Content-Encoding"
	HEADER_AccessControlRequestHeaders = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods   = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin    = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowHeaders   = "Access-Control-Allow-Headers"

	HEADER_Authorization  = "Authorization"
	HEADER_Accept         = "Accept"
	HEADER_ContentType    = "Content-Type"
	HEADER_ContentLength  = "Content-Length"
	HEADER_AcceptEncoding = "Accept-Encoding"
	HEADER_XCsrfToken     = "X-CSRF-Token"
	HEADER_ApiKey         = "api_key"
	HEADER_XRequestedWith = "x-requested-with"
)

func ActivateCors(container *restful.Container) {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders: []string{"X-CFI-Header"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin: *",
			HEADER_AccessControlAllowHeaders,
			HEADER_XRequestedWith,
			HEADER_ApiKey,
			HEADER_Accept,
			HEADER_ContentType,
			HEADER_ContentLength,
			HEADER_AcceptEncoding,
			HEADER_XCsrfToken,
			HEADER_Authorization,
			"Access-Control-Allow-Credentials: true"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}

	container.Filter(cors.Filter)
	container.Filter(corsOptionsFilter)
}

// corsOptionsFilter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// As for any filter, you can also install it for a particular WebService within a Container.
// Note: this filter is not needed when using CrossOriginResourceSharing (for CORS).
func corsOptionsFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}

	resp.AddHeader(HEADER_AccessControlAllowOrigin, "*")
	resp.AddHeader(HEADER_AccessControlAllowMethods, "PATCH,POST,GET,PUT,DELETE,HEAD,OPTIONS,TRACE")
	resp.AddHeader(HEADER_AccessControlRequestHeaders, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	resp.AddHeader(HEADER_AccessControlAllowHeaders, "Authorization, access_token, cache-control, currency, if-modified-since, locale, pragma, content-type, content-length")
	resp.AddHeader(HEADER_ContentEncoding, "application/json")
}
