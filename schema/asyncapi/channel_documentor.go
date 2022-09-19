// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"github.com/pkg/errors"
)

var ErrNilChannel = errors.New("No channel specified to document.")

type ChannelDocumentor struct {
	skip        bool
	channelItem *ChannelItem
	mutator     ops.DocumentElementMutator[ChannelItem]
}

func (d *ChannelDocumentor) WithSkip(skip bool) *ChannelDocumentor {
	d.skip = skip
	return d
}

func (d *ChannelDocumentor) WithChannelItem(channelItem *ChannelItem) *ChannelDocumentor {
	d.channelItem = channelItem
	return d
}

func (d *ChannelDocumentor) WithChannelItemMutator(fn ops.DocumentElementMutator[ChannelItem]) *ChannelDocumentor {
	d.mutator = fn
	return d
}

func (d ChannelDocumentor) DocType() string {
	return DocType
}

func (d ChannelDocumentor) Document(c *streamops.Channel) error {
	if c == nil {
		return errors.Wrap(ErrNilChannel, "")
	}

	if d.skip {
		return nil
	}

	// Initialize
	channelItem := d.channelItem
	if channelItem == nil {
		channelItem = new(ChannelItem)
	}

	// Servers
	if err := new(ServerDocumentor).Document(c.Binder()); err != nil {
		return err
	}
	channelItem.WithServers(c.Binder())

	// Bindings
	channelBindingsEns(channelItem)

	// Mutator
	if d.mutator != nil {
		d.mutator(channelItem)
	}

	// Publish
	documentationReflector.SpecEns().WithChannelsItem(c.Name(), *channelItem)

	// Children
	if publisher := streamops.RegisteredChannelPublisher(c.Name()); publisher != nil {
		err := ops.DocumentorWithType[streamops.ChannelPublisher](publisher, DocType).
			OrElse(ChannelPublisherDocumentor{}).
			Document(publisher)
		if err != nil {
			return err
		}
	}

	if subscriber := streamops.RegisteredChannelSubscriber(c.Name()); subscriber != nil {
		err := ops.DocumentorWithType[streamops.ChannelSubscriber](subscriber, DocType).
			OrElse(ChannelSubscriberDocumentor{}).
			Document(subscriber)
		if err != nil {
			return err
		}
	}

	return nil
}
