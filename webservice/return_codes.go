package webservice

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

var DefaultSuccessEnvelope = ResponseEnvelope{}

type RouteBuilderFunc func(*restful.RouteBuilder)

func Returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", DefaultSuccessEnvelope)
}
func Returns201(b *restful.RouteBuilder) {
	b.Returns(http.StatusCreated, "Created", DefaultSuccessEnvelope)
}
func Returns204(b *restful.RouteBuilder) {
	b.Returns(http.StatusNoContent, "No Content", DefaultSuccessEnvelope)
}
func Returns400(b *restful.RouteBuilder) {
	b.Returns(http.StatusBadRequest, "Bad Request", nil)
}
func Returns401(b *restful.RouteBuilder) {
	b.Returns(http.StatusUnauthorized, "Not Authorized", nil)
}
func Returns403(b *restful.RouteBuilder) {
	b.Returns(http.StatusForbidden, "Forbidden", nil)
}
func Returns404(b *restful.RouteBuilder) {
	b.Returns(http.StatusNotFound, "Not Found", nil)
}
func Returns409(b *restful.RouteBuilder) {
	b.Returns(http.StatusConflict, "Conflict", nil)
}
func Returns424(b *restful.RouteBuilder) {
	b.Returns(http.StatusFailedDependency, "Failed Dependency", nil)
}
func Returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Internal Server Error", nil)
}
func Returns503(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bad Gateway", nil)
}

func Returns(statuses ...int) RouteBuilderFunc {
	var statusFuncs []RouteBuilderFunc
	for _, status := range statuses {
		switch status {
		case 200:
			statusFuncs = append(statusFuncs, Returns200)
		case 201:
			statusFuncs = append(statusFuncs, Returns201)
		case 204:
			statusFuncs = append(statusFuncs, Returns204)
		case 400:
			statusFuncs = append(statusFuncs, Returns400)
		case 401:
			statusFuncs = append(statusFuncs, Returns401)
		case 404:
			statusFuncs = append(statusFuncs, Returns404)
		case 409:
			statusFuncs = append(statusFuncs, Returns409)
		case 424:
			statusFuncs = append(statusFuncs, Returns424)
		case 500:
			statusFuncs = append(statusFuncs, Returns500)
		case 503:
			statusFuncs = append(statusFuncs, Returns503)
		}
	}

	return func(b *restful.RouteBuilder) {
		for _, statusFunc := range statusFuncs {
			statusFunc(b)
		}
	}
}
