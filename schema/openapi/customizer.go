// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	openapi "github.com/swaggest/openapi-go/openapi3"
)

type SpecificationCustomizer struct {
	appInfo   schema.AppInfo
	serverUrl string
	version   string
}

func (c SpecificationCustomizer) PostBuildSpec(spec *openapi.Spec) error {
	return types.ErrorList{
		c.PostBuildVersion(spec),
		c.PostBuildInfo(spec),
		c.PostBuildServers(spec),
		c.PostBuildTags(spec),
	}.Filter()
}

func (c SpecificationCustomizer) PostBuildVersion(s *openapi.Spec) error {
	switch c.version {
	case "3":
		s.Openapi = "3.0.3"
	default:
		return errors.Errorf("Unknown OpenApi version %q", c.version)
	}
	return nil
}

func (c SpecificationCustomizer) PostBuildInfo(spec *openapi.Spec) error {
	spec.Info = openapi.Info{
		Title: "MSX API Documentation for " + c.appInfo.Name,
		Description: types.NewStringPtr("<h3>This is the REST API documentation for " + c.appInfo.Name + "</h3>\n \n" +
			c.appInfo.Description + "\n" +
			"+ API Authorization \n" +
			"    + Authorization header is <b>required</b>. \n" +
			"    + It should be in Bearer authentication scheme </br>(e.g <b> Authorization: BEARER &lt;access token&gt; </b>)\n"),
		TermsOfService: types.NewStringPtr("http://www.cisco.com"),
		Contact: &openapi.Contact{
			Name:  types.NewStringPtr("Cisco Systems Inc."),
			URL:   types.NewStringPtr("http://www.cisco.com"),
			Email: types.NewStringPtr("somecontact@cisco.com"),
		},
		License: &openapi.License{
			Name: "Apache License Version 2.0",
			URL:  types.NewStringPtr("http://www.apache.org/licenses/LICENSE-2.0.html"),
		},
		Version: c.appInfo.Version,
	}
	return nil
}

func (c SpecificationCustomizer) PostBuildServers(spec *openapi.Spec) error {
	spec.Servers = []openapi.Server{
		{
			URL: c.serverUrl,
		},
	}
	return nil
}

func (c SpecificationCustomizer) PostBuildTags(spec *openapi.Spec) error {
	return nil
}
