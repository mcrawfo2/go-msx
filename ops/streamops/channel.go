// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

// Channel maps to asyncapi.ChannelItem
type Channel struct {
	name        string
	binding     *stream.BindingConfiguration
	documentors ops.Documentors[Channel]
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) Destination() string {
	return c.binding.Destination
}

func (c *Channel) DefaultContentType() string {
	return c.binding.ContentType
}

func (c *Channel) DefaultContentEncoding() string {
	return c.binding.ContentEncoding
}

func (c *Channel) Binder() string {
	return c.binding.Binder
}

func (c *Channel) WithDocumentor(d ops.Documentor[Channel]) *Channel {
	c.documentors = c.documentors.WithDocumentor(d)
	return c
}

func (c Channel) Documentor(pred ops.DocumentorPredicate[Channel]) ops.Documentor[Channel] {
	return c.documentors.Documentor(pred)
}

func NewChannel(ctx context.Context, name string) (*Channel, error) {
	binding, err := stream.NewBindingConfiguration(ctx, name)
	if err != nil {
		return nil, err
	}

	result := &Channel{
		name:    name,
		binding: binding,
	}

	RegisterChannel(result)

	return result, nil
}

// Registry

var registeredChannels = make(map[string]*Channel)

func RegisterChannel(c *Channel) {
	registeredChannels[c.Name()] = c
}

func RegisteredChannels() map[string]*Channel {
	return registeredChannels
}

func RegisteredChannel(channel string) *Channel {
	return registeredChannels[channel]
}
