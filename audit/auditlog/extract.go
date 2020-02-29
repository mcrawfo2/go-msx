package auditlog

import (
	"github.com/emicklei/go-restful"
	"strconv"
)

const (
	XForwardedForHeader = "X-Forwarded-For"
)

type RequestDetails struct {
	Source   string
	Protocol string
	Host     string
	Port     string
}

func ExtractRequestDetails(req *restful.Request, host string, port int) *RequestDetails {
	remoteAddr := req.Request.RemoteAddr
	proxyAddr := req.Request.Header.Get(XForwardedForHeader)
	if proxyAddr != "" {
		remoteAddr = proxyAddr
	}

	return &RequestDetails{
		Source:   remoteAddr,
		Protocol: req.Request.Proto,
		Host:     host,
		Port:     strconv.Itoa(port),
	}
}
