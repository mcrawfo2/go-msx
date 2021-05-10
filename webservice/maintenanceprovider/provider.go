package maintenanceprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
)

const endpointName = "maintenance"

type MaintenanceProvider struct{}

func (h MaintenanceProvider) updateMaintenance(req *restful.Request) (interface{}, error) {
	maintenanceUpdateReq := MaintenanceUpdate{}
	bodyBytes, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &maintenanceUpdateReq)
	if err != nil {
		return nil, err
	}
	err =  validate.Validate(&maintenanceUpdateReq)
	if err != nil {
		return nil, err
	}
	resp := MaintenanceResponse{}
	var maintenanceDetail []MaintenanceTask
	if maintenanceUpdateReq.Mode == "NORMAL" {
		resp.Mode = maintenanceUpdateReq.Mode
		task := MaintenanceTask{
			Task:    "REST API Maintenance",
			Mode:    resp.Mode,
			Message: "Resumed processing API requests.",
		}
		maintenanceDetail = append(maintenanceDetail, task)

		task = MaintenanceTask{
			Task:    "Message Consumer Maintenance",
			Mode:    resp.Mode,
			Message: "Resumed consumer bindings: [].",
		}
		maintenanceDetail = append(maintenanceDetail, task)
		resp.Detail = maintenanceDetail

	} else if maintenanceUpdateReq.Mode == "MAINTENANCE" {
		resp.Mode = maintenanceUpdateReq.Mode
		task := MaintenanceTask{
			Task:    "REST API Maintenance",
			Mode:    resp.Mode,
			Message: "Stopped processing API requests. Excluding: [/admin/** ALL]",
		}
		maintenanceDetail = append(maintenanceDetail, task)

		task = MaintenanceTask{
			Task:    "Message Consumer Maintenance",
			Mode:    resp.Mode,
			Message: "Paused consumer bindings: [].",
		}
		maintenanceDetail = append(maintenanceDetail, task)
		resp.Detail = maintenanceDetail
	}

	return resp, nil
}


func (h MaintenanceProvider) EndpointName() string {
	return endpointName
}

func (h MaintenanceProvider) Actuate(healthService *restful.WebService) error {
	healthService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	healthService.Produces(restful.MIME_JSON, restful.MIME_XML)

	healthService.Path(healthService.RootPath() + "/admin/maintenance")

	healthService.Route(healthService.PUT("").
		Operation("admin.maintenance").
		To(adminprovider.RawAdminController(h.updateMaintenance)).
		Doc("Set Maintenance").
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(MaintenanceProvider))
		adminprovider.RegisterLink("maintenance", "maintenance", false)
	}
	return nil
}
