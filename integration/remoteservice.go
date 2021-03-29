package integration

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const (
	remoteServiceConfigRoot = "remoteservice"
)

type RemoteServiceConfig struct {
	ServiceName string   `config:"key=service,default="`
}


func NewRemoteServiceConfig(ctx context.Context, serviceName ServiceName) *RemoteServiceConfig {
	cfg := config.FromContext(ctx)
	var remoteService RemoteServiceConfig
	if err := cfg.Populate(&remoteService, remoteServiceConfigRoot + "." + string(serviceName)); err != nil {
		logger.Errorf("Not able to load remote service config for `%s`", string(serviceName))
	}
	if len(remoteService.ServiceName) == 0 {
		remoteService.ServiceName = string(serviceName)
	}
	return &remoteService
}

type ServiceName string

const (
	ServiceNameAdministration  ServiceName = "administrationservice"
	ServiceNameAlert           ServiceName = "alertservice"
	ServiceNameAuditing        ServiceName = "auditingservice"
	ServiceNameBilling         ServiceName = "billingservice"
	ServiceNameConsume         ServiceName = "consumeservice"
	ServiceNameManage          ServiceName = "manageservice"
	ServiceNameMonitor          ServiceName = "monitorservice"
	ServiceNameNotification     ServiceName = "notificationservice"
	ServiceNameOrchestration    ServiceName = "orchestrationservice"
	ServiceNameRouter           ServiceName = "routerservice"
	ServiceNameServiceConfig    ServiceName = "templateservice"
	ServiceNameServiceExtension ServiceName = "serviceextensionservice"
	ServiceNameUserManagement   ServiceName = "usermanagementservice"
	ServiceNameWorkflow         ServiceName = "workflowservice"
	ServiceNameIpam             ServiceName = "ipamservice"

	ServiceNameDnacBeat    = "probe_dnac"
	ServiceNameEncsBeat    = "probe_encs"
	ServiceNameHeartBeat   = "probe_ping"
	ServiceNameMerakiBeat  = "probe_meraki"
	ServiceNameSnmpBeat    = "probe_snmp"
	ServiceNameSshBeat     = "probe_ssh"
	ServiceNameViptelaBeat = "probe_viptela"

	ServiceNameResource      = "resourceservice"
	ServiceNameManagedDevice = "manageddeviceservice"

	ResourceProviderNameDnac    = "dnac"
	ResourceProviderNameAws     = "aws"
	ResourceProviderNameViptela = "viptelaservice"
)
