// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"github.com/iancoleman/strcase"
	"strings"
)

func humanTypeName(messageTypeName string) string {
	cleanMessageTypeName := dispatchName(messageTypeName)
	spacedHumanName := strcase.ToDelimited(cleanMessageTypeName, ' ')
	return strings.Title(spacedHumanName)
}

func dispatchName(messageId string) string {
	var suffixes = []string{
		"Request",
		"Response",
		"Event",
		"Message",
	}

	for _, suffix := range suffixes {
		messageId = strings.TrimSuffix(messageId, suffix)
	}

	return messageId
}

func channelShortName(channel string) string {
	return strings.TrimSuffix(channel, "_TOPIC")
}

func messageName(channelShortName, suffix string) string {
	channelCamel := strcase.ToCamel(strings.ToLower(channelShortName))
	return channelCamel + strcase.ToCamel(suffix)
}

func operationName(channelShortName string, op string) string {
	channelCamel := strcase.ToCamel(strings.ToLower(channelShortName))
	return "on" + channelCamel + strcase.ToCamel(op)
}

func packageName(domain string) string {
	domainSnake := strcase.ToSnake(domain)
	return strings.ReplaceAll(domainSnake, "_", "")
}
