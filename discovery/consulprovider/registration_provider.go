package consulprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"time"
)

const (
	configRootRegistrationProvider = "spring.cloud.consul.discovery"
)

var (
	logger      = log.NewLogger("msx.discovery.consulprovider")
	ErrDisabled = errors.New("Consul registration provider disabled")
)

type RegistrationProviderConfig struct {
	Enabled             bool          `config:"default=false"`
	Port                int           `config:"default=0"`
	HealthCheckEnabled  bool          `config:"default=false"`
	HealthCheckPath     string        `config:"default=/admin/health"`
	HealthCheckInterval time.Duration `config:"default=10s"`
	HealthCheckTimeout  time.Duration `config:"default=2s"`
}

type RegistrationProvider struct {
	config  *RegistrationProviderConfig
	conn    *consul.Connection
	service *api.AgentServiceRegistration
}

func (c *RegistrationProvider) Register() error {
	return nil

	//id := uuid.Must(uuid.NewV4())
	//check := api.AgentServiceCheck{
	//	Interval:      checkInterval,
	//	Timeout:       checkTimeout,
	//	TLSSkipVerify: checkTLSSkipVerify,
	//	HTTP:          "http://" + ipaddress + ":" + port + heathCheckPath,
	//}
	//
	//iport, _ := strconv.Atoi(port)
	//
	//var importedTags []string
	//tags, err := vmsglobal.GetApplicationMap(vmsconstants.ApplicationMapSpringCloudConsulDiscoveryTags)
	//if err == nil && tags != "" {
	//	importedTags = strings.Split(tags, ",")
	//}
	//
	//buildInfo, err := readBuildInfo()
	//if err != nil {
	//	vmsutil.Error.Print(err.Error())
	//	buildInfo = &BuildInfo{}
	//}
	//
	//tagsToUse := []string{
	//	"managedMicroservice",
	//	"contextPath=" + contextPath,
	//	"instanceUuid=" + id.String(),
	//	"name=" + appDisplayName,
	//	"version=" + buildInfo.Version,
	//	"buildDateTime=" + buildInfo.BuildDateTime,
	//	"buildNumber=" + buildInfo.BuildNumber,
	//	"componentAttributes=" + marshalComponentAttributes(map[string]string{
	//		"serviceName": appName,
	//		"context":     strings.TrimPrefix(contextPath, "/"),
	//		"name":        appDisplayName,
	//		"description": appDescription,
	//		"parent":      appParent,
	//		"type":        appType,
	//	}),
	//}
	//
	//tagsToUse = append(tagsToUse, importedTags...)
	//
	//// Service name is the name assigned to this service
	////
	//pc.Service = api.AgentServiceRegistration{
	//	Name:    appName,
	//	Address: ipaddress,
	//	Check:   &check,
	//	Port:    iport,
	//	ID:      appName + "-" + id.String(),
	//	Tags:    tagsToUse,
	//}

}

func NewRegistrationProviderFromConfig(cfg *config.Config) (*RegistrationProvider, error) {
	var providerConfig = &RegistrationProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootRegistrationProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled {
		logger.Warn(ErrDisabled)
		return nil, nil
	}

	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, nil
	}

	return &RegistrationProvider{
		config: providerConfig,
		conn:   conn,
	}, nil
}
