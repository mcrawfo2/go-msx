// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/thejerf/abtime"
	"reflect"
	"strings"
	"testing"
	"time"
)

func advanceCache(cache *HeapMapCache, clock *abtime.ManualTime, period time.Duration) {
	clock.Advance(period)
	ticks := period
	for ticks >= 0 {
		clock.Trigger(sleeperId)
		<-cache.expired
		ticks -= cache.expireFrequency
	}
}

func TestHeapMapCache(t *testing.T) {
	type kvpair struct {
		key   string
		value interface{}
	}

	tests := []struct {
		name    string
		preset  []kvpair
		setArgs kvpair
		getKey  string
		advance time.Duration
		want    interface{}
		wantOk  bool
	}{
		{
			name: "Clean",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key1",
			want:   "value1",
			wantOk: true,
		},
		{
			name: "Overwrite",
			preset: []kvpair{
				{
					key:   "key1",
					value: "value2",
				},
				{
					key:   "key2",
					value: "value3",
				},
			},
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key1",
			want:   "value1",
			wantOk: true,
		},
		{
			name: "Expand",
			preset: []kvpair{{
				key:   "key1",
				value: "value1",
			}},
			setArgs: kvpair{
				key:   "key2",
				value: "value2",
			},
			getKey: "key2",
			want:   "value2",
			wantOk: true,
		},
		{
			name: "Missing",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key2",
			wantOk: false,
		},
		{
			name: "Expired",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey:  "key1",
			advance: time.Second + time.Millisecond,
			wantOk:  false,
		},
		{
			name: "DeageClean",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key1",
			want:   "value1",
			wantOk: true,
		},
		{
			name: "DeageOverwrite",
			preset: []kvpair{
				{
					key:   "key1",
					value: "value2",
				},
				{
					key:   "key2",
					value: "value3",
				},
			},
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key1",
			want:   "value1",
			wantOk: true,
		},
		{
			name: "DeageExpand",
			preset: []kvpair{{
				key:   "key1",
				value: "value1",
			}},
			setArgs: kvpair{
				key:   "key2",
				value: "value2",
			},
			getKey: "key2",
			want:   "value2",
			wantOk: true,
		},
		{
			name: "DeageMissing",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey: "key2",
			wantOk: false,
		},
		{
			name: "DeageExpired",
			setArgs: kvpair{
				key:   "key1",
				value: "value1",
			},
			getKey:  "key1",
			advance: time.Second + time.Millisecond,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClock := types.NewMockClock()
			ttl := 1 * time.Second
			expireLimit := 1
			expireFrequency := 500 * time.Millisecond
			deage := strings.HasPrefix("Deage", tt.name)

			cache := NewCache2(ttl, expireLimit, expireFrequency, deage, mockClock, false, "")
			cache.expired = make(chan struct{})

			// Load the cache initial state with
			for _, preset := range tt.preset {
				cache.Set(preset.key, preset.value)
			}

			// Expire once
			advanceCache(cache, mockClock, 500*time.Millisecond)

			// Apply our test
			cache.Set(tt.setArgs.key, tt.setArgs.value)

			// Apply our wait
			advanceCache(cache, mockClock, tt.advance)

			// Check our result
			got, gotOk := cache.Get(tt.getKey)
			if gotOk != tt.wantOk || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = (%v,%v) want (%v,%v)", got, gotOk, tt.want, tt.wantOk)
			}

			// Allow everything to expire
			advanceCache(cache, mockClock, 2*time.Second)
			assert.Empty(t, cache.index)

			cache.shutdown = true
		})
	}
}
