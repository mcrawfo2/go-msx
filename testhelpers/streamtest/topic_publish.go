package streamtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type TopicPublishTestCall func(*testing.T, context.Context) error

type TopicPublishTest struct {
	TopicName    string
	PublishError error
	CallError    error
	Call         TopicPublishTestCall
	Checks       struct {
		Message MessageCheck
	}
	Errors struct {
		Message []error
	}
}

func (t *TopicPublishTest) WithTopic(name string) *TopicPublishTest {
	t.TopicName = name
	return t
}

func (t *TopicPublishTest) WithPublishError(err error) *TopicPublishTest {
	t.PublishError = err
	return t
}

func (t *TopicPublishTest) WithPredicate(predicate MessagePredicate) *TopicPublishTest {
	t.Checks.Message.Validators = append(t.Checks.Message.Validators, predicate)
	return t
}

func (t *TopicPublishTest) WithCallError(err error) *TopicPublishTest {
	t.CallError = err
	return t
}

func (t *TopicPublishTest) WithCall(call TopicPublishTestCall) *TopicPublishTest {
	t.Call = call
	return t
}

func (t *TopicPublishTest) Test(tt *testing.T) {
	mockPublisher := new(stream.MockPublisher)

	mockPublisher.
		On("Publish", mock.AnythingOfType("*message.Message")).
		Run(func(args mock.Arguments) {
			msg := args.Get(0).(*message.Message)
			t.Errors.Message = t.Checks.Message.Check(msg)
		}).
		Return(t.PublishError)

	mockPublisher.
		On("Close").
		Return(nil)

	mockProvider := new(stream.MockProvider)

	mockProvider.
		On("NewPublisher",
			mock.AnythingOfType("*config.Config"),
			t.TopicName,
			mock.AnythingOfType("*stream.BindingConfiguration")).
		Return(mockPublisher, nil)

	stream.RegisterProvider("mock", mockProvider)

	ctx := context.Background()
	ctx = configtest.ContextWithNewInMemoryConfig(ctx, map[string]string{
		"spring.application.name":                                 tt.Name(),
		"spring.cloud.stream.bindings." + t.TopicName + ".binder": "mock",
	})

	err := t.Call(tt, ctx)
	if t.CallError != nil {
		assert.Error(tt, err)
		assert.True(tt, errors.Is(err, t.CallError))
	} else if t.PublishError != nil {
		assert.Error(tt, err)
		assert.True(tt, errors.Is(err, t.PublishError))
	}

	testhelpers.ReportErrors(tt, "Message", t.Errors.Message)
}

func NewTopicPublishTest() *TopicPublishTest {
	return &TopicPublishTest{}
}
