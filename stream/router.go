package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
)

type ListenerAction func(msg *message.Message) error

var (
	logger = log.NewLogger("msx.stream")
	listenerMux sync.Mutex
	listeners = make(map[string][]ListenerAction)
	router *message.Router
	routerLogger = log.NewLogger("watermill.router")
	routerWatermillLogger = NewWatermillLoggerAdapter(log.NewLogger("watermill.router"))
	handlerCounter atomic.Int32
)

func StartRouter(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return errors.New("Failed to retrieve config from context")
	}

	routerConfig := message.RouterConfig{}

	var err error
	router, err = message.NewRouter(routerConfig, routerWatermillLogger)
	if err != nil {
		return errors.New("Failed to create router")
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.Recoverer)

	listenerMux.Lock()
	for topic, topicListeners := range listeners {
		for _, topicListener := range topicListeners {
			if err = addListener(cfg, topic, topicListener); err != nil {
				return err
			}
		}
	}
	listeners = nil
	defer listenerMux.Unlock()

	var exited = make(chan struct{})
	var finished = false
	go func() {
		err = router.Run(ctx)
		close(exited)
		if finished && err != nil {
			routerLogger.WithError(err).Error("Router exited abnormally")
		}
		router = nil
	}()

	// Wait until the router has started or exited
waiter:
	for {
		select {
		case <-router.Running():
			break waiter

		case <-exited:
			break waiter
		}
	}

	if err != nil {
		return errors.Wrap(err, "Failed to run router")
	}

	return nil
}

func StopRouter(context.Context) error {
	if router != nil {
		return router.Close()
	}
	return nil
}

func addListener(cfg *config.Config, topic string, action ListenerAction) error {
	subscriber, err := NewSubscriber(cfg, topic)
	if err != nil {
		return err
	}

	bindingConfig, err := NewBindingConfigurationFromConfig(cfg, topic)
	if err != nil {
		return err
	}

	index := handlerCounter.Inc()
	handlerName := fmt.Sprintf("%s-%d", topic, index)
	router.AddNoPublisherHandler(handlerName, topic, subscriber, listenerHandler(action, bindingConfig))
	return nil
}

func AddListener(topic string, action ListenerAction) error {
	listenerMux.Lock()
	defer listenerMux.Unlock()

	if listeners == nil {
		return errors.New("Router already running")
	}

	if _, ok := listeners[topic]; !ok {
		listeners[topic] = []ListenerAction{}
	}

	listeners[topic] = append(listeners[topic], action)

	return nil
}

func listenerHandler(action ListenerAction, cfg *BindingConfiguration) message.NoPublishHandlerFunc {
	action = TraceActionDecorator(cfg, StatsActionDecorator(action, cfg))

	return func(msg *message.Message) error {
		retryableAction := func() error {
			return action(msg)
		}

		if err := cfg.Retry.Retry(retryableAction); err != nil {
			logger.WithContext(msg.Context()).WithError(err).Error("Failed to process message")
		}

		msg.Ack()

		return nil
	}
}
