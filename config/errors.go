// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import "github.com/pkg/errors"

var ErrEmptyKey = errors.New("Empty key not allowed")
var ErrDuplicateKey = errors.New("Duplicate key not allowed")
var ErrNotLoaded = errors.New("Configuration not loaded")
var ErrNotFound = errors.New("Setting key not found")
var ErrEmptyValue = errors.New("Empty value found for required variable value")
var ErrCircularReference = errors.New("Circular reference detected")

var ErrParseInvalidCloseBrace = errors.New("Invalid close brace")
var ErrParseInvalidVariableReference = errors.New("Invalid reference in variable name")
var ErrParseUnexpectedInput = errors.New("Unexpected input")

var ErrInvalidValue = errors.New("Failed to parse value")
var ErrValueCannotBeSet = errors.New("Cannot set value of targets")
var ErrNoPopulator = errors.New("No config populator for type")
