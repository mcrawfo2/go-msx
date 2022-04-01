// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package trace

const (
	FieldHttpCode      = "http.code"
	FieldHttpUrl       = "http.url"
	FieldHttpMethod    = "http.method"
	FieldOperation     = "operation"
	FieldKeyspace      = "keyspace"
	FieldError         = "error"
	FieldTopic         = "stream.topic"
	FieldTransport     = "stream.transport"
	FieldDirection     = "stream.direction"
	FieldSpanKind      = "span.kind"
	FieldSpanType      = "span.type"
	FieldStatus        = "status"
	FieldDeviceId      = "beat.device.id"
	FieldDeviceAddress = "beat.device.address"
	FieldServiceId     = "beat.service.id"

	SpanKindProducer = "producer"
	SpanKindConsumer = "consumer"
	SpanKindClient   = "client"
	SpanKindServer   = "server"

	RefChildOf     = "childOf"
	RefFollowsFrom = "followsFrom"
)
