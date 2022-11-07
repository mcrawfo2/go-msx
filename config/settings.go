// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"crypto/sha1"
	"sort"
)

var nullSeparator = []byte{0}

func SettingsHash(settings map[string]string) []byte {
	var keys []string
	for k := range settings {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	h := sha1.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write(nullSeparator)
		h.Write([]byte(settings[k]))
		h.Write(nullSeparator)
	}

	return h.Sum(nil)
}
