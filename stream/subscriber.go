// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}
