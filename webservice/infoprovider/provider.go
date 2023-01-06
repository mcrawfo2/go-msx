// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package infoprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"runtime/debug"
)

const (
	configKeyInfo = "info"
	endpointName  = "info"
)

type GoModule struct {
	Path    string    `json:"path"`     // module path
	Version string    `json:"version"`  // module version
	Sum     string    `json:"checksum"` // checksum
	Replace *GoModule `json:"replace,omitempty"`
}

type GoInfo struct {
	GoVersion string            `json:"version"`        // Version of Go that produced this binary.
	Path      string            `json:"path"`           // The main package path
	Main      *GoModule         `json:"main,omitempty"` // The module containing the main package
	Settings  map[string]string `json:"settings"`
}

type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version" config:"default="`
	Description string `json:"description"`
	Attributes  struct {
		DisplayName string `json:"displayName"`
		Parent      string `json:"parent"`
		Type        string `json:"type"`
	} `json:"attributes"`
}

type BuildInfo struct {
	Version       string       `json:"version"`
	BuildNumber   string       `json:"number" config:"buildNumber"`
	BuildDateTime string       `json:"-"`
	Artifact      string       `json:"artifact"`
	Name          string       `json:"name"`
	Time          epochSeconds `json:"time" config:"default=0"`
	Group         string       `json:"group"`
	CommitHash    string       `json:"commitHash" config:"default="`
	DiffHash      string       `json:"diffHash" config:"default="`
}

type InfoProvider struct{}

func (h InfoProvider) infoReport(req *restful.Request) (interface{}, error) {
	type Info struct {
		App   AppInfo   `json:"app"`
		Build BuildInfo `json:"build"`
		Go    *GoInfo   `json:"go,omitempty" config:"-"`
	}

	i := Info{}
	if err := config.MustFromContext(req.Request.Context()).Populate(&i, configKeyInfo); err != nil {
		return nil, webservice.NewStatusError(err, 500)
	}

	i.App.Version = i.Build.Version
	buildTime, err := types.ParseTime(i.Build.BuildDateTime)
	if err == nil {
		i.Build.Time = newEpochSeconds(buildTime.ToTimeTime())
	}

	goBuildInfo, ok := debug.ReadBuildInfo()
	if ok {
		i.Go = new(GoInfo)
		i.Go.GoVersion = goBuildInfo.GoVersion
		i.Go.Path = goBuildInfo.Path
		if goBuildInfo.Main.Path != "" {
			i.Go.Main = &GoModule{
				Path:    goBuildInfo.Main.Path,
				Version: goBuildInfo.Main.Version,
				Sum:     goBuildInfo.Main.Sum,
			}
		}

		i.Go.Settings = make(map[string]string)
		for _, setting := range goBuildInfo.Settings {
			i.Go.Settings[setting.Key] = setting.Value
		}
	}

	return i, nil
}

func (h InfoProvider) EndpointName() string {
	return endpointName
}

func (h InfoProvider) Actuate(infoService *restful.WebService) error {
	infoService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	infoService.Produces(restful.MIME_JSON, restful.MIME_XML)

	infoService.Path(infoService.RootPath() + "/admin/info")

	// Unsecured routes for info
	infoService.Route(infoService.GET("").
		Operation("admin.info").
		To(adminprovider.RawAdminController(h.infoReport)).
		Doc("Get System info").
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(InfoProvider))
		adminprovider.RegisterLink("info", "info", false)
	}
	return nil
}
