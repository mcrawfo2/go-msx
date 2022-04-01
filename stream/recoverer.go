// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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

			bt := types.BackTraceFromDebugStackTrace(debug.Stack())
			logger.
				WithContext(msg.Context()).
				WithError(e).
				WithField(log.FieldStack, bt.Stanza()).
				Error("Recovered from panic")
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
