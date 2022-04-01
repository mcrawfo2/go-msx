// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package repository

import "github.com/pkg/errors"

var ErrNotFound = errors.New("Entity not found")
var ErrAlreadyExists = errors.New("Entity already exists")
