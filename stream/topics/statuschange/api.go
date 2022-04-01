// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package statuschange

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
)

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) (err error) {
	msg, err := producer.Message(ctx)
	if err != nil {
		return err
	}

	return Publish(ctx, msg)
}

func Publish(ctx context.Context, message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return stream.Publish(ctx, TopicStatusChange, bytes, nil)
}
