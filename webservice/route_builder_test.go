package webservice

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/contexttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type routeBuilderContextKey int

const routeBuilderContextKeyTesting routeBuilderContextKey = iota

type RouteBuilderTest struct {
	Request struct {
		Method          string
		Path            string
		QueryParameters url.Values
		Headers         http.Header
		Body            []byte
	}
	WebService   *restful.WebService
	RouteBuilder *restful.RouteBuilder
	Route        struct {
		Path       string
		Parameters []*restful.Parameter
		Funcs      []RouteBuilderFunc
		Filters    []restful.FilterFunction
		Target     restful.RouteFunction
	}
	Context   context.Context
	Injectors []types.ContextInjector
	Checks    struct {
		Route    webservicetest.RouteCheck
		Request  webservicetest.RequestCheck
		Response webservicetest.ResponseCheck
		Context  contexttest.ContextCheck
		Log      []logtest.Check
	}
	Verifiers struct {
		Request []webservicetest.RequestVerifier
	}
	Errors struct {
		Context  []error
		Request  []error
		Route    []error
		Response []error
		Log      []error
	}
	Recording *logtest.Recording
}

func (r *RouteBuilderTest) WithRecording(rec *logtest.Recording) *RouteBuilderTest {
	r.Recording = rec
	return r
}

func (r *RouteBuilderTest) WithRouteFilter(f restful.FilterFunction) *RouteBuilderTest {
	r.Route.Filters = append(r.Route.Filters, f)
	return r
}

func (r *RouteBuilderTest) WithRouteBuilder(rb *restful.RouteBuilder) *RouteBuilderTest {
	r.RouteBuilder = rb
	return r
}

func (r *RouteBuilderTest) WithWebService(ws *restful.WebService) *RouteBuilderTest {
	r.WebService = ws
	return r
}

func (r *RouteBuilderTest) WithRouteParameter(p *restful.Parameter) *RouteBuilderTest {
	r.Route.Parameters = append(r.Route.Parameters, p)
	return r
}

func (r *RouteBuilderTest) WithRouteTarget(target restful.RouteFunction) *RouteBuilderTest {
	r.Route.Target = target
	return r
}

func (r *RouteBuilderTest) WithRouteTargetReturn(status int) *RouteBuilderTest {
	return r.WithRouteTarget(func(request *restful.Request, response *restful.Response) {
		response.WriteHeader(status)
	})
}

func (r *RouteBuilderTest) WithRoutePath(p string) *RouteBuilderTest {
	r.Route.Path = p
	return r
}

func (r *RouteBuilderTest) WithRouteBuilderDo(fn RouteBuilderFunc) *RouteBuilderTest {
	r.Route.Funcs = append(r.Route.Funcs, fn)
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
	if r.Route.Path == "" {
		r.Route.Path = p
	}
	return r
}

func (r *RouteBuilderTest) WithRequestBody(body []byte) *RouteBuilderTest {
	r.Request.Body = body
	return r
}

func (r *RouteBuilderTest) WithRequestBodyString(body string) *RouteBuilderTest {
	return r.WithRequestBody([]byte(body))
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

func (r *RouteBuilderTest) WithRequestVerifier(fn webservicetest.RequestVerifier) *RouteBuilderTest {
	r.Verifiers.Request = append(r.Verifiers.Request, fn)
	return r
}

func (r *RouteBuilderTest) WithResponsePredicate(p webservicetest.ResponsePredicate) *RouteBuilderTest {
	r.Checks.Response.Validators = append(r.Checks.Response.Validators, p)
	return r
}

func (r *RouteBuilderTest) WithLogCheck(l logtest.Check) *RouteBuilderTest {
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
	if len(r.Verifiers.Request) > 0 {
		t := req.Request.Context().Value(routeBuilderContextKeyTesting).(*testing.T)
		for _, v := range r.Verifiers.Request {
			v(t, req)
		}
	}
	chain.ProcessFilter(req, resp)
}

func (r RouteBuilderTest) defaultTarget(_ *restful.Request, resp *restful.Response) {
	// No body required
}

func (r *RouteBuilderTest) Test(t *testing.T) {
	ws := r.WebService
	if ws == nil {
		ws = new(restful.WebService)
		ws.Path("")
	}

	if r.Recording == nil {
		r.Recording = logtest.RecordLogging()
	}

	var rb = r.RouteBuilder
	if rb == nil {
		target := r.Route.Target
		if target == nil {
			target = r.defaultTarget
		}

		method := r.Request.Method
		if method == "" {
			method = "POST"
		}

		path := r.Route.Path
		if path == "" {
			path = "/"
		}

		rb = ws.Method(method).Path(r.Route.Path).To(target)
	}

	//rb.Operation(r.Route.Operation)

	for _, p := range r.Route.Parameters {
		rb.Param(p)
	}

	for _, fn := range r.Route.Funcs {
		rb.Do(fn)
	}

	for _, filter := range r.Route.Filters {
		rb.Filter(filter)
	}
	rb.Filter(r.checkContext)
	rb.Filter(r.checkRequest)

	// Build the route
	route := rb.Build()
	ws.Route(rb)

	ctx := r.Context
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, routeBuilderContextKeyTesting, t)

	ctx = ContextWithRoute(ctx, &route)
	ctx = ContextWithRouteOperation(ctx, route.Operation)
	for _, injector := range r.Injectors {
		ctx = injector(ctx)
	}

	// Check the route
	r.Errors.Route = r.Checks.Route.Check(route)

	// Execute the request if required
	if len(r.Checks.Context.Validators) > 0 ||
		len(r.Checks.Request.Validators) > 0 ||
		len(r.Checks.Response.Validators) > 0 ||
		len(r.Verifiers.Request) > 0 {

		// Create a container for dispatching
		container := restful.NewContainer()
		container.Router(new(restful.CurlyRouter))
		container.Add(ws)

		// Create the request
		req := new(http.Request)
		req.Method = route.Method

		req.URL = new(url.URL)
		req.URL.Path = r.Request.Path
		if r.Request.Path == "" {
			req.URL.Path = "/"
		}
		req.URL.RawQuery = r.Request.QueryParameters.Encode()

		req.Header = r.Request.Headers
		if req.Header == nil {
			req.Header = make(http.Header)
		}

		if r.Request.Body != nil {
			bodyBuffer := bytes.NewBuffer(r.Request.Body)
			bodyReadCloser := ioutil.NopCloser(bodyBuffer)
			req.Body = bodyReadCloser
			req.ContentLength = int64(len(r.Request.Body))
		}

		req = req.WithContext(ctx)

		// Create the response
		rec := new(httptest.ResponseRecorder)
		rec.Body = new(bytes.Buffer)

		// Dispatch the request
		container.Dispatch(rec, req)

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
