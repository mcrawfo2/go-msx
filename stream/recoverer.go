package stream

import (
	"runtime/debug"

	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type PanicRecovererSubscriberAction struct {
	action ListenerAction
	cfg    *BindingConfiguration
}

func (a *PanicRecovererSubscriberAction) Call(msg *message.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var e error
			if err, ok := r.(error); ok {
				e = err
			} else {
				e = errors.Errorf("Exception: %v", r)
			}
			logger.WithContext(msg.Context()).WithError(e).Error("Recovered from panic")
			bt := types.BackTraceFromDebugStackTrace(debug.Stack())
			log.Stack(logger, msg.Context(), bt)

			msg.Ack()

			err = nil
		}
	}()

	err = a.action(msg)
	if err != nil {
		logger.WithContext(msg.Context()).WithError(err).Error("Failed to process message")
	}
	return err
}

func PanicRecovererActionInterceptor(cfg *BindingConfiguration, action ListenerAction) ListenerAction {
	panicRecovererAction := &PanicRecovererSubscriberAction{
		action: action,
		cfg:    cfg,
	}
	return panicRecovererAction.Call
}
