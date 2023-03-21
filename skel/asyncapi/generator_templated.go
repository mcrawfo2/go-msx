// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/gedex/inflector"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type TemplatedGenerator struct {
	cfg  GeneratorConfig
	dirs types.StringSet
}

// Generate generates templated asyncapi-enabled channels
func (g *TemplatedGenerator) Generate() (err error) {
	logger.Infof("Generating channel %s", g.cfg.ChannelName)

	if g.cfg.Publisher && g.cfg.Subscriber {
		logger.Warn("Generating loop-back topic (publisher and subscriber both specified)")
	}

	if err = g.RenderChannel(); err != nil {
		return err
	}

	return g.GoGenerate()
}

func (g *TemplatedGenerator) Dirs() []string {
	return g.dirs.Values()
}

// GoGenerate runs the generation step on any template target directories
func (g *TemplatedGenerator) GoGenerate() (err error) {
	if err = GoGenerate(g.Dirs()); err != nil {
		return
	}

	g.dirs = make(types.StringSet)
	return nil
}

// RenderChannel generates the channel implementation
func (g *TemplatedGenerator) RenderChannel() (err error) {
	opts := skel.NewRenderOptions()
	opts.AddVariables(map[string]string{
		"async.channel":         g.cfg.ChannelName,
		"async.channel.package": packageName(channelShortName(g.cfg.ChannelName)),
	})

	if err = g.render(opts, OperationNone, msgChannel); err != nil {
		return
	}

	if !g.cfg.Deep {
		return
	}

	if g.cfg.Publisher {
		err = g.RenderChannelOperation(OperationPublish)
		if err != nil {
			return
		}
	}

	if g.cfg.Subscriber {
		err = g.RenderChannelOperation(OperationSubscribe)
		if err != nil {
			return
		}
	}

	return
}

// RenderChannelOperation generates the channel operation (publisher/subscriber) implementation
func (g *TemplatedGenerator) RenderChannelOperation(op string) (err error) {
	opts := skel.NewRenderOptions()
	opts.AddVariables(map[string]string{
		"async.channel":         g.cfg.ChannelName,
		"async.channel.package": packageName(channelShortName(g.cfg.ChannelName)),
		"async.operation":       op,
		"async.operation.id":    operationName(channelShortName(g.cfg.ChannelName), op),
	})

	opts.AddCondition("CHANNEL_MULTI", g.cfg.Multi)

	if err = g.render(opts, op, msgOperation); err != nil {
		return err
	}

	if !g.cfg.Deep {
		return
	}

	// for each new message type
	for _, messageId := range g.cfg.Messages {
		if op == OperationPublish {
			messageId += "Response"
		} else {
			messageId += "Request"
		}
		err = g.RenderMessageOperation(op, messageId, nil)
		if err != nil {
			return
		}
	}

	return
}

var messageOperationSnippets = map[string]map[string]string{
	OperationPublish: {
		"imports": `
	api "${app.packageurl}/internal/stream/${async.channel.package}/api"
`,
		"dependencies": `
type ${async.upmsgtype}Publisher interface {
	Publish${async.upmsgtype}(ctx context.Context, payload api.${async.upmsgtype}) error
}

`,
		"implementation": `
type ${async.msgtype}Output struct {
//#if CHANNEL_MULTI
	EventType string             ` + "`" + `out:"header=${async.dispatch.header}" const:"${async.dispatch.value}"` + "`" + `
//#endif CHANNEL_MULTI
	Payload   api.${async.upmsgtype} ` + "`" + `out:"body"` + "`" + `
}

func (p ${async.msgtype}Publisher) Publish${async.upmsgtype}(ctx context.Context, payload api.${async.upmsgtype}) error {
	return p.messagePublisher.Publish(ctx, ${async.msgtype}Output{
		Payload: payload,
	})
}
`,
	},
	OperationSubscribe: {
		"imports": `
	api "${app.packageurl}/internal/stream/${async.channel.package}/api"
`,
		"dependencies": `
type ${async.upmsgtype}Handler interface {
	On${async.upmsgtype}(ctx context.Context, payload api.${async.upmsgtype}) error
}

`,
		"implementation": `
type ${async.msgtype}Input struct {
//#if CHANNEL_MULTI
	EventType string                       ` + "`" + `in:"header=${async.dispatch.header}" const:"${async.dispatch.value}"` + "`" + `
//#endif CHANNEL_MULTI
	Payload   api.${async.upmsgtype} ` + "`" + `in:"body"` + "`" + `
}

`,
		"handler": `func(ctx context.Context, in *${async.msgtype}Input) error {
	return handler.On${async.upmsgtype}(ctx, in.Payload)
},`,
	},
}

