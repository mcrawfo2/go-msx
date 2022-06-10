package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 1. dest is empty, source is not
func TestRecursiveMerge_DestBaseSrcEmpty(t *testing.T) {
	src := make(map[string]interface{})
	src["key1"] = "key1val"
	src["key2"] = "key2val"

	dest := make(map[string]interface{})
	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = "key1val"
	expectedResult["key2"] = "key2val"
	assert.Equal(t, expectedResult, result)
}

// 2. source is empty, dest is not
func TestRecursiveMerge_DestEmptySrcBase(t *testing.T) {
	dest := make(map[string]interface{})
	dest["key1"] = "key1val"
	dest["key2"] = "key2val"

	src := make(map[string]interface{})
	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = "key1val"
	expectedResult["key2"] = "key2val"
	assert.Equal(t, expectedResult, result)
}

// 3. src, dest are all basic values
func TestRecursiveMerge_DestBasicSrcBasic(t *testing.T) {
	dest := make(map[string]interface{})
	dest["key1"] = "key1val"
	dest["key2"] = "key2val"

	src := make(map[string]interface{})
	src["key3"] = 123
	src["key2"] = 120.012
	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = "key1val"
	expectedResult["key2"] = 120.012
	expectedResult["key3"] = 123
	assert.Equal(t, expectedResult, result)
}

// 4. dest complex value gets replaced with a simple src value
func TestRecursiveMerge_DestInterfaceSrcBasic(t *testing.T) {
	dest := make(map[string]interface{})
	dest["key1"] = map[string]interface{}{
		"name":  123,
		"value": "value",
	}
	dest["key2"] = "key2val"

	src := make(map[string]interface{})
	src["key3"] = 123
	src["key2"] = "key2val"
	src["key1"] = "testvalue-basictype"
	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = "testvalue-basictype"
	expectedResult["key2"] = "key2val"
	expectedResult["key3"] = 123
	assert.Equal(t, expectedResult, result)
}

// 5. nested object: only one of the child keys gets updated
func TestRecursiveMerge_DestInterfaceSrcInterface(t *testing.T) {
	dest := make(map[string]interface{})
	dest["key1"] = map[string]interface{}{
		"name":  123,
		"value": "value",
	}
	dest["key2"] = "key2val"

	src := make(map[string]interface{})
	src["key3"] = 123
	src["key2"] = "key2val"
	src["key1"] = map[string]interface{}{
		"name": 4567,
	}
	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = map[string]interface{}{
		"name":  4567,
		"value": "value",
	}
	expectedResult["key2"] = "key2val"
	expectedResult["key3"] = 123
	assert.Equal(t, expectedResult, result)
}

// 6. src and destination are identical
func TestRecursiveMerge_DestSrcIdentical(t *testing.T) {
	dest := make(map[string]interface{})
	dest["key1"] = map[string]interface{}{
		"name":  123,
		"value": "value",
	}
	dest["key2"] = "key2val"

	src := make(map[string]interface{})
	src["key1"] = map[string]interface{}{
		"name":  123,
		"value": "value",
	}
	src["key2"] = "key2val"

	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	expectedResult["key1"] = map[string]interface{}{
		"name":  123,
		"value": "value",
	}
	expectedResult["key2"] = "key2val"
	assert.Equal(t, expectedResult, result)
}

// 7. src and destination are empty
func TestRecursiveMerge_DestSrcEmpty(t *testing.T) {
	dest := make(map[string]interface{})
	src := make(map[string]interface{})

	merge := Merge{}
	result := merge.RecursiveMerge(src, dest)

	expectedResult := make(map[string]interface{})
	assert.Equal(t, expectedResult, result)
}
