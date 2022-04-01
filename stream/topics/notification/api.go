// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package notification

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

const TopicName = "NOTIFICATION_TOPIC"

func Publish(ctx context.Context, message Message) error {
	return stream.PublishObject(ctx, TopicName, message, nil)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) error {
	message, err := producer.Message(ctx)
	if err != nil {
		return err
	}
	return Publish(ctx, message)
}
