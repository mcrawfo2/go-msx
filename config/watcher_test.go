// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatcher(t *testing.T) {
	staticProvider := NewInMemoryProvider("static", map[string]string{
		"alpha":   "a",
		"bravo":   "b",
		"charlie": "c",
		"delta":   "d",
		"echo":    "e",
	})

	cache := NewCacheProvider(staticProvider)

	watcher := NewWatcher(cache)

	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		expectedChanges := 2
		version := 1

		for {
			select {
			case <-ctx.Done():
				return

			case n := <-watcher.Notify():
				if n.Error != nil {
					logger.WithError(n.Error).Errorf("Error received from watcher")
				}

				for _, c := range n.Delta {
					if c.IsSet() {
						logger.Infof("[%d] SET: %q = %q", version, c.NewEntry.Name, c.NewEntry.Value)
					} else {
						logger.Infof("[%d] UNSET: %q", version, c.OldEntry.Name)
					}
					expectedChanges--
				}

				if expectedChanges <= 0 {
					cancelCtx()
				}

				version++
			}
		}
	}()

	go func() {
		err := staticProvider.SetValue("foxtrot", "f")
		assert.NoError(t, err)

		err = staticProvider.UnsetValue("charlie")
		assert.NoError(t, err)
	}()

	watcher.Watch(ctx)
}
