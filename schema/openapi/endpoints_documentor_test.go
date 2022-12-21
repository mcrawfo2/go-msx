// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAppInfo = schema.AppInfo{
	Name:        "TestApp",
	DisplayName: "Test Application",
	Description: "The Test Application",
	Version:     "5.0.0",
}

func TestEndpointsDocumentor_Document(t *testing.T) {
	type args struct {
		endpoints []*restops.Endpoint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SingleEndpoint",
			args: args{
				endpoints: []*restops.Endpoint{
					{
						Method:      "GET",
						Path:        "/api/v2/single",
						OperationID: "getSingle",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "MultipleEndpoints",
			args: args{
				endpoints: []*restops.Endpoint{
					{
						Method:      "GET",
						Path:        "/api/v2/single",
						OperationID: "getSingle",
					},
					{
						Method:      "GET",
						Path:        "/api/v2/single/{singleId}",
						OperationID: "getSingle",
						Request: restops.EndpointRequest{
							Parameters: []restops.EndpointRequestParameter{
								{
									In:   "path",
									Name: "singleId",
								},
							},
						},
					},
					{
						Method:      "PUT",
						Path:        "/api/v2/single/{singleId:bob}",
						OperationID: "putSingle",
						Request: restops.EndpointRequest{
							Parameters: []restops.EndpointRequestParameter{
								{
									In:   "path",
									Name: "singleId",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DuplicateParameter",
			args: args{
				endpoints: []*restops.Endpoint{
					{
						Method:      "GET",
						Path:        "/api/v2/single/{singleId}",
						OperationID: "getSingle",
						Request: restops.EndpointRequest{
							Parameters: []restops.EndpointRequestParameter{
								{
									In:   "path",
									Name: "singleId",
								},
								{
									In:   "path",
									Name: "singleId",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MissingParameter",
			args: args{
				endpoints: []*restops.Endpoint{
					{
						Method:      "GET",
						Path:        "/api/v2/single",
						OperationID: "getSingle",
						Request: restops.EndpointRequest{
							Parameters: []restops.EndpointRequestParameter{
								{
									In:   "path",
									Name: "singleId",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewEndpointsDocumentor(&testAppInfo, "http://localhost:30303", "3")

			err := d.Document(tt.args.endpoints, "", "")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
