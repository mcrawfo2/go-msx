package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/filesystemtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewWebRoot(t *testing.T) {
	err := fs.SetSources()
	assert.NoError(t, err)

	webroot, err := NewWebRoot("/webservice/testdata/webroot")
	assert.NoError(t, err)
	assert.NotNil(t, webroot)

	errs := filesystemtest.FileSystemCheck{
		Validators: []filesystemtest.FileSystemPredicate{
			{
				Description: "NoRootIndex",
				Matches: func(fs http.FileSystem) bool {
					_, err := fs.Open("/")
					return err != nil
				},
			},
			{
				Description: "FileIndex",
				Matches: func(fs http.FileSystem) bool {
					_, err := fs.Open("/secondary/index.html")
					return err == nil
				},
			},
		},
	}.Check(webroot)

	if len(errs) > 0 {
		testhelpers.ReportErrors(t, "FS", errs)
	}
}
