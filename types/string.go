package types

import (
	"math/rand"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const suffixCharset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = suffixCharset[seededRand.Intn(len(suffixCharset))]
	}
	return string(b)
}

type StringSlice []string

func (ss StringSlice) ToInterfaceSlice() (target []interface{}) {
	for _, s := range ss {
		target = append(target, s)
	}
	return
}
