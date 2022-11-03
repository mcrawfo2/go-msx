// Copyright © 2022, Cisco Systems Inc. All rights reserved.

package platform

import (
	"cto-github.cisco.com/NFV-BU/go-msx/build/maven"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var (
	logger                  = log.NewLogger("build.maven.platform")
	ArtifactVersionNotFound = errors.New("artifact version not found")
	platformManifestStore   = map[string]*Manifest{}
)

const (
	NFVGroupIDPrefix = "com.cisco.nfv"
	VMSGroupIDPrefix = "com.cisco.vms"
)

type Manifest struct {
	GroupId  string
	Version  string
	Snapshot bool
	Content  *ManifestContent
}

func (m *Manifest) readPlatformManifest(artifact maven.ArtifactDescriptor, repository maven.ArtifactRepository) (result *ManifestContent, err error) {
	platformArtifactFactory, err := maven.NewArtifactFactory(artifact, repository)
	if err != nil {
		return
	}

	resourcesZipArtifact := platformArtifactFactory.CreateArtifact(artifact.ResourcesZipName())
	resoucesZipData, err := resourcesZipArtifact.Retrieve(repository)
	if err != nil {
		return
	}
	return NewManifestContent(resoucesZipData, artifact)
}

// ResolveDependencyVersion set dependency version by groupId
func (m *Manifest) ResolveDependencyVersion(artifact maven.ArtifactDescriptor, sourceArtifact maven.ArtifactDescriptor) (version string) {
	if len(artifact.Version) > 0 {
		return artifact.Version
	}
	// defaults to sourceArtifact version
	version = sourceArtifact.Version
	if IsNFVLib(artifact) {
		version = m.NfvLibVersion()
	} else if IsVMSLib(artifact) {
		version = m.FrameworkVersion()
	}
	return
}

func (m *Manifest) FrameworkVersion() string {
	if m.Snapshot {
		return m.Version
	}
	return m.Content.ManifestInfo.Versions.PlatformBuildVersion
}

func (m *Manifest) NfvLibVersion() string {
	if m.Snapshot {
		return m.Version
	}
	v := m.Content.ManifestInfo.Versions.PlatformNfvVersion
	if len(v) > 0 {
		return v
	} else {
		return m.FrameworkVersion()
	}
}

// download manifest for maven artifactory
func (m *Manifest) downloadManifest(mavenRepository *maven.HttpRepository) (err error) {

	logger.Infof("platform manifest version: %s", m.Version)

	platformArtifact := maven.ArtifactDescriptor{
		GroupId:    m.GroupId,
		ArtifactId: "platform-manifest",
		Version:    m.Version,
		Scope:      "",
	}

	manifestContent, err := m.readPlatformManifest(platformArtifact, mavenRepository)

	if err != nil {
		return
	}

	frameworkVersion := manifestContent.ManifestInfo.Versions.PlatformBuildVersion
	if len(frameworkVersion) == 0 {
		err = errors.Wrapf(ArtifactVersionNotFound, "versions.msx.platform.build.version not found in %s:%s:%s", platformArtifact.GroupId, platformArtifact.ArtifactId, platformArtifact.Version)
		return
	}
	m.Content = manifestContent
	logger.Infof("platform framework version: %s", frameworkVersion)
	return
}

func (m *Manifest) ResolveArtifactVersion(descriptor maven.ArtifactDescriptor) (result *maven.ArtifactDescriptor) {
	var version string
	if IsNFVLib(descriptor) {
		version = m.NfvLibVersion()
	} else {
		version = m.FrameworkVersion()
	}
	descriptor = descriptor.WithVersion(version)
	return &descriptor
}

func GetPlatformManifest(descriptor maven.ArtifactDescriptor) (result *Manifest, err error) {

	var platformManifest *Manifest
	if m, ok := platformManifestStore[descriptor.Version]; ok {
		platformManifest = m
	} else {
		platformManifest, err = NewManifest(descriptor)
	}
	if err != nil {
		return
	}

	platformManifestStore[descriptor.Version] = platformManifest
	platformManifestStore[platformManifest.Version] = platformManifest

	return platformManifest, nil
}

func HasVersionNumber(v string) bool {
	vs := strings.Split(v, "-")
	if len(vs) == 2 {
		if n, e := strconv.Atoi(vs[1]); e == nil && n > 0 {
			return true
		}
	}
	return false
}

func IsNFVLib(d maven.ArtifactDescriptor) bool {
	return strings.HasPrefix(d.GroupId, NFVGroupIDPrefix)
}

func IsVMSLib(d maven.ArtifactDescriptor) bool {
	return strings.HasPrefix(d.GroupId, VMSGroupIDPrefix)
}

func NewManifest(descriptor maven.ArtifactDescriptor) (result *Manifest, err error) {
	var platformManifestGroupId string

	if descriptor.IsStable() {
		platformManifestGroupId = "com/cisco/vms/manifest/EI-Stable"
	} else if descriptor.IsEdge() || HasVersionNumber(descriptor.Version) {
		platformManifestGroupId = "com/cisco/vms/manifest/Build-Stable"
	}

	platformManifest := Manifest{
		GroupId:  platformManifestGroupId,
		Version:  descriptor.Version,
		Snapshot: false,
	}

	if len(platformManifestGroupId) == 0 {
		platformManifest.Snapshot = true
		return &platformManifest, nil
	}
	var mavenRepository = maven.NewDefaultHttpRepository()
	if !HasVersionNumber(descriptor.Version) {
		var artifactCollection *maven.ReleaseArtifactCollection

		artifactCollection, err = mavenRepository.GetReleaseArtifactCollection(
			platformManifestGroupId,
			"platform-manifest")
		if err != nil {
			return
		}
		var version string
		version, err = artifactCollection.LatestVersion(descriptor.Release())
		if err != nil {
			return
		}
		platformManifest.Version = version
	}

	err = platformManifest.downloadManifest(mavenRepository)
	return &platformManifest, err
}
