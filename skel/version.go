// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

const Repository = "https://engci-maven-master.cisco.com/artifactory/api/storage/symphony-release/com/cisco/vms/go-msx-skel"

func init() {
	AddTarget("version", "Show current and available versions", ShowVersion)
}

type entry struct {
	Uri    string `json:"uri"`
	Folder bool   `json:"folder"`
}

type folder struct {
	Repo         string    `json:"repo"`
	Path         string    `json:"path"`
	Created      time.Time `json:"created"`
	CreatedBy    string    `json:"createdBy"`
	LastModified time.Time `json:"lastModified"`
	ModifiedBy   string    `json:"modifiedBy"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Children     []entry   `json:"children"`
	Uri          string    `json:"uri"`
}

var versionFolderMatcher = regexp.MustCompile(`^/\d+$`)

func ShowVersion(_ []string) error {
	logger.Infof("Current build: %d", buildNumber)

	client := &http.Client{Transport: &http.Transport{}}
	req, _ := http.NewRequest("GET", Repository, http.NoBody)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var f folder
	if err = json.Unmarshal(respBodyBytes, &f); err != nil {
		return err
	}

	latestBuild := 0
	var v int64
	for _, c := range f.Children {
		if !c.Folder {
			continue
		}
		if !versionFolderMatcher.MatchString(c.Uri) {
			continue
		}
		if v, err = strconv.ParseInt(c.Uri[1:], 10, 64); err != nil {
			continue
		}
		if int(v) > latestBuild {
			latestBuild = int(v)
		}
	}

	if latestBuild > buildNumber {
		logger.Infof("Newer build available: %d", latestBuild)
		logger.Infof("  - %s/%d/go-msx-skel-%s-%d.tar.gz", f.Uri, latestBuild, runtime.GOOS, latestBuild)
	} else {
		logger.Infof("Latest build: %d", latestBuild)
	}

	return nil
}

