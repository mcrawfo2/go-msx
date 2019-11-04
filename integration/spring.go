package integration

import "net/http"

type SpringHttpStatus struct {
	code int
	name string
}

func (s SpringHttpStatus) Code() int {
	return s.code
}

func (s SpringHttpStatus) Name() string {
	return s.name
}

var springHttpStatusByCode map[int]SpringHttpStatus
var springHttpStatusByName map[string]SpringHttpStatus

func init() {
	springHttpStatusByCode = make(map[int]SpringHttpStatus)
	springHttpStatusByName = make(map[string]SpringHttpStatus)

	setSpringHttpStatusMaps(http.StatusContinue, "CONTINUE")
	setSpringHttpStatusMaps(http.StatusSwitchingProtocols, "SWITCHING_PROTOCOLS")
	setSpringHttpStatusMaps(http.StatusProcessing, "PROCESSING")
	setSpringHttpStatusMaps(103, "CHECKPOINT")

	setSpringHttpStatusMaps(http.StatusOK, "OK")
	setSpringHttpStatusMaps(http.StatusCreated, "CREATED")
	setSpringHttpStatusMaps(http.StatusAccepted, "ACCEPTED")
	setSpringHttpStatusMaps(http.StatusNonAuthoritativeInfo, "NON_AUTHORITATIVE_INFORMATION")
	setSpringHttpStatusMaps(http.StatusNoContent, "NO_CONTENT")
	setSpringHttpStatusMaps(http.StatusResetContent, "RESET_CONTENT")
	setSpringHttpStatusMaps(http.StatusPartialContent, "PARTIAL_CONTENT")
	setSpringHttpStatusMaps(http.StatusMultiStatus, "MULTI_STATUS")
	setSpringHttpStatusMaps(http.StatusAlreadyReported, "ALREADY_REPORTED")
	setSpringHttpStatusMaps(http.StatusIMUsed, "IM_USED")

	setSpringHttpStatusMaps(http.StatusMultipleChoices, "MULTIPLE_CHOICES")
	setSpringHttpStatusMaps(http.StatusMovedPermanently, "MOVED_PERMANENTLY")
	setSpringHttpStatusMaps(http.StatusFound, "FOUND")
	setSpringHttpStatusMaps(http.StatusSeeOther, "SEE_OTHER")
	setSpringHttpStatusMaps(http.StatusNotModified, "NOT_MODIFIED")
	setSpringHttpStatusMaps(http.StatusUseProxy, "USE_PROXY")
	setSpringHttpStatusMaps(http.StatusTemporaryRedirect, "TEMPORARY_REDIRECT")
	setSpringHttpStatusMaps(http.StatusPermanentRedirect, "PERMANENT_REDIRECT")

	setSpringHttpStatusMaps(http.StatusBadRequest, "BAD_REQUEST")
	setSpringHttpStatusMaps(http.StatusUnauthorized, "UNAUTHORIZED")
	setSpringHttpStatusMaps(http.StatusPaymentRequired, "PAYMENT_REQUIRED")
	setSpringHttpStatusMaps(http.StatusForbidden, "FORBIDDEN")
	setSpringHttpStatusMaps(http.StatusNotFound, "NOT_FOUND")
	setSpringHttpStatusMaps(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
	setSpringHttpStatusMaps(http.StatusNotAcceptable, "NOT_ACCEPTABLE")

	setSpringHttpStatusMaps(http.StatusProxyAuthRequired, "PROXY_AUTHENTICATION_REQUIRED")
	setSpringHttpStatusMaps(http.StatusRequestTimeout, "REQUEST_TIMEOUT")
	setSpringHttpStatusMaps(http.StatusConflict, "CONFLICT")
	setSpringHttpStatusMaps(http.StatusGone, "GONE")
	setSpringHttpStatusMaps(http.StatusLengthRequired, "LENGTH_REQUIRED")
	setSpringHttpStatusMaps(http.StatusPreconditionFailed, "PRECONDITION_FAILED")
	setSpringHttpStatusMaps(http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE")
	setSpringHttpStatusMaps(http.StatusRequestURITooLong, "URI_TOO_LONG")
	setSpringHttpStatusMaps(http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE")
	setSpringHttpStatusMaps(http.StatusRequestedRangeNotSatisfiable, "REQUESTED_RANGE_NOT_SATISFIABLE")
	setSpringHttpStatusMaps(http.StatusExpectationFailed, "EXPECTATION_FAILED")
	setSpringHttpStatusMaps(http.StatusTeapot, "I_AM_A_TEAPOT")
	setSpringHttpStatusMaps(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY")
	setSpringHttpStatusMaps(http.StatusLocked, "LOCKED")
	setSpringHttpStatusMaps(http.StatusFailedDependency, "FAILED_DEPENDENCY")
	setSpringHttpStatusMaps(http.StatusUpgradeRequired, "UPGRADE_REQUIRED")
	setSpringHttpStatusMaps(http.StatusPreconditionRequired, "PRECONDITION_REQUIRED")
	setSpringHttpStatusMaps(http.StatusTooManyRequests, "TOO_MANY_REQUESTS")
	setSpringHttpStatusMaps(http.StatusRequestHeaderFieldsTooLarge, "REQUEST_HEADER_FIELDS_TOO_LARGE")
	setSpringHttpStatusMaps(http.StatusUnavailableForLegalReasons, "UNAVAILABLE_FOR_LEGAL_REASONS")

	setSpringHttpStatusMaps(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR")
	setSpringHttpStatusMaps(http.StatusNotImplemented, "NOT_IMPLEMENTED")
	setSpringHttpStatusMaps(http.StatusBadGateway, "BAD_GATEWAY")
	setSpringHttpStatusMaps(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE")
	setSpringHttpStatusMaps(http.StatusGatewayTimeout, "GATEWAY_TIMEOUT")
	setSpringHttpStatusMaps(http.StatusHTTPVersionNotSupported, "HTTP_VERSION_NOT_SUPPORTED")
	setSpringHttpStatusMaps(http.StatusVariantAlsoNegotiates, "VARIANT_ALSO_NEGOTIATES")
	setSpringHttpStatusMaps(http.StatusInsufficientStorage, "INSUFFICIENT_STORAGE")
	setSpringHttpStatusMaps(http.StatusLoopDetected, "LOOP_DETECTED")
	setSpringHttpStatusMaps(509, "BANDWIDTH_LIMIT_EXCEEDED")
	setSpringHttpStatusMaps(http.StatusNotExtended, "NOT_EXTENDED")
	setSpringHttpStatusMaps(http.StatusNetworkAuthenticationRequired, "NETWORK_AUTHENTICATION_REQUIRED")
}

func setSpringHttpStatusMaps(code int, name string) {
	springHttpStatus := SpringHttpStatus{code: code, name: name}
	springHttpStatusByCode[code] = springHttpStatus
	springHttpStatusByName[name] = springHttpStatus
}

func GetSpringStatusNameForCode(code int) string {
	return springHttpStatusByCode[code].Name()
}

func GetSpringStatusCodeForName(name string) int {
	return springHttpStatusByName[name].Code()
}
