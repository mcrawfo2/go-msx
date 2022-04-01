// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package metricsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
)

const (
	endpointName = "metrics"
)

var metricNames = []string{}

type Report struct {
	Names []string `json:"names"`
}

type Provider struct{}

func (h Provider) Report(req *restful.Request) (interface{}, error) {
	return Report{Names: metricNames}, nil
}

func (h Provider) EndpointName() string {
	return endpointName
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.metrics").
		To(adminprovider.RawAdminController(h.Report)).
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
