// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package controllertest

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/contexttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/mohae/deepcopy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "cto-github.cisco.com/NFV-BU/go-msx/ops/restops/httperrors"
)

// Allow the testing.T object to be injected via the context.Context
type contextKey int

const contextKeyTesting contextKey = iota

type ControllerFactory func(ctx context.Context) (webservice.RestController, error)

type EndpointProducerSourceFactory func(ctx context.Context) (restops.EndpointsProducer, error)

type ControllerTest struct {
	Server struct {
		ContextPath string
	}
	Request struct {
		Method          string
		Path            string
		QueryParameters url.Values
		Headers         http.Header
		Body            []byte
	}
	Controller struct {
		RootPath          string
		Factory           ControllerFactory
		EndpointsProducer EndpointProducerSourceFactory
	}
	Context struct {
		Base         context.Context
		Config       *config.Config
		TokenDetails *securitytest.MockTokenDetailsProvider
		Injectors    []types.ContextInjector
	}
	Setup  []func()
	Checks struct {
		Route    webservicetest.RouteCheck
		Request  webservicetest.RequestCheck
		Response webservicetest.ResponseCheck
		Context  contexttest.ContextCheck
		Log      []logtest.Check
	}
	Verifiers struct {
		Request  []webservicetest.RequestVerifier
		Response []webservicetest.ResponseVerifier
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

func (r *ControllerTest) Clone() *ControllerTest {
	return deepcopy.Copy(r).(*ControllerTest)
}

func (r *ControllerTest) WithContextPath(contextPath string) *ControllerTest {
	r.Server.ContextPath = contextPath
	return r
}

func (r *ControllerTest) WithRecording(rec *logtest.Recording) *ControllerTest {
	r.Recording = rec
	return r
}

func (r *ControllerTest) WithControllerRootPath(rootPath string) *ControllerTest {
	r.Controller.RootPath = rootPath
	return r
}

func (r *ControllerTest) WithControllerFactory(factory ControllerFactory) *ControllerTest {
	r.Controller.Factory = factory
	return r
}

func (r *ControllerTest) WithEndpointProducerSourceFactory(factory EndpointProducerSourceFactory) *ControllerTest {
	r.Controller.EndpointsProducer = factory
	return r
}

func (r *ControllerTest) WithContext(ctx context.Context) *ControllerTest {
	r.Context.Base = ctx
	return r
}

func (r *ControllerTest) WithContextInjector(i types.ContextInjector) *ControllerTest {
	r.Context.Injectors = append(r.Context.Injectors, i)
	return r
}

func (r *ControllerTest) WithConfig(cfg *config.Config) *ControllerTest {
	r.Context.Config = cfg
	return r
}

func (r *ControllerTest) WithTokenDetailsProvider(provider *securitytest.MockTokenDetailsProvider) *ControllerTest {
	r.Context.TokenDetails = provider
	return r
}

func (r *ControllerTest) WithRequestMethod(m string) *ControllerTest {
	r.Request.Method = m
	return r
}

func (r *ControllerTest) WithRequestPath(p string, vars map[string]string) *ControllerTest {
	if vars != nil {
		for k, v := range vars {
			p = strings.ReplaceAll(p, "{"+k+"}", v)
		}
	}
	r.Request.Path = p
	return r
}

func (r *ControllerTest) WithRequestBody(body []byte) *ControllerTest {
	r.Request.Body = body
	return r
}

func (r *ControllerTest) WithRequestBodyString(body string) *ControllerTest {
	return r.WithRequestBody([]byte(body))
}

func (r *ControllerTest) WithRequestBodyJson(dto interface{}) *ControllerTest {
	b, _ := json.Marshal(dto)
	return r.WithRequestBody(b)
}

func (r *ControllerTest) WithRequestQueryParameter(name, value string) *ControllerTest {
	if r.Request.QueryParameters == nil {
		r.Request.QueryParameters = make(url.Values)
	}
	r.Request.QueryParameters.Add(name, value)
	return r
}

func (r *ControllerTest) WithRequestHeader(name, value string) *ControllerTest {
	if r.Request.Headers == nil {
		r.Request.Headers = make(http.Header)
	}
	r.Request.Headers.Add(name, value)
	return r
}

func (r *ControllerTest) WithRoutePredicate(p ...webservicetest.RoutePredicate) *ControllerTest {
	r.Checks.Route.Validators = append(r.Checks.Route.Validators, p...)
	return r
}

func (r *ControllerTest) WithContextPredicate(c contexttest.ContextPredicate) *ControllerTest {
	r.Checks.Context.Validators = append(r.Checks.Context.Validators, c)
	return r
}

func (r *ControllerTest) WithRequestPredicate(p webservicetest.RequestPredicate) *ControllerTest {
	r.Checks.Request.Validators = append(r.Checks.Request.Validators, p)
	return r
}

func (r *ControllerTest) WithRequestVerifier(fn webservicetest.RequestVerifier) *ControllerTest {
	r.Verifiers.Request = append(r.Verifiers.Request, fn)
	return r
}

func (r *ControllerTest) WithResponsePredicate(p webservicetest.ResponsePredicate) *ControllerTest {
	r.Checks.Response.Validators = append(r.Checks.Response.Validators, p)
	return r
}

func (r *ControllerTest) WithLogCheck(l logtest.Check) *ControllerTest {
	r.Checks.Log = append(r.Checks.Log, l)
	return r
}

func (r *ControllerTest) checkContext(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()
	r.Errors.Context = r.Checks.Context.Check(ctx)
	chain.ProcessFilter(req, resp)
}

func (r *ControllerTest) checkRequest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	r.Errors.Request = r.Checks.Request.Check(req)
	if len(r.Verifiers.Request) > 0 {
		t := req.Request.Context().Value(contextKeyTesting).(*testing.T)
		for _, v := range r.Verifiers.Request {
			v(t, req)
		}
	}
	chain.ProcessFilter(req, resp)
}

func (r *ControllerTest) checkRoute(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()
	route := webservice.RouteFromContext(ctx)
	r.Errors.Context = r.Checks.Route.Check(*route)
	chain.ProcessFilter(req, resp)
}

func (r *ControllerTest) newWebServer(ctx context.Context, t *testing.T) *webservice.WebServer {
	webServerConfig := webservice.WebServerConfig{
		Enabled: true,
		Host:    "0.0.0.0",
		Port:    9999,
		Cors: webservice.CorsConfig{
			Enabled: true,
		},
		ContextPath: r.Server.ContextPath,
		StaticPath:  "/www",
	}
	managementConfig := webservice.ManagementSecurityConfig{}

	server, err := webservice.NewWebServer(&webServerConfig, &managementConfig, ctx)
	assert.NoError(t, err)
	return server
}

func (r *ControllerTest) Test(t *testing.T) {
	var cfg *config.Config
	if r.Context.Config != nil {
		cfg = r.Context.Config
	} else {
		cfg = configtest.NewInMemoryConfig(nil)
	}

	err := fs.ConfigureFileSystem(cfg)
	assert.NoError(t, err)

	if r.Recording == nil {
		r.Recording = logtest.RecordLogging()
	}

	ctx := r.Context.Base
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	ctx = context.WithValue(ctx, contextKeyTesting, t)
	if r.Context.TokenDetails != nil {
		ctx = r.Context.TokenDetails.Inject(ctx)
	}

	for _, injector := range r.Context.Injectors {
		ctx = injector(ctx)
	}

	// These filters will be auto-appended to the route filter chain during the filter chain execution
	ctx = webservice.ContextWithFilters(ctx,
		r.checkContext,
		r.checkRequest,
		r.checkRoute)

	// Create a web server
	s := r.newWebServer(context.Background(), t)
	if r.Controller.Factory != nil {
		var controller webservice.RestController
		controller, err = r.Controller.Factory(ctx)
		assert.NoError(t, err)
		err = s.RegisterRestController(r.Controller.RootPath, controller)
	} else if r.Controller.EndpointsProducer != nil {
		var producer restops.EndpointsProducer
		producer, err = r.Controller.EndpointsProducer(ctx)
		assert.NoError(t, err)
		err = restops.
			NewRestfulEndpointRegisterer(s).
			RegisterEndpoints(producer)
	}
	assert.NoError(t, err)
	ctx = webservice.ContextWithWebServerValue(ctx, s)
	ctx = trace.ContextWithUntracedContext(ctx)
	s.SetContext(ctx)

	// Create the request
	req := new(http.Request)
	req.Method = r.Request.Method

	req.URL = new(url.URL)
	req.URL.Path = r.Request.Path
	if r.Request.Path == "" {
		req.URL.Path = "/"
	}
	if r.Server.ContextPath != "" {
		req.URL.Path = r.Server.ContextPath + req.URL.Path
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
	s.Handler().ServeHTTP(rec, req)

	// Check the response
	r.Errors.Response = r.Checks.Response.Check(rec)
	if len(r.Verifiers.Response) > 0 {
		for _, v := range r.Verifiers.Response {
			v(t, rec)
		}
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
