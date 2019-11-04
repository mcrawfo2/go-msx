package webservice

import (
	"context"
	"crypto/tls"
	"github.com/emicklei/go-restful"
	"net/http"
	"sync"
	"time"
)

type WebServer struct {
	ctx          context.Context
	cfg          *WebServerConfig
	container    *restful.Container
	containerMtx sync.Mutex
	router       *restful.CurlyRouter
	services     []*restful.WebService
	server       *http.Server
}

func NewWebServer(cfg *WebServerConfig, ctx context.Context) *WebServer {
	return &WebServer{
		ctx:    ctx,
		cfg:    cfg,
		router: &restful.CurlyRouter{},
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
	s.container.Filter(s.contextInjectorFilter(s.container, s.router))

	// Add all web services
	for _, service := range s.services {
		s.container.Add(service)
	}

	// Add documentation provider
	if documentationProvider != nil {
		documentationService := new(restful.WebService)
		documentationService.Path(s.cfg.ContextPath)
		if err := documentationProvider.Actuate(s.container, documentationService); err != nil {
			logger.WithError(err).Errorf("Failed to register actuator")
		} else {
			s.container.Add(documentationService)
		}
	}

	// Add actuators
	for _, actuatorProvider := range actuators {
		if actuatorProvider == nil {
			continue
		}

		actuatorService := new(restful.WebService)
		actuatorService.Path(s.cfg.ContextPath)
		if err := actuatorProvider.Actuate(actuatorService); err != nil {
			logger.WithError(err).Errorf("Failed to register actuator")
		} else {
			s.container.Add(actuatorService)
		}
	}

	return s.container
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

func (s *WebServer) contextInjectorFilter(container *restful.Container, router restful.RouteSelector) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// Inject the container and router
		ctx := ContextWithContainer(s.ctx, container)
		ctx = ContextWithRouter(ctx, router)

		// Inject the webservice and route
		service, route, _ := router.SelectRoute(container.RegisteredWebServices(), req.Request)
		ctx = ContextWithService(ctx, service)
		ctx = ContextWithRoute(ctx, route)
		if route != nil {
			ctx = ContextWithRouteOperation(ctx, route.Operation)
		} else {
			ctx = ContextWithRouteOperation(ctx, "unknown")
		}

		// Inject the request variables
		ctx = ContextWithHttpRequest(ctx, req.Request)
		ctx = ContextWithPathParameters(ctx, req.PathParameters())

		// Inject security provider
		ctx = ContextWithSecurityProvider(ctx, securityProvider)

		// Inject anything from the registered injectors
		ctx = injectContextValues(ctx)

		// Add the context into the request
		req.Request = req.Request.WithContext(ctx)

		chain.ProcessFilter(req, resp)
	}
}
