package maven

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

const (
	metadataFileName    = "maven-metadata.xml"
	ciscoArtifactoryUrl = "http://engci-maven-master.cisco.com/artifactory/symphony-group"
)

type ArtifactRepository interface {
	GetReleaseArtifactCollection(groupId, artifactId string) (*ReleaseArtifactCollection, error)
	GetSnapshotArtifactCollection(groupId, artifactId, version string) (*SnapshotArtifactCollection, error)

	GetReleaseArtifactData(artifact *ReleaseArtifact) (*HashedData, error)
	GetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string) (*HashedData, error)
}

type MutableArtifactRepository interface {
	SetReleaseArtifactCollection(collection *ReleaseArtifactCollection) error
	SetSnapshotArtifactCollection(collection *SnapshotArtifactCollection) error
	SetReleaseArtifactData(artifact *ReleaseArtifact, hashedData *HashedData) error
	SetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string, hashedData *HashedData) error
}

type FileRepository struct {
	BasePath string
}

func (r *FileRepository) getArtifactCollectionPath(groupId, artifactId string) string {
	return path.Join(r.BasePath, groupId, artifactId)
}

func (r *FileRepository) getArtifactPath(groupId, artifactId, version, fileName string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return path.Join(collectionPath, version, fileName)
}

func (r *FileRepository) getReleaseFilePath(groupId, artifactId, subPath string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return path.Join(collectionPath, subPath)
}

func (r *FileRepository) getSnapshotFilePath(groupId, artifactId, version, subPath string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return path.Join(collectionPath, version, subPath)
}

func (r *FileRepository) getHashedData(filePath string) (*HashedData, error) {
	// Load the file
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Load the file's sha256
	hashPath := filePath + ".sha256"
	hashBytes, err := ioutil.ReadFile(hashPath)
	if err != nil {
		return nil, err
	}

	return NewHashedData(fileBytes, hashBytes)
}

func (r *FileRepository) GetReleaseArtifactCollection(groupId, artifactId string) (*ReleaseArtifactCollection, error) {
	filePath := r.getReleaseFilePath(groupId, artifactId, metadataFileName)

	if manifestData, err := r.getHashedData(filePath); err != nil {
		return nil, err
	} else {
		return NewReleaseArtifactCollection(manifestData)
	}
}

func (r *FileRepository) GetSnapshotArtifactCollection(groupId, artifactId, version string) (*SnapshotArtifactCollection, error) {
	filePath := r.getSnapshotFilePath(groupId, artifactId, version, metadataFileName)

	if manifestData, err := r.getHashedData(filePath); err != nil {
		return nil, err
	} else {
		return NewSnapshotArtifactCollection(manifestData)
	}
}

func (r *FileRepository) GetReleaseArtifactData(artifact *ReleaseArtifact) (*HashedData, error) {
	filePath := r.getArtifactPath(artifact.Collection.GroupId, artifact.Collection.ArtifactId, artifact.Version, artifact.FileName)
	return r.getHashedData(filePath)
}

func (r *FileRepository) GetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string) (*HashedData, error) {
	filePath := r.getArtifactPath(collection.GroupId, collection.ArtifactId, collection.Version, fileName)
	return r.getHashedData(filePath)
}

func (r *FileRepository) SetReleaseArtifactCollection(collection *ReleaseArtifactCollection) error {
	filePath := r.getReleaseFilePath(collection.GroupId, collection.ArtifactId, metadataFileName)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		_ = os.MkdirAll(fileDir, 0755)
	}
	if err := ioutil.WriteFile(filePath, collection.Metadata.Data, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath+".sha256", collection.Metadata.Hash, 0644)
}

func (r *FileRepository) SetSnapshotArtifactCollection(collection *SnapshotArtifactCollection) error {
	filePath := r.getSnapshotFilePath(collection.GroupId, collection.ArtifactId, collection.Version, metadataFileName)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		_ = os.MkdirAll(fileDir, 0755)
	}
	if err := ioutil.WriteFile(filePath, collection.Metadata.Data, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath+".sha256", collection.Metadata.Hash, 0644)
}

