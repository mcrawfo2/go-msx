package types

import (
	"strconv"
	"strings"
)

type Version []int

func (v Version) Lt(other Version) bool {
	maxLen := len(v)
	if len(other) < maxLen {
		maxLen = len(other)
	}

	for n := 0; n < maxLen; n++ {
		if v[n] < other[n] {
			return true
		}
		if v[n] > other[n] {
			return false
		}
	}

	return len(v) < len(other)
}

func NewVersion(source string) Version {
	parts := strings.Split(source, ".")
	var numbers []int
	for _, part := range parts {
		if number, err := strconv.Atoi(part); err != nil {
			return Version{}
		} else {
			numbers = append(numbers, number)
		}
	}
	return numbers
}
