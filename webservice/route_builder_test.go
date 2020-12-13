package webservice

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/contexttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type RouteBuilderTest struct {
	Request struct {
		Method          string
		Path            string
		QueryParameters url.Values
		Headers         http.Header
	}
	RouteBuilder *restful.RouteBuilder
	Funcs        []RouteBuilderFunc
	Filters      []restful.FilterFunction
	Target       restful.RouteFunction
	Context      context.Context
	Injectors    []types.ContextInjector
	Checks       struct {
		Route    webservicetest.RouteCheck
		Request  webservicetest.RequestCheck
		Response webservicetest.ResponseCheck
		Context  contexttest.ContextCheck
		Log      []log.Check
	}
	Errors struct {
		Context  []error
		Request  []error
		Route    []error
		Response []error
		Log      []error
	}
	Recording *log.Recording
}

func (r *RouteBuilderTest) WithRecording(rec *log.Recording) *RouteBuilderTest {
	r.Recording = rec
	return r
}

func (r *RouteBuilderTest) WithRouteFilter(f restful.FilterFunction) *RouteBuilderTest {
	r.Filters = append(r.Filters, f)
	return r
}

func (r *RouteBuilderTest) WithRouteBuilder(rb *restful.RouteBuilder) *RouteBuilderTest {
	r.RouteBuilder = rb
	return r
}

func (r *RouteBuilderTest) WithRouteTarget(target restful.RouteFunction) *RouteBuilderTest {
	r.Target = target
	return r
}

func (r *RouteBuilderTest) WithRouteTargetReturn(status int) *RouteBuilderTest {
	return r.WithRouteTarget(func(request *restful.Request, response *restful.Response) {
		response.WriteHeader(status)
	})
}

func (r *RouteBuilderTest) WithRouteBuilderDo(fn RouteBuilderFunc) *RouteBuilderTest {
	r.Funcs = append(r.Funcs, fn)
	return r
}

func (r *RouteBuilderTest) WithContext(ctx context.Context) *RouteBuilderTest {
	r.Context = ctx
	return r
}

func (r *RouteBuilderTest) WithContextInjector(i types.ContextInjector) *RouteBuilderTest {
	r.Injectors = append(r.Injectors, i)
	return r
}

func (r *RouteBuilderTest) WithRequestMethod(m string) *RouteBuilderTest {
	r.Request.Method = m
	return r
}

func (r *RouteBuilderTest) WithRequestPath(p string) *RouteBuilderTest {
	r.Request.Path = p
	return r
}

func (r *RouteBuilderTest) WithRequestQueryParameter(name, value string) *RouteBuilderTest {
	if r.Request.QueryParameters == nil {
		r.Request.QueryParameters = make(url.Values)
	}
	r.Request.QueryParameters.Add(name, value)
	return r
}

func (r *RouteBuilderTest) WithRequestHeader(name, value string) *RouteBuilderTest {
	if r.Request.Headers == nil {
		r.Request.Headers = make(http.Header)
	}
	r.Request.Headers.Add(name, value)
	return r
}

func (r *RouteBuilderTest) WithRoutePredicate(p ...webservicetest.RoutePredicate) *RouteBuilderTest {
	r.Checks.Route.Validators = append(r.Checks.Route.Validators, p...)
	return r
}

func (r *RouteBuilderTest) WithContextPredicate(c contexttest.ContextPredicate) *RouteBuilderTest {
	r.Checks.Context.Validators = append(r.Checks.Context.Validators, c)
	return r
}

func (r *RouteBuilderTest) WithRequestPredicate(p webservicetest.RequestPredicate) *RouteBuilderTest {
	r.Checks.Request.Validators = append(r.Checks.Request.Validators, p)
	return r
}

func (r *RouteBuilderTest) WithResponsePredicate(p webservicetest.ResponsePredicate) *RouteBuilderTest {
	r.Checks.Response.Validators = append(r.Checks.Response.Validators, p)
	return r
}

func (r *RouteBuilderTest) WithLogCheck(l log.Check) *RouteBuilderTest {
	r.Checks.Log = append(r.Checks.Log, l)
	return r
}

func (r *RouteBuilderTest) checkContext(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()
	r.Errors.Context = r.Checks.Context.Check(ctx)
	chain.ProcessFilter(req, resp)
}

func (r *RouteBuilderTest) checkRequest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	r.Errors.Request = r.Checks.Request.Check(req)
	chain.ProcessFilter(req, resp)
}

func (r *RouteBuilderTest) defaultTarget(_ *restful.Request, resp *restful.Response) {
	// No body required
}

func (r *RouteBuilderTest) Test(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("")

	if r.Recording == nil {
		r.Recording = log.RecordLogging()
	}

	var rb = r.RouteBuilder
	if rb == nil {
		target := r.Target
		if target == nil {
			target = r.defaultTarget
		}

		rb = ws.Method(r.Request.Method).Path(r.Request.Path).To(target)
	}

	for _, fn := range r.Funcs {
		rb.Do(fn)
	}

	for _, filter := range r.Filters {
		rb.Filter(filter)
	}
	rb.Filter(r.checkContext)
	rb.Filter(r.checkRequest)

	// Build the route
	route := rb.Build()

	ctx := r.Context
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = ContextWithRoute(ctx, &route)
	ctx = ContextWithRouteOperation(ctx, "RouteBuilderTest")
	for _, injector := range r.Injectors {
		ctx = injector(ctx)
	}

	// Check the route
	r.Errors.Route = r.Checks.Route.Check(route)

	// Execute the request if required
	if len(r.Checks.Context.Validators) > 0 ||
		len(r.Checks.Request.Validators) > 0 ||
		len(r.Checks.Response.Validators) > 0 {

		httpRequest := new(http.Request)
		httpRequest.Method = r.Request.Method
		httpRequest.URL = new(url.URL)
		httpRequest.URL.Path = r.Request.Path
		httpRequest.URL.RawQuery = r.Request.QueryParameters.Encode()
		httpRequest.Header = r.Request.Headers
		if httpRequest.Header == nil {
			httpRequest.Header = make(http.Header)
		}
		httpRequest = httpRequest.WithContext(ctx)
		req := restful.NewRequest(httpRequest)

		rec := new(httptest.ResponseRecorder)
		rec.Body = new(bytes.Buffer)
		resp := restful.NewResponse(rec)

		// Build the filter chain
		filterChain := restful.FilterChain{
			Filters: route.Filters,
			Target:  route.Function,
		}

		// Execute the request
		filterChain.ProcessFilter(req, resp)

		// Check the response
		r.Errors.Response = r.Checks.Response.Check(rec)
	}

	// Check the logs
	for _, logCheck := range r.Checks.Log {
		errs := logCheck.Check(r.Recording)
		r.Errors.Log = append(r.Errors.Log, errs...)
	}

	// Report any errors
	testhelpers.ReportErrors(t, "Route", r.Errors.Route)
	testhelpers.ReportErrors(t, "Context", r.Errors.Context)
	testhelpers.ReportErrors(t, "Request", r.Errors.Request)
	testhelpers.ReportErrors(t, "Response", r.Errors.Response)
	testhelpers.ReportErrors(t, "Log", r.Errors.Log)
}
