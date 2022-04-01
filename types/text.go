// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

type TextUnmarshaler interface {
	UnmarshalText(data string) error
}

type TextMarshaler interface {
	MarshalText() (string, error)
}
