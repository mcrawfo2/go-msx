package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/ThreeDotsLabs/watermill/message"
	"time"
)

const (
	statsSubsystemKafka                 = "stream"
	statsGaugeConnections               = "connections"
	statsGaugePublisherCount            = "publisher_count"
	statsCounterPublisherSends          = "publisher_sends"
	statsCounterPublisherSendErrors     = "publisher_send_errors"
	statsHistogramPublisherSendTime     = "publisher_send_time"
	statsGaugeSubscriberCount           = "subscriber_count"
	statsCounterSubscriberReceives      = "subscriber_receives"
	statsCounterSubscriberReceiveErrors = "subscriber_receive_errors"
	statsHistogramSubscriberReceiveTime = "subscriber_receive_time"
)

var (
	gaugeConnections = stats.NewGauge(statsSubsystemKafka, statsGaugeConnections)

	gaugeVecPublishers            = stats.NewGaugeVec(statsSubsystemKafka, statsGaugePublisherCount, "binder", "topic")
	counterVecPublisherSends      = stats.NewGaugeVec(statsSubsystemKafka, statsCounterPublisherSends, "binder", "topic")
	counterVecPublisherSendErrors = stats.NewGaugeVec(statsSubsystemKafka, statsCounterPublisherSendErrors, "binder", "topic")
	histVecPublisherSendTime      = stats.NewHistogramVec(statsSubsystemKafka, statsHistogramPublisherSendTime, nil, "binder", "topic")

	gaugeVecSubscribers               = stats.NewGaugeVec(statsSubsystemKafka, statsGaugeSubscriberCount, "binder", "topic")
	counterVecSubscriberReceives      = stats.NewGaugeVec(statsSubsystemKafka, statsCounterSubscriberReceives, "binder", "topic")
	counterVecSubscriberReceiveErrors = stats.NewGaugeVec(statsSubsystemKafka, statsCounterSubscriberReceiveErrors, "binder", "topic")
	histVecSubscriberReceiveTime      = stats.NewHistogramVec(statsSubsystemKafka, statsHistogramSubscriberReceiveTime, nil, "binder", "topic")
)

type StatsPublisher struct {
	publisher Publisher
	cfg       *BindingConfiguration
}

func (s *StatsPublisher) Publish(message *message.Message) error {
	now := time.Now()
	count := float64(1)
	counterVecPublisherSends.WithLabelValues(s.cfg.Binder, s.cfg.Destination).Add(count)
	defer func() {
		histVecPublisherSendTime.WithLabelValues(s.cfg.Binder, s.cfg.Destination).Observe(float64(time.Since(now)) / float64(time.Millisecond))
	}()

	err := s.publisher.Publish(message)
	if err != nil {
		counterVecPublisherSendErrors.WithLabelValues(s.cfg.Binder, s.cfg.Destination).Add(count)
	}
	return err
}

func (s *StatsPublisher) Close() error {
	gaugeConnections.Dec()
	gaugeVecPublishers.WithLabelValues(s.cfg.Binder, s.cfg.Destination).Dec()
	return s.publisher.Close()
}

func NewStatsPublisher(publisher Publisher, cfg *BindingConfiguration) *StatsPublisher {
	gaugeConnections.Inc()
	gaugeVecPublishers.WithLabelValues(cfg.Binder, cfg.Destination).Inc()
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
	gaugeConnections.Dec()
	gaugeVecSubscribers.WithLabelValues(s.cfg.Binder, s.cfg.Destination).Dec()
	return s.subscriber.Close()
}

func NewStatsSubscriber(subscriber message.Subscriber, cfg *BindingConfiguration) *StatsSubscriber {
	gaugeConnections.Inc()
	gaugeVecSubscribers.WithLabelValues(cfg.Binder, cfg.Destination).Inc()
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
	counterVecSubscriberReceives.WithLabelValues(a.cfg.Binder, a.cfg.Destination).Inc()
	defer func() {
		histVecSubscriberReceiveTime.WithLabelValues(a.cfg.Binder, a.cfg.Destination).Observe(float64(time.Since(now)) / float64(time.Millisecond))
	}()

	if err = a.action(msg); err != nil {
		counterVecSubscriberReceiveErrors.WithLabelValues(a.cfg.Binder, a.cfg.Destination).Inc()
	}

	return err
}

func StatsActionInterceptor(cfg *BindingConfiguration, action ListenerAction) ListenerAction {
	statsAction := &StatsSubscriberAction{
		action: action,
		cfg:    cfg,
	}
	return statsAction.Call
}
