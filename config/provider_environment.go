// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"os"
	"strings"
)

type EnvironmentProvider struct {
	Describer
	SilentNotifier
}

func (p *EnvironmentProvider) Load(ctx context.Context) (ProviderEntries, error) {
	lines := os.Environ()

	var results = make(ProviderEntries, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		entry := NewEntry(p, parts[0], parts[1])
		if entry.NormalizedName == "" || entry.NormalizedName[0] == '.' {
			continue
		}

		results = append(results, entry)
	}

	return results, nil
}

func NewEnvironmentProvider(name string) *EnvironmentProvider {
	return &EnvironmentProvider{
		Describer: Named{
			name: name,
		},
	}
}
