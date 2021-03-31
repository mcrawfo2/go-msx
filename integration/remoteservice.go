package integration

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const (
	remoteServiceConfigRoot = "remoteservice"

	ServiceNameAdministration   = "administrationservice"
	ServiceNameAlert            = "alertservice"
	ServiceNameAuditing         = "auditingservice"
	ServiceNameBilling          = "billingservice"
	ServiceNameConsume          = "consumeservice"
	ServiceNameManage           = "manageservice"
	ServiceNameMonitor          = "monitorservice"
	ServiceNameNotification     = "notificationservice"
	ServiceNameOrchestration    = "orchestrationservice"
	ServiceNameRouter           = "routerservice"
	ServiceNameServiceConfig    = "serviceconfigmanager"
	ServiceNameServiceExtension = "serviceextensionservice"
	ServiceNameUserManagement   = "usermanagementservice"
	ServiceNameWorkflow         = "workflowservice"
	ServiceNameIpam             = "ipamservice"

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

type RemoteServiceConfig struct {
	ServiceName string `config:"key=service,default="`
}

func NewRemoteServiceConfig(ctx context.Context, serviceName string) *RemoteServiceConfig {
	cfg := config.FromContext(ctx)
	var remoteService RemoteServiceConfig
	if err := cfg.Populate(&remoteService, config.PrefixWithName(remoteServiceConfigRoot, serviceName)); err != nil {
		logger.Errorf("Not able to load remote service config for `%s`", serviceName)
	}
	if len(remoteService.ServiceName) == 0 {
		remoteService.ServiceName = serviceName
	}
	return &remoteService
}
