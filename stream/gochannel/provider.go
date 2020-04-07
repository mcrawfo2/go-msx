package gochannel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"sync"
)

const (
	providerNameGoChannel = "gochannel"
)

var loggerAdapter = stream.NewWatermillLoggerAdapter(log.NewLogger("watermill.gochannel"))

type Provider struct {
	channels   map[string]*gochannel.GoChannel
	channelMtx sync.Mutex
}

func (p *Provider) channel(cfg *config.Config, key string, streamBinding *stream.BindingConfiguration) (channel *gochannel.GoChannel, err error) {
	p.channelMtx.Lock()
	defer p.channelMtx.Unlock()

	if channel, ok := p.channels[key]; ok {
		return channel, nil
	}

	bindingConfig, err := NewBindingConfigurationFromConfig(cfg, key, streamBinding)
	if err != nil {
		return
	}

	gochannelConfig := gochannel.Config{
		OutputChannelBuffer:            bindingConfig.Producer.OutputChannelBuffer,
		Persistent:                     bindingConfig.Producer.Persistent,
		BlockPublishUntilSubscriberAck: bindingConfig.Producer.BlockPublishUntilSubscriberAck,
	}

	channel = gochannel.NewGoChannel(
		gochannelConfig,
		loggerAdapter)

	p.channels[key] = channel

	return channel, err
}

func (p *Provider) NewPublisher(cfg *config.Config, name string, streamBinding *stream.BindingConfiguration) (stream.Publisher, error) {
	channel, err := p.channel(cfg, name, streamBinding)
	if err != nil {
		return nil, err
	}

	var publisher = stream.NewTopicPublisher(channel, streamBinding)

	// Do not close the go channel when we are done with the publisher
	publisher = stream.NewIntransientPublisher(publisher)

	return publisher, nil
}

func (p *Provider) NewSubscriber(cfg *config.Config, name string, streamBinding *stream.BindingConfiguration) (stream.Subscriber, error) {
	channel, err := p.channel(cfg, name, streamBinding)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func RegisterProvider(_ *config.Config) error {
	stream.RegisterProvider(providerNameGoChannel, &Provider{
		channels: make(map[string]*gochannel.GoChannel),
	})
	return nil
}
