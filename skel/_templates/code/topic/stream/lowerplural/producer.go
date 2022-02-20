package lowerplural

import "cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api"

type UpperCamelSingularMessageProducer interface {
	Produce() api.UpperCamelSingularMessage
}
