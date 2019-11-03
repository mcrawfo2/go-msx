package webservice

import "github.com/emicklei/go-restful"

const (
	MIME_XML  = "application/xml"  // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_JSON = "application/json" // Accept or Content-Type used in Consumes() and/or Produces()

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

	METHOD_GET     = "GET"
	METHOD_POST    = "POST"
	METHOD_PUT     = "PUT"
	METHOD_DELETE  = "DELETE"
	METHOD_HEAD    = "HEAD"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_PATCH   = "PATCH"
)

func ActivateCORS(container *restful.Container) {
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
		AllowedMethods: []string{METHOD_GET, METHOD_POST, METHOD_PUT, METHOD_DELETE, METHOD_OPTIONS, METHOD_HEAD, METHOD_PATCH},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}

	container.Filter(cors.Filter)
	container.Filter(OPTIONSFilter)
}

// OPTIONSFilter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// As for any filter, you can also install it for a particular WebService within a Container.
// Note: this filter is not needed when using CrossOriginResourceSharing (for CORS).
func OPTIONSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}

	resp.AddHeader(HEADER_AccessControlAllowOrigin, "*")
	resp.AddHeader(HEADER_AccessControlAllowMethods, "PATCH, POST, GET, OPTIONS, PUT, DELETE")
	resp.AddHeader(HEADER_AccessControlRequestHeaders, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	resp.AddHeader(HEADER_AccessControlAllowHeaders, "Authorization, access_token, cache-control, currency, if-modified-since, locale, pragma, content-type, content-length")
	resp.AddHeader(HEADER_ContentEncoding, "application/json")
}
