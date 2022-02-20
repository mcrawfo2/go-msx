package api

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type UpperCamelSingularMessage struct {
	Id   types.UUID `json:"id"`
	Data string     `json:"data"`
}