func (r *FileRepository) SetReleaseArtifactData(artifact *ReleaseArtifact, hashedData *HashedData) error {
	filePath := r.getArtifactPath(artifact.Collection.GroupId, artifact.Collection.ArtifactId, artifact.Version, artifact.FileName)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		_ = os.MkdirAll(fileDir, 0755)
	}
	if err := ioutil.WriteFile(filePath, hashedData.Data, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath+".sha256", hashedData.Hash, 0644)
}

func (r *FileRepository) SetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string, hashedData *HashedData) error {
	filePath := r.getArtifactPath(collection.GroupId, collection.ArtifactId, collection.Version, fileName)
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		_ = os.MkdirAll(fileDir, 0755)
	}
	if err := ioutil.WriteFile(filePath, hashedData.Data, 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath+".sha256", hashedData.Hash, 0644)
}

func NewDefaultFileRepository() *FileRepository {
	return &FileRepository{
		BasePath: path.Join(HomeDir(), ".go-msx", "build"),
	}
}

type HttpRepository struct {
	BaseUrl string
}

func (r *HttpRepository) getArtifactCollectionPath(groupId, artifactId string) string {
	return fmt.Sprintf("%s/%s/%s", r.BaseUrl, groupId, artifactId)
}

func (r *HttpRepository) getArtifactPath(groupId, artifactId, version, fileName string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return collectionPath + "/" + version + "/" + fileName
}

func (r *HttpRepository) getReleaseFileUrl(groupId, artifactId, subPath string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return collectionPath + "/" + subPath
}

func (r *HttpRepository) getSnapshotFileUrl(groupId, artifactId, version, subPath string) string {
	collectionPath := r.getArtifactCollectionPath(groupId, artifactId)
	return collectionPath + "/" + version + "/" + subPath
}

func (r *HttpRepository) readUrl(sourceUrl string) ([]byte, error) {
	response, err := http.Get(sourceUrl)
	if err != nil {
		var message string
		if urlErr, ok := err.(*url.Error); ok {
			message = urlErr.Err.Error()
		} else {
			message = "unknown error"
		}
		return nil, fmt.Errorf("unable to connect to Artifactory: %s", message)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received response status code %d from Artifactory", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Artifactory")
	}

	return body, nil
}

func (r *HttpRepository) getHashedData(filePath string) (*HashedData, error) {
	// Load the file
	fileBytes, err := r.readUrl(filePath)
	if err != nil {
		return nil, err
	}

	// Load the file's sha256
	hashPath := filePath + ".sha256"
	hashBytes, err := r.readUrl(hashPath)
	if err != nil {
		return nil, err
	}

	return NewHashedData(fileBytes, hashBytes)
}

func (r *HttpRepository) GetReleaseArtifactCollection(groupId, artifactId string) (*ReleaseArtifactCollection, error) {
	fileUrl := r.getReleaseFileUrl(groupId, artifactId, metadataFileName)

	if manifestData, err := r.getHashedData(fileUrl); err != nil {
		return nil, err
	} else {
		return NewReleaseArtifactCollection(manifestData)
	}
}

func (r *HttpRepository) GetSnapshotArtifactCollection(groupId, artifactId, version string) (*SnapshotArtifactCollection, error) {
	filePath := r.getSnapshotFileUrl(groupId, artifactId, version, metadataFileName)

	if manifestData, err := r.getHashedData(filePath); err != nil {
		return nil, err
	} else {
		return NewSnapshotArtifactCollection(manifestData)
	}
}

func (r *HttpRepository) GetReleaseArtifactData(artifact *ReleaseArtifact) (*HashedData, error) {
	filePath := r.getArtifactPath(artifact.Collection.GroupId, artifact.Collection.ArtifactId, artifact.Version, artifact.FileName)
	return r.getHashedData(filePath)
}

func (r *HttpRepository) GetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string) (*HashedData, error) {
	filePath := r.getArtifactPath(collection.GroupId, collection.ArtifactId, collection.Version, fileName)
	return r.getHashedData(filePath)
}

