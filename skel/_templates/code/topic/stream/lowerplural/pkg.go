package lowerplural

import "cto-github.cisco.com/NFV-BU/go-msx/log"

type contextKey string

var logger = log.NewLogger("${app.name}.internal.stream.lowerplural")

const (
	topicUpperCamelSingular = "SCREAMING_SNAKE_SINGULAR_TOPIC"
)
