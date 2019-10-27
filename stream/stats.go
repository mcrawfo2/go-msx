package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"time"
)

const (
	statsGaugeConnections               = "%s.connections"
	statsGaugePublishers                = "%s.publishers.%s"
	statsTimerPublisher                 = "%s.publisher.timer"
	statsCounterPublisherMessages       = "%s.publisherMessages"
	statsCounterPublisherMessageErrors  = "%s.publisherMessageErrors"
	statsGaugeSubscribers               = "%s.subscribers.%s"
	statsTimerSubscriber                = "%s.subscriber.timer"
	statsCounterSubscriberMessages      = "%s.subscriberMessages"
	statsCounterSubscriberMessageErrors = "%s.subscriberMessageErrors"
)

type StatsPublisher struct {
	publisher Publisher
	cfg       *BindingConfiguration
}

func (s *StatsPublisher) Publish(messages ...*message.Message) error {
	now := time.Now()
	count := int64(len(messages))
	stats.Incr(fmt.Sprintf(statsCounterPublisherMessages, s.cfg.Binder), count)
	defer func() {
		stats.PrecisionTiming(stats.Name(fmt.Sprintf(statsTimerPublisher, s.cfg.Binder), s.cfg.Destination, ""), time.Since(now))
	}()

	err := s.publisher.Publish(messages...)
	if err != nil {
		stats.Incr(fmt.Sprintf(statsCounterPublisherMessageErrors, s.cfg.Binder), count)
	}
	return err
}

func (s *StatsPublisher) Close() error {
	stats.GaugeDelta(fmt.Sprintf(statsGaugeConnections, s.cfg.Binder), -1)
	stats.GaugeDelta(fmt.Sprintf(statsGaugePublishers, s.cfg.Binder, s.cfg.Destination), -1)
	return s.publisher.Close()
}

func NewStatsPublisher(publisher Publisher, cfg *BindingConfiguration) *StatsPublisher {
	stats.GaugeDelta(fmt.Sprintf(statsGaugeConnections, cfg.Binder), 1)
	stats.GaugeDelta(fmt.Sprintf(statsGaugePublishers, cfg.Binder, cfg.Destination), 1)
	return &StatsPublisher{
		publisher: publisher,
		cfg:       cfg,
	}
}

type StatsSubscriber struct {
	subscriber message.Subscriber
	cfg        *BindingConfiguration
}

func (s *StatsSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	result, err := s.subscriber.Subscribe(ctx, topic)
	return result, err
}

func (s *StatsSubscriber) Close() error {
	stats.GaugeDelta(fmt.Sprintf(statsGaugeConnections, s.cfg.Binder), -1)
	stats.GaugeDelta(fmt.Sprintf(statsGaugeSubscribers, s.cfg.Binder, s.cfg.Destination), -1)
	return s.subscriber.Close()
}

func NewStatsSubscriber(subscriber message.Subscriber, cfg *BindingConfiguration) *StatsSubscriber {
	stats.GaugeDelta(fmt.Sprintf(statsGaugeConnections, cfg.Binder), 1)
	stats.GaugeDelta(fmt.Sprintf(statsGaugeSubscribers, cfg.Binder, cfg.Destination), 1)
	return &StatsSubscriber{
		subscriber: subscriber,
		cfg:        cfg,
	}
}

type StatsSubscriberAction struct {
	action ListenerAction
	cfg    *BindingConfiguration
}

func (a *StatsSubscriberAction) Call(msg *message.Message) (err error) {
	now := time.Now()
	stats.Incr(fmt.Sprintf(statsCounterSubscriberMessages, a.cfg.Binder), 1)
	defer func() {
		stats.PrecisionTiming(stats.Name(fmt.Sprintf(statsTimerSubscriber, a.cfg.Binder), a.cfg.Destination, ""), time.Since(now))
	}()

	if err = a.action(msg); err != nil {
		stats.Incr(fmt.Sprintf(statsCounterSubscriberMessageErrors, a.cfg.Binder), 1)
	}

	return err
}

func DecorateSubscriberAction(action ListenerAction, cfg *BindingConfiguration) ListenerAction {
	statsAction := &StatsSubscriberAction{
		action: action,
		cfg:    cfg,
	}
	return statsAction.Call
}
