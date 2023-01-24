// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

type Transformer func(string) string

type Transformers []Transformer

func (t Transformers) Transform(target string) string {
	for _, transformer := range t {
		target = transformer(target)
	}

	return target
}
