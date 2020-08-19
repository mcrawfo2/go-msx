package build

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func uploadArtifactory(sourceFile string, uploadUrl string) (err error) {
	logger.Infof("Uploading %q to %q", sourceFile, uploadUrl)

	var req = new(http.Request)
	req.URL, err = url.Parse(uploadUrl)
	if err != nil {
		return err
	}

	req.Method = http.MethodPut
	req.Header = make(http.Header)
	req.Header.Set("Authorization", BuildConfig.Binaries.Authorization())

	req.Body, err = os.Open(sourceFile)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to upload binary %q", filepath.Base(sourceFile))
	}

	return nil
}
