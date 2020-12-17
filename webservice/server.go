package webservice

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/background"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"path"
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
	handlers      map[string]http.Handler
	documentation []DocumentationProvider
	security      AuthenticationProvider
	actuators     []ServiceProvider
	actuatorCfg   *ManagementSecurityConfig
	aliases       []StaticAlias
	server        *http.Server
	injectors     *types.ContextInjectors
	webRoot       http.FileSystem
}

// NewService returns an existing restful.WebService if one exists at the specified path.
// Otherwise, it registers a new restful.WebService and returns it.
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

	for _, service := range s.services {
		if service.RootPath() == s.cfg.ContextPath+path {
			return service, nil
		}
	}

	s.resetContainer()
	webService := new(restful.WebService)
	webService.Path(s.cfg.ContextPath + path)

	s.services = append(s.services, webService)
	return webService, nil
}

// SetHandler registers an http.Handler at the specified path
func (s *WebServer) SetHandler(p string, handler http.Handler) {
	p = path.Join(s.cfg.ContextPath, p)
	s.handlers[p] = handler
}

// AddDocumentationProvider registers a documentation provider
func (s *WebServer) AddDocumentationProvider(provider DocumentationProvider) {
	if provider != nil {
		s.documentation = append(s.documentation, provider)
	}
}

// RegisterActuator registers an (admin) endpoint actuator
func (s *WebServer) RegisterActuator(provider ServiceProvider) {
	if provider != nil {
		s.actuators = append(s.actuators, provider)
	}
}

// RegisterRestController registers a RestController
func (s *WebServer) RegisterRestController(path string, controller RestController) (err error) {
	svc, err := s.NewService(path)
	if err != nil {
		return err
	}
	controller.Routes(svc)
	return nil
}

// SetAuthenticationProvider registers the specified AuthenticationProvider
func (s *WebServer) SetAuthenticationProvider(provider AuthenticationProvider) {
	if provider != nil {
		s.security = provider
	}
}

// RegisterInjector registers a context injector
func (s *WebServer) RegisterInjector(injector types.ContextInjector) {
	s.injectors.Register(injector)
}

// RegisterAlias registers a static file alias
func (s *WebServer) RegisterAlias(path, file string) {
	s.aliases = append(s.aliases, StaticAlias{
		Path: path,
		File: file,
	})
}

// Url returns the base URL of the web server
func (s *WebServer) Url() string {
	return s.cfg.Url()
}

// ContextPath returns the context path (prefix) of the web server
func (s *WebServer) ContextPath() string {
	return s.cfg.ContextPath
}

// Handler returns the root http.Handler for the server, building a restful.Container if necessary.
func (s *WebServer) Handler() http.Handler {
	return s.generateContainer()
}

func (s *WebServer) resetContainer() {
	s.containerMtx.Lock()
	defer s.containerMtx.Unlock()

	s.container = nil
}

func (s *WebServer) generateContainer() *restful.Container {
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
		s.injectors.Slice()))
	s.container.Filter(tracingFilter)
	s.container.Filter(recoveryFilter)
	if s.cfg.Cors {
		ActivateCors(s.container)
	}
	s.container.Filter(securityContextFilter)
	s.container.Filter(authenticationFilter)
	s.container.Filter(auditContextFilter)

	// Add all web services
	for _, svc := range s.services {
		s.container.Add(svc)
	}

	// Add all handlers
	for p, handler := range s.handlers {
		s.container.HandleWithFilter(p, handler)
	}

	// Add documentation provider
	for _, documentation := range s.documentation {
		s.activateDocumentation(documentation)
	}

	// Add admin actuators
	for _, provider := range s.actuators {
		s.activateActuator(provider)
	}

	// Add static file server
	s.activateStatic(s.aliases)

	return s.container
}

