// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

// Standard manifest artifact (just points to a file for all of the data)
type Artifact struct {
	TemplateFileName string `json:"templateFileName"`
}

// Standard manifest
type Manifest map[string][]Artifact
