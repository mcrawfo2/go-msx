package maven

import (
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ArtifactDescriptor struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

func (d ArtifactDescriptor) IsSnapshot() bool {
	return strings.HasSuffix(d.Version, "-DEV-SNAPSHOT")
}

func (d ArtifactDescriptor) IsStable() bool {
	return strings.HasSuffix(d.Version, "-STABLE")
}

func (d ArtifactDescriptor) IsEdge() bool {
	return strings.HasSuffix(d.Version, "-UNSTABLE") ||
		strings.HasSuffix(d.Version, "-EDGE")
}

func (d ArtifactDescriptor) Release() string {
	parts := strings.SplitN(d.Version, "-", 2)
	return parts[0]
}

func (d ArtifactDescriptor) WithVersion(version string) ArtifactDescriptor {
	return ArtifactDescriptor{
		GroupId:    d.GroupId,
		ArtifactId: d.ArtifactId,
		Version:    version,
		Scope:      d.Scope,
	}
}

func (d ArtifactDescriptor) PomFileName() string {
	return fmt.Sprintf("%s-%s.pom", d.ArtifactId, d.Version)
}

func (d ArtifactDescriptor) JarFileName() string {
	return fmt.Sprintf("%s-%s.jar", d.ArtifactId, d.Version)
}

func (d ArtifactDescriptor) Triple() string {
	return fmt.Sprintf("%s:%s:%s", d.GroupId, d.ArtifactId, d.Version)
}

func (d ArtifactDescriptor) GroupPath() string {
	return strings.ReplaceAll(d.GroupId, ".", "/")
}

func NewArtifactDescriptor(triple string) ArtifactDescriptor {
	artifactDescriptor := ArtifactDescriptor{}
	parts := strings.SplitN(triple, ":", 4)
	if len(parts) > 0 {
		artifactDescriptor.GroupId = parts[0]
	}
	if len(parts) > 1 {
		artifactDescriptor.ArtifactId = parts[1]
	}
	if len(parts) > 2 {
		artifactDescriptor.Version = parts[2]
	}
	return artifactDescriptor
}

type HashedData struct {
	Data []byte
	Hash []byte
}

func (d HashedData) Len() int64 {
	return int64(len(d.Data))
}

type ReleaseArtifactCollection struct {
	Metadata   *HashedData
	GroupId    string
	ArtifactId string
	Versions   []string
}

func (c *ReleaseArtifactCollection) Artifact(version, fileName string) *ReleaseArtifact {
	return &ReleaseArtifact{
		Collection: c,
		Version:    version,
		FileName:   fileName,
	}
}

func (c *ReleaseArtifactCollection) LatestVersion(baseVersion string) (string, error) {
	// Find all matching versions
	var versionRegex = regexp.MustCompile("([\\d.]+)-(\\d+)")
	var allVersions []int
	for _, v := range c.Versions {
		matches := versionRegex.FindStringSubmatch(v)
		if matches == nil {
			// Unknown build number format - skip this version
			continue
		}
		buildNumber, _ := strconv.ParseInt(matches[2], 10, 32)

		if matches[1] == baseVersion {
			allVersions = append(allVersions, int(buildNumber))
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(allVersions)))

	// Make sure we found something
	if len(allVersions) == 0 {
		return "", errors.Errorf("no builds found")
	}

	return fmt.Sprintf("%s-%d", baseVersion, allVersions[0]), nil
}

func (c *ReleaseArtifactCollection) ForVersion(version string) ArtifactFactory {
	return &ReleaseArtifactCollectionForVersion{
		collection: c,
		version:    version,
	}
}

type ReleaseArtifactCollectionForVersion struct {
	collection *ReleaseArtifactCollection
	version    string
}

func (r *ReleaseArtifactCollectionForVersion) CreateArtifact(fileName string) Artifact {
	return r.collection.Artifact(r.version, fileName)
}

type ReleaseArtifact struct {
	Collection *ReleaseArtifactCollection
	Version    string
	FileName   string
}

func (a *ReleaseArtifact) Retrieve(repository ArtifactRepository) (*HashedData, error) {
	return repository.GetReleaseArtifactData(a)
}

type SnapshotArtifactCollection struct {
	Metadata    *HashedData
	GroupId     string
	ArtifactId  string
	Version     string
	BuildNumber string
	Timestamp   string
}

type ArtifactFactory interface {
	CreateArtifact(fileName string) Artifact
}

type Artifact interface {
	Retrieve(repository ArtifactRepository) (*HashedData, error)
}

func (c *SnapshotArtifactCollection) CreateArtifact(fileName string) Artifact {
	return &SnapshotArtifact{
		Collection: c,
		FileName:   fileName,
	}
}

type SnapshotArtifact struct {
	Collection *SnapshotArtifactCollection
	FileName   string
}

func (a *SnapshotArtifact) Retrieve(repository ArtifactRepository) (*HashedData, error) {
	return repository.GetSnapshotArtifactData(a.Collection, a.FileName)
}

func NewHashedData(data []byte, hash []byte) (*HashedData, error) {
	// TODO: Validate the hash against the data
	return &HashedData{
		Data: data,
		Hash: hash,
	}, nil
}

func NewReleaseArtifactCollection(manifest *HashedData) (*ReleaseArtifactCollection, error) {
	type metadata struct {
		XMLName    xml.Name `xml:"metadata"`
		GroupId    string   `xml:"groupId"`
		ArtifactId string   `xml:"artifactId"`
		Versions   []string `xml:"versioning>versions>version"`
	}

	result := metadata{}
	if err := xml.Unmarshal(manifest.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse build metadata")
	}

	return &ReleaseArtifactCollection{
		Metadata:   manifest,
		ArtifactId: result.ArtifactId,
		GroupId:    strings.Replace(result.GroupId, ".", "/", -1),
		Versions:   result.Versions,
	}, nil
}

func NewSnapshotArtifactCollection(manifest *HashedData) (*SnapshotArtifactCollection, error) {
	type metadata struct {
		XMLName     xml.Name `xml:"metadata"`
		GroupId     string   `xml:"groupId"`
		ArtifactId  string   `xml:"artifactId"`
		Version     string   `xml:"version"`
		BuildNumber string   `xml:"versioning>snapshot>buildNumber"`
		TimeStamp   string   `xml:"versioning>snapshot>timestamp"`
	}

	result := metadata{}
	if err := xml.Unmarshal(manifest.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse build metadata")
	}

	return &SnapshotArtifactCollection{
		Metadata:    manifest,
		ArtifactId:  result.ArtifactId,
		GroupId:     strings.Replace(result.GroupId, ".", "/", -1),
		Version:     result.Version,
		BuildNumber: result.BuildNumber,
		Timestamp:   result.TimeStamp,
	}, nil
}

func NewReleaseArtifact(groupId, artifact, version, fileName string) (*ReleaseArtifact, error) {
	collection := ReleaseArtifactCollection{
		GroupId:    groupId,
		ArtifactId: artifact,
	}

	return collection.Artifact(version, fileName), nil
}

func ResolveArtifactVersion(descriptor ArtifactDescriptor) (result *ArtifactDescriptor, err error) {
	var version string
	var artifactCollection *ReleaseArtifactCollection
	var mavenRepository = NewDefaultHttpRepository()

	if descriptor.IsStable() {
		artifactCollection, err = mavenRepository.GetReleaseArtifactCollection(
			"com/cisco/vms/manifest/EI-Stable",
			"platform-manifest")
		if err == nil {
			version, err = artifactCollection.LatestVersion(descriptor.Release())
		}
	} else if descriptor.IsEdge() {
		artifactCollection, err = mavenRepository.GetReleaseArtifactCollection(
			"com/cisco/vms/manifest/Build-Stable",
			"platform-manifest")
		if err == nil {
			version, err = artifactCollection.LatestVersion(descriptor.Release())
		}
	} else {
		// Just use the specified version for snapshots and specific releases
		version = descriptor.Version
	}

	if err == nil {
		descriptor = descriptor.WithVersion(version)
		result = &descriptor
	}

	return

}

func NewArtifactFactory(descriptor ArtifactDescriptor, repository ArtifactRepository) (ArtifactFactory, error) {
	if descriptor.IsSnapshot() {
		collection, err := repository.GetSnapshotArtifactCollection(
			descriptor.GroupPath(),
			descriptor.ArtifactId,
			descriptor.Version)
		if err != nil {
			return nil, err
		}
		return collection, nil
	} else {
		collection, err := repository.GetReleaseArtifactCollection(
			descriptor.GroupPath(),
			descriptor.ArtifactId)
		if err != nil {
			return nil, err
		}
		return collection.ForVersion(descriptor.Version), nil
	}
}
