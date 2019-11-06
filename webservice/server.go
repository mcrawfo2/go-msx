package webservice

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"net/http"
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
	security      SecurityProvider
	actuators     []ServiceProvider
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

func (s *WebServer) NewService() *restful.WebService {
	s.resetContainer()
	webService := new(restful.WebService)
	webService.Path(s.cfg.ContextPath)
	s.services = append(s.services, webService)
	return webService
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

	// Add all web services
	for _, service := range s.services {
		s.container.Add(service)
	}

	// Add documentation provider
	s.actuateDocumentation(s.documentation)

	// Add actuators
	for _, provider := range s.actuators {
		s.actuateService(provider)
	}

	return s.container
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
	s.ctx = ctx

	s.server = &http.Server{
		Addr:    s.cfg.Address(),
		Handler: s.Handler(),
	}

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

func (s *WebServer) SetSecurityProvider(provider SecurityProvider) {
	if provider != nil {
		s.security = provider
	}
}

func (s *WebServer) RegisterInjector(injector types.ContextInjector) {
	s.injectors.Register(injector)
}

func requestContextInjectorFilter(ctx context.Context, container *restful.Container, router restful.RouteSelector, security SecurityProvider, injectors *types.ContextInjectors) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// Inject the container, router, security provider, request
		ctx := ContextWithContainer(ctx, container)
		ctx = ContextWithRouter(ctx, router)
		ctx = ContextWithSecurityProvider(ctx, security)

		// Inject the webservice and route
		service, route, _ := router.SelectRoute(container.RegisteredWebServices(), req.Request)
		ctx = ContextWithService(ctx, service)
		ctx = ContextWithRoute(ctx, route)
		if route != nil {
			ctx = ContextWithRouteOperation(ctx, route.Operation)
		} else {
			ctx = ContextWithRouteOperation(ctx, "unknown")
		}

		// Execute external injectors
		ctx = injectors.Inject(ctx)

		req.Request = req.Request.WithContext(ctx)

		chain.ProcessFilter(req, resp)
	}
}
