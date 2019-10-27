package stream

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}
