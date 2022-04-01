// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package args

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"os"
	"strings"
	"sync"
)

type FlagExistsFunc func(string) bool

type CommandLineExtrasProvider struct {
	provider config.Provider
	exists   FlagExistsFunc
	entries  config.ProviderEntries
	once     sync.Once
}

func (c *CommandLineExtrasProvider) Load(_ context.Context) (config.ProviderEntries, error) {
	c.once.Do(func() {
		strip := c.parseArgs()
		c.stripArgs(strip)
	})

	return c.entries, nil
}

func (c *CommandLineExtrasProvider) parseArgs() []int {
	var result config.ProviderEntries
	args := os.Args[1:]
	var strip []int
	for n := 0; n < len(args); n++ {
		v := args[n]
		if len(v) < 2 {
			continue
		}
		if v == "--" {
			break
		}
		if !strings.HasPrefix(v, "--") {
			continue
		}
		v = v[2:]
		split := strings.SplitN(v, "=", 2)
		if len(split) == 2 {
			result = append(result, config.NewEntry(c.provider, split[0], split[1]))
			strip = append(strip, n)
		} else if n == len(args)-1 {
			continue
		} else if strings.HasPrefix(args[n+1], "--") {
			continue
		} else if c.exists(v) {
			// Flag exists
			n++
		} else {
			result = append(result, config.NewEntry(c.provider, v, args[n+1]))
			strip = append(strip, n, n+1)
			n++
		}
	}

	c.entries = result

	return strip
}

func (c *CommandLineExtrasProvider) stripArgs(strip []int) {
	// Strip arguments we processed
	args := os.Args[1:]
	for n := len(strip) - 1; n >= 0; n-- {
		idx := strip[n]
		args = append(args[:idx], args[idx+1:]...)
	}
	os.Args = append(os.Args[:1], args...)
}

func NewCommandLineExtrasProvider(p config.Provider, exists FlagExistsFunc) *CommandLineExtrasProvider {
	return &CommandLineExtrasProvider{
		provider: p,
		exists:   exists,
	}
}
