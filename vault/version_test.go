package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_Decode_VersionedMetadata(t *testing.T) {
	metadataBytes, err := ioutil.ReadFile("testdata/version-metadata.json")
	assert.NoError(t, err)

	var metadataPojo types.Pojo
	err = json.Unmarshal(metadataBytes, &metadataPojo)
	assert.NoError(t, err)

	var results VersionedMetadata
	err = mapstructure.Decode(metadataPojo, &results)
	assert.NoError(t, err)

	assert.Equal(t, 4, results.CurrentVersion)
	assert.Equal(t, 0, results.OldestVersion)
	assert.Len(t, results.Versions, 4)
}