func NewDefaultHttpRepository() *HttpRepository {
	artifactoryUrl := os.Getenv("ARTIFACTORY_URL")
	if artifactoryUrl == "" {
		artifactoryUrl = ciscoArtifactoryUrl
	}

	return &HttpRepository{
		BaseUrl: artifactoryUrl,
	}
}

type CachingRepository struct {
	Source ArtifactRepository
	Cache  MutableArtifactRepository
}

func (r *CachingRepository) GetReleaseArtifactCollection(groupId, artifactId string) (*ReleaseArtifactCollection, error) {
	artifactCollection, err := r.Source.GetReleaseArtifactCollection(groupId, artifactId)
	if err != nil {
		return nil, err
	}

	err = r.Cache.SetReleaseArtifactCollection(artifactCollection)
	if err != nil {
		return nil, err
	}

	return artifactCollection, err
}

func (r *CachingRepository) GetSnapshotArtifactCollection(groupId, artifactId, version string) (*SnapshotArtifactCollection, error) {
	artifactCollection, err := r.Source.GetSnapshotArtifactCollection(groupId, artifactId, version)
	if err != nil {
		return nil, err
	}

	err = r.Cache.SetSnapshotArtifactCollection(artifactCollection)
	if err != nil {
		return nil, err
	}

	return artifactCollection, err
}

func (r *CachingRepository) GetReleaseArtifactData(artifact *ReleaseArtifact) (*HashedData, error) {
	hashedData, err := r.Source.GetReleaseArtifactData(artifact)
	if err != nil {
		return nil, err
	}

	err = r.Cache.SetReleaseArtifactData(artifact, hashedData)
	if err != nil {
		return nil, err
	}

	return hashedData, err
}

func (r *CachingRepository) GetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string) (*HashedData, error) {
	artifact, err := r.Source.GetSnapshotArtifactData(collection, fileName)
	if err != nil {
		return nil, err
	}

	err = r.Cache.SetSnapshotArtifactData(collection, fileName, artifact)
	if err != nil {
		return nil, err
	}

	return artifact, err
}

type FailoverRepository struct {
	PrimaryRepository   ArtifactRepository
	SecondaryRepository ArtifactRepository
}

func (r *FailoverRepository) GetReleaseArtifactCollection(groupId, artifactId string) (*ReleaseArtifactCollection, error) {
	artifactCollection, err := r.PrimaryRepository.GetReleaseArtifactCollection(groupId, artifactId)
	if err == nil {
		return artifactCollection, err
	}

	return r.SecondaryRepository.GetReleaseArtifactCollection(groupId, artifactId)
}

func (r *FailoverRepository) GetSnapshotArtifactCollection(groupId, artifactId, version string) (*SnapshotArtifactCollection, error) {
	artifactCollection, err := r.PrimaryRepository.GetSnapshotArtifactCollection(groupId, artifactId, version)
	if err == nil {
		return artifactCollection, err
	}

	return r.SecondaryRepository.GetSnapshotArtifactCollection(groupId, artifactId, version)
}

func (r *FailoverRepository) GetReleaseArtifactData(artifact *ReleaseArtifact) (*HashedData, error) {
	if hashedData, err := r.PrimaryRepository.GetReleaseArtifactData(artifact); err == nil {
		return hashedData, err
	} else {
		return r.SecondaryRepository.GetReleaseArtifactData(artifact)
	}
}

func (r *FailoverRepository) GetSnapshotArtifactData(collection *SnapshotArtifactCollection, fileName string) (*HashedData, error) {
	if artifact, err := r.PrimaryRepository.GetSnapshotArtifactData(collection, fileName); err == nil {
		return artifact, err
	} else {
		return r.SecondaryRepository.GetSnapshotArtifactData(collection, fileName)
	}
}

func NewDefaultFailoverRepository() ArtifactRepository {
	fileRepository := NewDefaultFileRepository()
	return &FailoverRepository{
		PrimaryRepository: &CachingRepository{
			Source: NewDefaultHttpRepository(),
			Cache:  fileRepository,
		},
		SecondaryRepository: fileRepository,
	}
}

func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
