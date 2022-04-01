// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package filesystemtest

import (
	"fmt"
	"net/http"
	"testing"
)

type FileSystemCheck struct {
	Validators []FileSystemPredicate
}

func (c FileSystemCheck) Check(fs http.FileSystem) []error {
	var results []error

	for _, predicate := range c.Validators {
		if !predicate.Matches(fs) {
			results = append(results, FileSystemCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type FileSystemVerifier func(t *testing.T, fs http.FileSystem)

type FileSystemPredicate struct {
	Description string
	Matches     func(fs http.FileSystem) bool
}

type FileSystemCheckError struct {
	Validator FileSystemPredicate
}

func (c FileSystemCheckError) Error() string {
	return fmt.Sprintf("Failed fs validator: %s", c.Validator.Description)
}
