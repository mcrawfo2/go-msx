// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"sync"
)

// ChannelSubscriber maps to asyncapi.Operation
type ChannelSubscriber struct {
	channel          *Channel
	name             string
	dispatchHeader   types.Optional[string]
	dispatcher       stream.Dispatcher
	dispatchTable    map[stream.MetadataHeader]stream.ListenerAction
	dispatchTableMtx sync.RWMutex
	documentors      ops.Documentors[ChannelSubscriber]
}

func (p *ChannelSubscriber) Channel() *Channel {
	return p.channel
}

func (p *ChannelSubscriber) Name() string {
	return p.name
}

func (p *ChannelSubscriber) Documentor(predicate ops.DocumentorPredicate[ChannelSubscriber]) ops.Documentor[ChannelSubscriber] {
	return p.documentors.Documentor(predicate)
}

func (p *ChannelSubscriber) AddDocumentor(d ...ops.Documentor[ChannelSubscriber]) *ChannelSubscriber {
	p.documentors = p.documentors.WithDocumentor(d...)
	return p
}

func (p *ChannelSubscriber) AddMessageConsumer(mc MessageConsumer) (err error) {
	p.dispatchTableMtx.Lock()
	defer p.dispatchTableMtx.Unlock()

	if p.dispatchHeader.IsPresent() {
		var dispatchValues []string
		dispatchValues, err = mc.MetadataFilterValues(p.dispatchHeader.Value())
		if err != nil {
			return err
		}

		for _, v := range dispatchValues {
			p.dispatchTable[stream.MetadataHeader(v)] = mc.OnMessage
		}

		p.dispatcher, err = stream.NewMetadataDispatcherIndirect(p.dispatchHeader.Value(), p.lookupListenerAction)
		if err != nil {
			return err
		}
	} else if p.dispatcher != nil {
		return errors.New("Cannot register multiple message consumers on the same channel without setting a dispatch header")
	} else {
		p.dispatcher = messageConsumerDispatcher{messageConsumer: mc}
	}

	return nil
}

func (p *ChannelSubscriber) lookupListenerAction(value string) stream.ListenerAction {
	p.dispatchTableMtx.RLock()
	defer p.dispatchTableMtx.RUnlock()

	return p.dispatchTable[stream.MetadataHeader(value)]
}

func (p *ChannelSubscriber) OnMessage(msg *message.Message) error {
	if p.dispatcher == nil {
		return errors.Errorf("No consumers registered for channel %q", p.name)
	}
	err := p.dispatcher.Dispatch(msg)
	if err != nil {
		err = errors.Wrapf(err, "Message dispatch failed for channel %q", p.name)
	}

	return err
}

func NewChannelSubscriber(_ context.Context, channel *Channel, name string, dispatchHeader types.Optional[string]) (*ChannelSubscriber, error) {
	if channel == nil {
		return nil, errors.Errorf("Nil channel passed to subscriber %q", name)
	}

	result := &ChannelSubscriber{
		channel:        channel,
		name:           name,
		dispatchHeader: dispatchHeader,
		dispatchTable:  map[stream.MetadataHeader]stream.ListenerAction{},
	}

	if RegisterChannelSubscriber(result) {
		if err := stream.AddMessageListener(result.Channel().Name(), result); err != nil {
			return nil, errors.Wrapf(err, "Failed to listen on channel %q", channel.Name())
		}
	}

	return result, nil
}

var registeredChannelSubscribers = make(map[string]*ChannelSubscriber)

func RegisterChannelSubscriber(p *ChannelSubscriber) bool {
	// Do not add registration twice
	if _, ok := registeredChannelSubscribers[p.Channel().Name()]; ok {
		return false
	}

	registeredChannelSubscribers[p.channel.Name()] = p
	return true
}

func RegisteredChannelSubscriber(channel string) *ChannelSubscriber {
	return registeredChannelSubscribers[channel]
}