func (s *WebServer) activateStatic(aliases []StaticAlias) {
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

	fs := http.FileServer(noIndexFileSystem{s.webRoot})

	staticUiHandler := http.StripPrefix(
		staticService.RootPath(), fs).ServeHTTP

	for _, alias := range aliases {
		staticService.Route(staticService.
			GET(alias.Path).
			Operation("static-alias").
			To(StaticFileAlias(alias, staticUiHandler)))
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

func (s *WebServer) activateDocumentation(provider DocumentationProvider) {
	documentationService := new(restful.WebService)
	documentationService.Path(s.cfg.ContextPath)
	if err := provider.Actuate(s.container, documentationService); err != nil {
		logger.WithError(err).Errorf("Failed to register actuator")
	} else {
		s.container.Add(documentationService)
	}
}

func (s *WebServer) activateActuator(provider ServiceProvider) {
	if provider == nil {
		return
	}

	actuatorService := new(restful.WebService)
	actuatorService.Path(s.cfg.ContextPath + "/admin/" + provider.EndpointName())

	if s.actuatorCfg.EndpointSecurityEnabled(provider.EndpointName()) {
		securityFilter := NewManagementSecurityFilter(s.actuatorCfg)
		actuatorService.Filter(securityFilter)
	}

	if err := provider.Actuate(actuatorService); err != nil {
		logger.WithError(err).Errorf("Failed to register actuator")
	} else {
		s.container.Add(actuatorService)
	}
}

// Serve starts a web server in the background.
func (s *WebServer) Serve(ctx context.Context) error {
	s.ctx = trace.UntracedContextFromContext(ctx)

	s.server = &http.Server{
		Addr:     s.cfg.Address(),
		Handler:  s.Handler(),
		ErrorLog: stdlog.New(logger.Level(logrus.ErrorLevel).(*log.LevelLogger), "", 0),
	}

	restful.EnableTracing(s.cfg.TraceEnabled)

	// Start the server
	go func() {
		var err error
		if s.cfg.Tls.Enabled {
			var tlsConfig *tls.Config

			logger.Infof("Serving on https://%s%s", s.cfg.Address(), s.cfg.ContextPath)
			tlsConfig, err = buildTlsConfig(s.ctx, s.cfg)
			var ln net.Listener
			if err == nil {
				tlsConfig.BuildNameToCertificate()
				ln, err = s.getTLSListener(tlsConfig)
				if err == nil {
					err = s.server.Serve(ln)
				}
			}

			logger.WithError(err).Error("Error starting TLS listener")
		} else {
			logger.Infof("Serving on http://%s%s", s.cfg.Address(), s.cfg.ContextPath)

			err = s.server.ListenAndServe()
		}

		if err == http.ErrServerClosed {
			logger.Info("Web server exited normally")
			background.ErrorReporterFromContext(ctx).NonFatal(err)
		} else if err != nil {
			logger.Error("Web Server exited abnormally")
			background.ErrorReporterFromContext(ctx).Fatal(err)
		}
	}()

	return nil
}

// StopServing stops the running background webserver.
func (s *WebServer) StopServing(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func requestContextInjectorFilter(ctx context.Context, container *restful.Container, router restful.RouteSelector, security AuthenticationProvider, injectors types.ContextInjectors) restful.FilterFunction {
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

func buildTlsConfig(ctx context.Context, cfg *WebServerConfig) (*tls.Config, error) {
	ca, err := ioutil.ReadFile(cfg.Tls.CaFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read CA certificate file %q", cfg.Tls.CaFile)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	ciphers, err := ParseCiphers(cfg.Tls.CipherSuites)
	if err != nil {
		return nil, err
	}

	w, err := certificate.NewSource(ctx, cfg.Tls.CertificateSource)
	if err != nil {
		return nil, err
	}

	tlsconfig := &tls.Config{
		ClientAuth:     tls.VerifyClientCertIfGiven,
		ClientCAs:      caCertPool,
		MinVersion:     TLSLookup[cfg.Tls.MinVersion],
		CipherSuites:   ciphers,
		GetCertificate: w.TlsCertificate,
	}

	return tlsconfig, nil
}

func (s *WebServer) getTLSListener(tlscfg *tls.Config) (net.Listener, error) {
	addr := s.server.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return ln, err
	}
	return tls.NewListener(ln, tlscfg), nil
}

func NewWebServer(cfg *WebServerConfig, actuatorConfig *ManagementSecurityConfig, ctx context.Context) (*WebServer, error) {
	webRoot, err := NewWebRoot(cfg.StaticPath)
	if err != nil {
		return nil, err
	}

	return &WebServer{
		ctx:         ctx,
		cfg:         cfg,
		actuatorCfg: actuatorConfig,
		router:      &restful.CurlyRouter{},
		injectors:   new(types.ContextInjectors),
		handlers:    make(map[string]http.Handler),
		webRoot:     webRoot,
	}, nil
}
