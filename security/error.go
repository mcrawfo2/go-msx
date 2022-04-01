// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package security

import "github.com/pkg/errors"

var (
	ErrTokenNotFound = errors.New("Token missing from context")
)
