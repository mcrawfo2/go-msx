package fs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSourceFileSystem(t *testing.T) {
	fs, err := newSourceFileSystem()
	assert.NoError(t, err)
	assert.NotNil(t, fs)

	_, err = fs.Open("go.mod")
	assert.NoError(t, err)
}

func TestNewDistFileSystem(t *testing.T) {
	fs, err := newStagingFileSystem()
	assert.NoError(t, err)
	assert.NotNil(t, fs)

	_, err = fs.Open("/etc/someservice/buildinfo.yml")
	assert.NoError(t, err)
}
