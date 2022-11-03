// Copyright © 2022, Cisco Systems Inc. All rights reserved.

package platform

import (
	"archive/zip"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/build/maven"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type ManifestContent struct {
	ManifestInfo ManifestInfo
}

type ManifestInfo struct {
	Versions Versions `yaml:"versions"`
}

type Versions struct {
	PlatformBuildVersion string `yaml:"msx.platform.build.version"`
	PlatformNfvVersion   string `yaml:"msx.platform.nfv.version"`
}

func NewManifestContent(resourcesZipData *maven.HashedData, descriptor maven.ArtifactDescriptor) (result *ManifestContent, err error) {
	byteReader := bytes.NewReader(resourcesZipData.Data)
	zipReader, err := zip.NewReader(byteReader, resourcesZipData.Len())
	manifestInfoFileName := fmt.Sprintf("%s-%s/manifest-info.yaml", descriptor.ArtifactId, descriptor.Version)
	fileData, err := zipReader.Open(manifestInfoFileName)
	if err != nil {
		return
	}
	defer fileData.Close()

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		return
	}
	if err = fileData.Close(); err != nil {
		return
	}

	manifestInfo := ManifestInfo{}
	err = yaml.Unmarshal(fileBytes, &manifestInfo)
	if err != nil {
		return
	}
	return &ManifestContent{ManifestInfo: manifestInfo}, nil
}
