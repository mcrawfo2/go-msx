// Copyright © 2022, Cisco Systems Inc. All rights reserved.

package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

const ciscoArtifactoryUrl = "https://engci-maven-master.cisco.com/artifactory/api/npm/vms-npm-local"

type Package struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Dist        struct {
		Tarball string `json:"tarball"`
		Shasum  string `json:"shasum"`
	} `json:"dist"`
	Files         []string          `json:"files"`
	Scripts       map[string]string `json:"scripts"`
	PublishConfig struct {
		Access string `json:"access"`
	} `json:"publishConfig"`
}

type HttpRepository struct {
	BaseUrl string
}

func NewHttpRepository() *HttpRepository {
	artifactoryUrl := os.Getenv("ARTIFACTORY_URL")
	if artifactoryUrl == "" {
		artifactoryUrl = ciscoArtifactoryUrl
	}

	return &HttpRepository{
		BaseUrl: artifactoryUrl,
	}

}

func (r *HttpRepository) PackageInfo(packageName string, version string) (*Package, error) {
	sourceUrl := r.BaseUrl + "/" + path.Join(packageName, version)
	packageInfoBytes, err := r.Get(sourceUrl)
	if err != nil {
		return nil, err
	}

	var p Package
	if err := json.Unmarshal(packageInfoBytes, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *HttpRepository) Get(sourceUrl string) ([]byte, error) {
	sourceParsedUrl, _ := url.Parse(sourceUrl)
	baseParsedUrl, _ := url.Parse(r.BaseUrl)

	// Always resolve from our preferred artifactory location
	sourceParsedUrl.Scheme = baseParsedUrl.Scheme
	sourceParsedUrl.User = baseParsedUrl.User
	sourceParsedUrl.Host = baseParsedUrl.Host
	sourceUrl = sourceParsedUrl.String()

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
