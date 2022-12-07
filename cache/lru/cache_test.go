// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func testCache(ttl time.Duration, expireLimit int, age bool,
	metrics bool, prefix string) (*HeapMapCache, *clock.Mock) {
	mockClock := clock.NewMock()
	return NewCache2(ttl, expireLimit, 500*time.Millisecond, age, mockClock, metrics, prefix), mockClock
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
	}

	for _, age := range []bool{false, true} {
		cache, mockClock := testCache(1*time.Second, 1, age, false, "")
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cache.Clear()
				assert.Empty(t, cache.index)

				for _, preset := range tt.preset {
					cache.Set(preset.key, preset.value)
				}
				mockClock.Add(500 * time.Millisecond)
				cache.Set(tt.setArgs.key, tt.setArgs.value)
				mockClock.Add(tt.advance)
				got, gotOk := cache.Get(tt.getKey)
				if gotOk != tt.wantOk || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = (%v,%v) want (%v,%v)", got, gotOk, tt.want, tt.wantOk)
				}

				// Allow everything to expire
				mockClock.Add(2 * time.Second)
				assert.Empty(t, cache.index)
			})
		}
	}
}