// RenderMessageOperation generates the message operation (publisher/subscriber) implementation
func (g *TemplatedGenerator) RenderMessageOperation(op string, messageId string, snippets map[string]string) (err error) {
	vars := map[string]string{
		"async.channel":         g.cfg.ChannelName,
		"async.channel.package": packageName(channelShortName(g.cfg.ChannelName)),
		"async.operation":       op,
		// message type inflections
		"async.msgtype":         strcase.ToLowerCamel(messageId),
		"async.upmsgtype":       strcase.ToCamel(messageId),
		"async.snakemsgtype":    strcase.ToSnake(messageId),
		"async.dispatch.value":  dispatchName(messageId),
		"async.dispatch.header": "eventType",
		"async.msgtype.human":   humanTypeName(messageId),
		// domain inflections
		"async.domain.package":            packageName(inflector.Pluralize(g.cfg.Domain)),
		"async.domain.uppercamelsingular": strcase.ToCamel(inflector.Singularize(g.cfg.Domain)),
	}

	if snippets == nil {
		snippets = map[string]string{}

		opts := skel.NewRenderOptions()
		opts.AddVariables(vars)
		opts.AddCondition("CHANNEL_MULTI", g.cfg.Multi)

		for k, v := range messageOperationSnippets[op] {
			template := skel.Template{
				Name:       fmt.Sprintf("Creating message operation snippet %s", k),
				SourceData: []byte(v),
				Format:     text.FileFormatGo,
			}

			var contents string
			contents, err = template.RenderContents(opts)
			if err != nil {
				return err
			}
			snippets[k] = contents
		}
	}

	opts := skel.NewRenderOptions()
	opts.AddVariables(vars)
	opts.AddVariables(snippets)
	opts.AddCondition("CHANNEL_MULTI", g.cfg.Multi)

	if err = g.render(opts, op, msgMessage); err != nil {
		return
	}

	if !g.cfg.Deep {
		return
	}

	return g.RenderMessagePayload(op, messageId)
}

// RenderMessagePayload generates the message payload DTO implementation
func (g *TemplatedGenerator) RenderMessagePayload(op, messageId string) (err error) {
	opts := skel.NewRenderOptions()

	opts.AddVariables(map[string]string{
		"async.channel":         g.cfg.ChannelName,
		"async.channel.package": packageName(channelShortName(g.cfg.ChannelName)),
		// message type inflections
		"async.msgtype":          messageId,
		"async.upmsgtype":        strcase.ToCamel(messageId),
		"async.snakemsgtype":     strcase.ToSnake(messageId),
		"async.msgtype.dispatch": dispatchName(messageId),
		"async.msgtype.human":    humanTypeName(messageId),
	})

	opts.AddCondition("CHANNEL_MULTI", g.cfg.Multi)

	return g.render(opts, op, msgPayload)
}

func (g *TemplatedGenerator) render(opts skel.RenderOptions, op string, component msgComponent) (err error) {
	templates := g.templates(op, component)
	if err = templates.Render(opts); err != nil {
		err = errors.Wrap(err, "Failed rendering templates")
		return
	}
	g.dirs.AddAll(templates.Dirs(opts)...)
	return nil
}

// Convenience values for AsyncAPI

type msgComponent int

const (
	msgChannel msgComponent = iota
	msgOperation
	msgMessage
	msgPayload
)

func (g *TemplatedGenerator) templates(operation string, component msgComponent) skel.TemplateSet {
	type templateSetKey struct {
		operation string
		component msgComponent
	}

	key := templateSetKey{
		operation: operation,
		component: component,
	}

	switch key {
	case templateSetKey{operation: OperationPublish, component: msgChannel},
		templateSetKey{operation: OperationSubscribe, component: msgChannel},
		templateSetKey{operation: OperationNone, component: msgChannel}:
		return skel.TemplateSet{
			{Name: "Creating channel",
				Operation:  skel.OpAddNoOverwrite,
				SourceFile: "asyncapi/pkg.go",
				DestFile:   "internal/stream/${async.channel.package}/pkg.go",
				Format:     text.FileFormatGo,
			},
		}
	case templateSetKey{operation: OperationPublish, component: msgPayload},
		templateSetKey{operation: OperationSubscribe, component: msgPayload},
		templateSetKey{operation: OperationNone, component: msgPayload}:
		return skel.TemplateSet{{
			Name:       "Creating message payload",
			Operation:  skel.OpAddNoOverwrite,
			SourceFile: "asyncapi/payload.go",
			DestFile:   "internal/stream/${async.channel.package}/api/${async.snakemsgtype}.go",
			Format:     text.FileFormatGo,
		}}
	case templateSetKey{operation: OperationPublish, component: msgOperation}:
		return skel.TemplateSet{
			{Name: "Creating channel publisher",
				Operation:  skel.OpAddNoOverwrite,
				SourceFile: "asyncapi/publisher_channel.go",
				DestFile:   "internal/stream/${async.channel.package}/publisher_channel.go",
				Format:     text.FileFormatGo,
			},
		}
	case templateSetKey{operation: OperationPublish, component: msgMessage}:
		return skel.TemplateSet{
			{Name: "Creating message publisher",
				Operation:  skel.OpNew,
				SourceFile: "asyncapi/publisher_message.go",
				DestFile:   "internal/stream/${async.channel.package}/publisher_${async.snakemsgtype}.go",
				Format:     text.FileFormatGo,
			},
		}
	case templateSetKey{operation: OperationSubscribe, component: msgOperation}:
		return skel.TemplateSet{
			{Name: "Creating channel subscriber",
				Operation:  skel.OpAddNoOverwrite,
				SourceFile: "asyncapi/subscriber_channel.go",
				DestFile:   "internal/stream/${async.channel.package}/subscriber_channel.go",
				Format:     text.FileFormatGo,
			},
		}
	case templateSetKey{operation: OperationSubscribe, component: msgMessage}:
		return skel.TemplateSet{
			{Name: "Creating message subscriber",
				Operation:  skel.OpNew,
				SourceFile: "asyncapi/subscriber_message.go",
				DestFile:   "internal/stream/${async.channel.package}/subscriber_${async.snakemsgtype}.go",
				Format:     text.FileFormatGo,
			},
		}
	}

	return skel.TemplateSet{}
}

func newTemplatedGenerator(cfg GeneratorConfig) *TemplatedGenerator {
	return &TemplatedGenerator{
		cfg:  cfg,
		dirs: types.StringSet{},
	}
}
