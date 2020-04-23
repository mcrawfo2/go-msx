package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDuration_MarshalJSON(t *testing.T) {
	duration := Duration(5 * time.Minute)
	durationJson, err := duration.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, durationJson)

	durationString := string(durationJson)
	assert.Equal(t, "\"5m0s\"", durationString)
}

func TestDuration_UnmarshalJSON_String(t *testing.T) {
	durationString := "\"5m\""
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.Nil(t, err)
	assert.Equal(t, float64(300), duration.Duration().Seconds())
}

func TestDuration_UnmarshalJSON_Float64(t *testing.T) {
	durationString := "300000000000"
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.Nil(t, err)
	assert.Equal(t, float64(300), duration.Duration().Seconds())
}

func TestDuration_UnmarshalJSON_Null(t *testing.T) {
	durationString := "null"
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.Nil(t, err)
	assert.Equal(t, float64(0), duration.Duration().Seconds())
}

func TestDuration_UnmarshalJSON_InvalidString(t *testing.T) {
	durationString := "\"44z\""
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.NotNil(t, err)
}

func TestDuration_UnmarshalJSON_InvalidType(t *testing.T) {
	durationString := "{}"
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.NotNil(t, err)
}

func TestDuration_UnmarshalJSON_InvalidJson(t *testing.T) {
	durationString := "zzzz"
	durationBytes := []byte(durationString)
	var duration = new(Duration)
	err := duration.UnmarshalJSON(durationBytes)
	assert.NotNil(t, err)
}
