// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package serviceconfigevent

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
)

func Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return stream.Publish(ctx, TopicServiceConfigEventTopic, payload, event.Headers)
}
