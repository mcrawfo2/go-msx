package api

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type ${async.upmsgtype} struct {
	Id types.UUID `json:"id"`
	Timestamp types.Time `json:"timestamp"`
}
