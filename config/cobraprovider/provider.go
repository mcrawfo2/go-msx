// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cobraprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/args"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var logger = log.NewLogger("msx.cli.cobra")

const (
	configKeyAppName = "spring.application.name"
)

type Provider struct {
	config.Named
	prefix  string
	appName string
	flagSet *pflag.FlagSet
	extras  *args.CommandLineExtrasProvider
	config.SilentNotifier
}

func (p *Provider) flagExists(name string) bool {
	return p.flagSet.Lookup(name) != nil
}

func (p *Provider) Load(ctx context.Context) (entries config.ProviderEntries, err error) {
	entries, err = p.extras.Load(ctx)
	if err != nil {
		return nil, err
	}
	entries = entries.Clone()

	if !p.flagSet.Parsed() {
		if err := p.flagSet.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	p.flagSet.VisitAll(func(flag *pflag.Flag) {
		key := config.PrefixWithName(p.prefix, flag.Name)
		entries = append(entries, config.NewEntry(p, key, flag.Value.String()))
	})

	// Apply application name
	entries = append(entries, config.NewEntry(p, configKeyAppName, p.appName))

	return entries, nil
}

func extractFlagSet(command *cobra.Command) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(command.Name(), pflag.ContinueOnError)
	flagSet.ParseErrorsWhitelist.UnknownFlags = true

	command.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		flagSet.AddFlag(flag)
	})

	command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		flagSet.AddFlag(flag)
	})

	return flagSet
}

func NewProvider(name string, command *cobra.Command, prefix string) *Provider {
	flagSet := extractFlagSet(command)
	p := &Provider{
		Named:   config.NewNamed(fmt.Sprintf("%s: cobra", name)),
		prefix:  prefix,
		flagSet: flagSet,
		appName: command.Root().Name(),
	}
	p.extras = args.NewCommandLineExtrasProvider(p, p.flagExists)
	return p
}
