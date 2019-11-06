package kafka

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/opentracing/opentracing-go"
)

var (
	logger = log.NewLogger("msx.stream.kafka")
)

type TraceMarshaler struct {
	upstream kafka.DefaultMarshaler
}

func (t TraceMarshaler) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	producerMessage, err := t.upstream.Marshal(topic, msg)

	span := trace.SpanFromContext(msg.Context())

	err = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(msg.Metadata))
	if err != nil {
		logger.WithError(err).WithContext(msg.Context()).Warn("Failed to apply trace context to outgoing message")
	}

	for k, v := range msg.Metadata {
		producerMessage.Headers = append(producerMessage.Headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	return producerMessage, nil
}

func (t TraceMarshaler) Unmarshal(consumerMessage *sarama.ConsumerMessage) (*message.Message, error) {
	msg, err := t.upstream.Unmarshal(consumerMessage)
	if err != nil {
		return msg, err
	}

	for _, header := range consumerMessage.Headers {
		if string(header.Key) != kafka.UUIDHeaderKey {
			msg.Metadata.Set(string(header.Key), string(header.Value))
		}
	}

	return msg, err
}
