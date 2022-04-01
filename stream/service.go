// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"encoding/json"
)

type PublisherService interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error)
	PublishObject(ctx context.Context, topic string, payload interface{}, metadata map[string]string) (err error)
}

type publisherFunc func(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error)

func (p publisherFunc) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error) {
	return p(ctx, topic, payload, metadata)
}

func (p publisherFunc) PublishObject(ctx context.Context, topic string, payload interface{}, metadata map[string]string) (err error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p(ctx, topic, payloadBytes, metadata)
}

var ProductionPublisherService PublisherService = publisherFunc(Publish)

func NewPublisherService(ctx context.Context) PublisherService {
	service := PublisherServiceFromContext(ctx)
	if service == nil {
		service = ProductionPublisherService
	}
	return service
}
