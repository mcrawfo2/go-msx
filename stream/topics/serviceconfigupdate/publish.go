// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package serviceconfigupdate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics/serviceconfigevent"
)

func PublishApplicationStatusUpdateRequest(ctx context.Context, request ApplicationStatusUpdateRequest) error {
	return stream.PublishObject(ctx,
		TopicServiceConfigUpdateTopic,
		actionRequest{Request: request},
		map[string]string{
			MetaDataEventType: EventTypeApplicationStatusUpdate,
		})
}

func PublishApplicationResult(ctx context.Context, sourceEvent serviceconfigevent.ApplicationEvent, err error) error {
	var status string
	var statusDetails string

	if err != nil {
		switch sourceEvent.EventType() {
		case serviceconfigevent.EventTypeApplicationCreated:
			status = ApplicationStatusFailed
			statusDetails = err.Error()

		case serviceconfigevent.EventTypeApplicationDeleted:
			status = ApplicationStatusDeleteFailed
			statusDetails = err.Error()
		}
	} else {
		switch sourceEvent.EventType() {
		case serviceconfigevent.EventTypeApplicationCreated:
			status = ApplicationStatusSuccess

		case serviceconfigevent.EventTypeApplicationDeleted:
			status = ApplicationStatusDeleted
		}
	}

	request := ApplicationStatusUpdateRequest{
		ApplicationId:   sourceEvent.ApplicationId,
		ServiceConfigId: sourceEvent.ServiceConfiguration.ServiceConfigID,
		Service:         sourceEvent.ServiceConfiguration.Service,
		Status:          status,
		StatusDetails:   statusDetails,
	}

	return PublishApplicationStatusUpdateRequest(ctx, request)
}
