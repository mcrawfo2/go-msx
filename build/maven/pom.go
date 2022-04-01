// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package maven

import (
	"encoding/xml"
	"fmt"
)

type PomFile struct {
	Packaging    string
	Parent       ArtifactDescriptor
	Dependencies []ArtifactDescriptor
}

func NewPomFile(data *HashedData) (*PomFile, error) {
	type pom struct {
		XMLName      xml.Name             `xml:"project"`
		Packaging    string               `xml:"packaging"`
		Parent       ArtifactDescriptor   `xml:"parent"`
		Dependencies []ArtifactDescriptor `xml:"dependencies>dependency"`
	}

	result := pom{}
	if err := xml.Unmarshal(data.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse pom file")
	}

	return &PomFile{
		Packaging:    result.Packaging,
		Parent:       result.Parent,
		Dependencies: result.Dependencies,
	}, nil
}
