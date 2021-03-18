package streamtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TopicReceiveTestCall func(*testing.T, context.Context) error

type contextTopicKey int

const contextTopicReceiveTestKey contextTopicKey = iota

func ContextWithTopicReceiveTest(ctx context.Context, test *TopicReceiveTest) context.Context {
	return context.WithValue(ctx, contextTopicReceiveTestKey, test)
}

func TopicReceiveTestFromContext(ctx context.Context) *TopicReceiveTest {
	return ctx.Value(contextTopicReceiveTestKey).(*TopicReceiveTest)
}

type TopicReceiveTest struct {
	TopicName string
	Action    stream.ListenerAction
	Metadata  map[string]string
	Payload   []byte
	Receive   struct {
		Want bool
		Got  bool
	}
	Error struct {
		Want bool
		Got  bool
	}
}

func (t *TopicReceiveTest) WithTopic(name string) *TopicReceiveTest {
	t.TopicName = name
	return t
}

func (t *TopicReceiveTest) WithAction(action stream.ListenerAction) *TopicReceiveTest {
	t.Action = action
	return t
}

func (t *TopicReceiveTest) WithPayload(payload []byte) *TopicReceiveTest {
	t.Payload = payload
	return t
}

func (t *TopicReceiveTest) WithMetaData(metadata map[string]string) *TopicReceiveTest {
	t.Metadata = metadata
	return t
}

func (t *TopicReceiveTest) WithWantReceive(want bool) *TopicReceiveTest {
	t.Receive.Want = want
	return t
}

func (t *TopicReceiveTest) WithWantError(want bool) *TopicReceiveTest {
	t.Error.Want = want
	return t
}

func (t *TopicReceiveTest) Received() {
	t.Receive.Got = true
}

func (t *TopicReceiveTest) Test(tt *testing.T) {
	ctx := context.Background()
	ctx = ContextWithTopicReceiveTest(ctx, t)

	msg := NewMessage()
	msg.SetContext(ctx)
	msg.Metadata = t.Metadata

	if len(t.Payload) > 0 {
		msg.Payload = t.Payload
	} else {
		msg.Payload = []byte("{}")
	}

	err := t.Action(msg)
	t.Error.Got = err != nil
	assert.Equal(tt, t.Error.Got, t.Error.Want)
	assert.Equal(tt, t.Receive.Got, t.Receive.Want)
}

func NewTopicReceiveTest() *TopicReceiveTest {
	return &TopicReceiveTest{}
}

func NewMessage() *message.Message {
	uuid := watermill.NewUUID()
	return message.NewMessage(uuid, []byte("{}"))
}
