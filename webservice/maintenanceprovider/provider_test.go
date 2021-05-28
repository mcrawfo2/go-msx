package maintenanceprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func mockMaintenanceResponse() MaintenanceResponse {
	var maintenanceDetail []MaintenanceTask
	const ModeMaintenance = "MAINTENANCE"
	task := MaintenanceTask{
		Task:    "REST API Maintenance",
		Mode:    ModeMaintenance,
		Message: "Stopped processing API requests. Excluding: [/admin/** ALL]",
	}
	maintenanceDetail = append(maintenanceDetail, task)

	task = MaintenanceTask{
		Task:    "Message Consumer Maintenance",
		Mode:    ModeMaintenance,
		Message: "Paused consumer bindings: [].",
	}
	maintenanceDetail = append(maintenanceDetail, task)

	return MaintenanceResponse{
		Mode:   ModeMaintenance,
		Detail: maintenanceDetail,
	}
}

func mockMaintenanceRequestBody() *restful.Request {
	maintenanceUpdate := MaintenanceUpdate{
		Mode: "MAINTENANCE",
	}

	normalBodyBytes, _ := json.Marshal(maintenanceUpdate)
	normalBodyIOReaderCloser := ioutil.NopCloser(bytes.NewReader(normalBodyBytes))

	return &restful.Request{
		Request: &http.Request{
			RemoteAddr: "10.10.10.10",
			Proto:      "http",
			Method:     "PUT",
			Body:       normalBodyIOReaderCloser,
		},
	}
}

func mockNormalRequestBody() *restful.Request {
	maintenanceUpdate := MaintenanceUpdate{
		Mode: "NORMAL",
	}

	normalBodyBytes, _ := json.Marshal(maintenanceUpdate)
	normalBodyIOReaderCloser := ioutil.NopCloser(bytes.NewReader(normalBodyBytes))

	return &restful.Request{
		Request: &http.Request{
			RemoteAddr: "10.10.10.10",
			Proto:      "http",
			Method:     "PUT",
			Body:       normalBodyIOReaderCloser,
		},
	}

}

func mockInvalidModeRequestBody() *restful.Request {
	maintenanceUpdate := MaintenanceUpdate{
		Mode: "Invalid",
	}

	normalBodyBytes, _ := json.Marshal(maintenanceUpdate)
	normalBodyIOReaderCloser := ioutil.NopCloser(bytes.NewReader(normalBodyBytes))

	return &restful.Request{
		Request: &http.Request{
			RemoteAddr: "10.10.10.10",
			Proto:      "http",
			Method:     "PUT",
			Body:       normalBodyIOReaderCloser,
		},
	}

}

func mockInvalidRequestBody() *restful.Request {
	maintenanceUpdate := "Invalid Request"

	normalBodyBytes, _ := json.Marshal(maintenanceUpdate)
	normalBodyIOReaderCloser := ioutil.NopCloser(bytes.NewReader(normalBodyBytes))

	return &restful.Request{
		Request: &http.Request{
			RemoteAddr: "10.10.10.10",
			Proto:      "http",
			Method:     "PUT",
			Body:       normalBodyIOReaderCloser,
		},
	}

}

func mockNormalResponse() MaintenanceResponse {
	var maintenanceDetail []MaintenanceTask
	const ModeNormal = "NORMAL"
	task := MaintenanceTask{
		Task:    "REST API Maintenance",
		Mode:    ModeNormal,
		Message: "Resumed processing API requests.",
	}
	maintenanceDetail = append(maintenanceDetail, task)

	task = MaintenanceTask{
		Task:    "Message Consumer Maintenance",
		Mode:    ModeNormal,
		Message: "Resumed consumer bindings: [].",
	}
	maintenanceDetail = append(maintenanceDetail, task)

	return MaintenanceResponse{
		Mode:   ModeNormal,
		Detail: maintenanceDetail,
	}
}
func TestMaintenanceProvider_Actuate(t *testing.T) {
	type args struct {
		healthService *restful.WebService
	}
	testArgs := args{
		healthService: new(restful.WebService),
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test",
			args:    testArgs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := MaintenanceProvider{}
			if err := h.Actuate(tt.args.healthService); (err != nil) != tt.wantErr {
				t.Errorf("Actuate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaintenanceProvider_EndpointName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: endpointName,
			want: endpointName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := MaintenanceProvider{}
			if got := h.EndpointName(); got != tt.want {
				t.Errorf("EndpointName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaintenanceProvider_updateMaintenance(t *testing.T) {
	type args struct {
		req *restful.Request
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "MODE_MAINTENANCE",
			args: args{
				req: mockMaintenanceRequestBody(),
			},
			want:    mockMaintenanceResponse(),
			wantErr: false,
		},
		{
			name:    "MODE_NORMAL",
			args:    args{req: mockNormalRequestBody()},
			want:    mockNormalResponse(),
			wantErr: false,
		},
		{
			name: "MODE_INVALID",
			args: args{
				req: mockInvalidModeRequestBody(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "INVALID_REQUEST",
			args: args{
				req: mockInvalidRequestBody(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := MaintenanceProvider{}
			got, err := h.updateMaintenance(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateMaintenance() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateMaintenance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterProvider(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: endpointName,
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterProvider(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RegisterProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
