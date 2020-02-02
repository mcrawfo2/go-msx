package webservice

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type WebServer struct {
	ctx           context.Context
	cfg           *WebServerConfig
	container     *restful.Container
	containerMtx  sync.Mutex
	router        *restful.CurlyRouter
	services      []*restful.WebService
	documentation DocumentationProvider
	security      AuthenticationProvider
	actuators     []ServiceProvider
	aliases       []StaticAlias
	server        *http.Server
	injectors     *types.ContextInjectors
}

func NewWebServer(cfg *WebServerConfig, ctx context.Context) *WebServer {
	return &WebServer{
		ctx:       ctx,
		cfg:       cfg,
		router:    &restful.CurlyRouter{},
		injectors: new(types.ContextInjectors),
	}
}

func (s *WebServer) NewService(path string) (*restful.WebService, error) {
	if s == nil {
		return nil, ErrDisabled
	}

	if path == "" {
		return nil, errors.New("Web service path must be specified")
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	s.resetContainer()
	webService := new(restful.WebService)
	webService.Path(s.cfg.ContextPath + path)
	s.services = append(s.services, webService)
	return webService, nil
}

func (s *WebServer) resetContainer() {
	s.containerMtx.Lock()
	defer s.containerMtx.Unlock()

	s.container = nil
}

func (s *WebServer) Handler() http.Handler {
	s.containerMtx.Lock()
	defer s.containerMtx.Unlock()

	if s.container != nil {
		return s.container
	}

	s.container = restful.NewContainer()
	s.container.Router(s.router)
	s.container.Filter(requestContextInjectorFilter(
		s.ctx,
		s.container,
		s.router,
		s.security,
		s.injectors.Clone()))
	s.container.Filter(tracingFilter)
	s.container.Filter(recoveryFilter)
	s.container.Filter(optionsFilter)
	s.container.Filter(securityContextFilter)
	s.container.Filter(authenticationFilter)
	if s.cfg.Cors {
		ActivateCors(s.container)
	}

	// Add all web services
	for _, service := range s.services {
		s.container.Add(service)
	}

	// Add documentation provider
	s.actuateDocumentation(s.documentation)

	// Add services
	for _, provider := range s.actuators {
		s.actuateService(provider)
	}

	// Add static file server
	s.actuateStatic(s.aliases)

	return s.container
}

func (s *WebServer) actuateStatic(aliases []StaticAlias) {
	staticService := new(restful.WebService)
	staticService.Path(s.cfg.ContextPath)

	logger.Infof("Serving static files at %s", s.cfg.Url())

	// Add NOT FOUND for unclaimed paths of other services
	for _, webService := range s.container.RegisteredWebServices() {
		webServiceRoot := webService.RootPath()
		staticService.Route(staticService.
			GET(webServiceRoot + "/{subPath:*}").
			To(HttpHandlerController(http.NotFound)))
	}

	staticUiHandler := http.StripPrefix(
		staticService.RootPath(),
		http.FileServer(http.Dir(s.cfg.StaticPath))).ServeHTTP

	for _, alias := range aliases {
		staticService.Route(staticService.
			GET(alias.Path).
			Operation("static-alias").
			To(StaticFile(alias.File)))
	}

	staticService.Route(staticService.
		GET("").
		Operation("static-root").
		To(EnsureSlash(staticUiHandler)))

	staticService.Route(staticService.
		GET("/{subPath:*}").
		Operation("static-file").
		To(HttpHandlerController(staticUiHandler)))

	s.container.Add(staticService)
}

func (s *WebServer) actuateDocumentation(provider DocumentationProvider) {
	if provider == nil {
		return
	}

	documentationService := new(restful.WebService)
	documentationService.Path(s.cfg.ContextPath)
	if err := s.documentation.Actuate(s.container, documentationService); err != nil {
		logger.WithError(err).Errorf("Failed to register actuator")
	} else {
		s.container.Add(documentationService)
	}
}

func (s *WebServer) actuateService(provider ServiceProvider) {
	if provider == nil {
		return
	}

	actuatorService := new(restful.WebService)
	actuatorService.Path(s.cfg.ContextPath)
	if err := provider.Actuate(actuatorService); err != nil {
		logger.WithError(err).Errorf("Failed to register actuator")
	} else {
		s.container.Add(actuatorService)
	}
}

func (s *WebServer) Serve(ctx context.Context) error {
	s.ctx = trace.UntracedContextFromContext(ctx)

	s.server = &http.Server{
		Addr:    s.cfg.Address(),
		Handler: s.Handler(),
	}

	restful.EnableTracing(s.cfg.TraceEnabled)

	// Start the server
	go func() {
		var err error
		if s.cfg.Tls {
			logger.Infof("Serving on https://%s%s", s.cfg.Address(), s.cfg.ContextPath)

			tlsConfig := &tls.Config{
				ClientAuth: tls.VerifyClientCertIfGiven,
			}
			tlsConfig.BuildNameToCertificate()
			err = s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		} else {
			logger.Infof("Serving on http://%s%s", s.cfg.Address(), s.cfg.ContextPath)

			err = s.server.ListenAndServe()
		}

		if err != nil && err.Error() != "http: Server closed" {
			logger.WithError(err).Error("Web Server exited")
		} else {
			logger.WithError(err).Info("Web server exited")
		}
	}()

	return nil
}

func (s *WebServer) StopServing(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *WebServer) SetDocumentationProvider(provider DocumentationProvider) {
	if provider != nil {
		s.documentation = provider
	}
}

func (s *WebServer) RegisterActuator(provider ServiceProvider) {
	if provider != nil {
		s.actuators = append(s.actuators, provider)
	}
}

func (s *WebServer) SetAuthenticationProvider(provider AuthenticationProvider) {
	if provider != nil {
		s.security = provider
	}
}

func (s *WebServer) RegisterInjector(injector types.ContextInjector) {
	s.injectors.Register(injector)
}

func (s *WebServer) RegisterAlias(path, file string) {
	s.aliases = append(s.aliases, StaticAlias{
		Path: path,
		File: file,
	})
}

func (s *WebServer) StaticFolder() string {
	folder, err := filepath.Abs(s.cfg.StaticPath)
	if err != nil {
		return s.cfg.StaticPath
	}
	return folder
}

func (s *WebServer) Url() string {
	return s.cfg.Url()
}

func (s *WebServer) ContextPath() string {
	return s.cfg.ContextPath
}

func requestContextInjectorFilter(ctx context.Context, container *restful.Container, router restful.RouteSelector, security AuthenticationProvider, injectors *types.ContextInjectors) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// Inject the container, router, security provider, request
		ctx2 := ContextWithContainer(ctx, container)
		ctx2 = ContextWithRouter(ctx2, router)
		ctx2 = ContextWithSecurityProvider(ctx2, security)

		// Inject the webservice and route
		service, route, _ := router.SelectRoute(container.RegisteredWebServices(), req.Request)
		ctx2 = ContextWithService(ctx2, service)
		ctx2 = ContextWithRoute(ctx2, route)
		if route != nil {
			ctx2 = ContextWithRouteOperation(ctx2, route.Operation)
		} else {
			ctx2 = ContextWithRouteOperation(ctx2, "unknown")
		}

		// Execute external injectors
		ctx2 = injectors.Inject(ctx2)

		req.Request = req.Request.WithContext(ctx2)

		chain.ProcessFilter(req, resp)
	}
}
