package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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
	logger                = log.NewLogger("msx.stream")
	listenerMux           sync.Mutex
	listeners             = make(map[string][]ListenerAction)
	router                *message.Router
	routerLogger          = log.NewLogger("watermill.router")
	routerWatermillLogger = NewWatermillLoggerAdapter(routerLogger)
	handlerCounter        atomic.Int32

	ErrRouterRunning      = errors.New("Router already running")
	ErrTopicNotSpecified  = errors.New("Topic not specified")
	ErrActionNotSpecified = errors.New("ListenerAction not specified")
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
	registered := 0
	for topic, topicListeners := range listeners {
		for _, topicListener := range topicListeners {
			err = addListener(ctx, topic, topicListener)
			switch err {
			case ErrBinderNotEnabled, ErrConsumerNotEnabled:
				// Ignore
			case nil:
				registered++
			default:
				return err
			}
		}
	}

	if registered == 0 {
		logger.WithContext(ctx).Warn("No listeners registered.  Router disabled.")
		router = nil
		return nil
	}

	listeners = nil
	defer listenerMux.Unlock()

	var exited = make(chan struct{})
	var finished = false
	go func() {
		ctx = trace.UntracedContextFromContext(ctx)
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

func addListener(ctx context.Context, topic string, action ListenerAction) error {
	subscriber, err := NewSubscriber(ctx, topic)
	if err != nil {
		return err
	}

	bindingConfig, err := NewBindingConfigurationFromConfig(config.FromContext(ctx), topic)
	if err != nil {
		return err
	}

	index := handlerCounter.Inc()
	handlerName := fmt.Sprintf("%s-%d", topic, index)
	router.AddNoPublisherHandler(handlerName, topic, subscriber, listenerHandler(topic, action, bindingConfig))
	return nil
}

func AddListener(topic string, action ListenerAction) error {
	listenerMux.Lock()
	defer listenerMux.Unlock()

	var err error
	if listeners == nil {
		err = ErrRouterRunning
	} else if action == nil {
		return ErrActionNotSpecified
	} else if topic == "" {
		return ErrTopicNotSpecified
	}
	if err != nil {
		return errors.Wrap(err, "Failed to add listener")
	}

	if _, ok := listeners[topic]; !ok {
		listeners[topic] = []ListenerAction{}
	}

	listeners[topic] = append(listeners[topic], action)

	return nil
}

func listenerHandler(topic string, action ListenerAction, cfg *BindingConfiguration) message.NoPublishHandlerFunc {
	traceAction := TraceActionInterceptor(cfg, StatsActionInterceptor(cfg, func(msg *message.Message) error {
		entry := logger.
			WithContext(msg.Context()).
			WithField("messageId", msg.UUID).
			WithField("topic", topic)
		if cfg.LogMessages {
			entry.Infof("received message payload: %s", string(msg.Payload))
		} else {
			entry.Info("received message (payload hidden)")
		}

		retryableAction := func() error {
			return action(msg)
		}

		if err := retry.NewRetry(msg.Context(), cfg.Retry).Retry(retryableAction); err != nil {
			bt := types.BackTraceFromError(err)
			logger.
				WithContext(msg.Context()).
				WithError(err).
				WithField(log.FieldStack, bt.Stanza()).
				Error("Failed to process message")
			log.Stack(logger, msg.Context(), bt)
		}

		msg.Ack()

		return nil
	}))

	result := PanicRecovererActionInterceptor(cfg, traceAction)

	return message.NoPublishHandlerFunc(result)
}
