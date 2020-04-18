package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	configRootRegistrationProvider = "spring.cloud.consul.discovery"
	RowSeparator                   = '~'
	FieldSeparator                 = ':'

	InstanceIdUuid     = "uuid"
	InstanceIdHostname = "hostname"

	ConfigKeyInfoAppName                  = "info.app.name"
	ConfigKeyInfoAppDescription           = "info.app.description"
	ConfigKeyInfoAppAttributesDisplayName = "info.app.attributes.displayName"
	ConfigKeyInfoAppAttributesParent      = "info.app.attributes.parent"
	ConfigKeyInfoAppAttributesType        = "info.app.attributes.type"
	ConfigKeyServerContextPath            = "server.contextPath"
	ConfigKeyServerSwaggerPath            = "swagger.ui.endpoint"
	ConfigKeyServerPort                   = "server.port"

	ConfigKeyInfoBuildVersion       = "info.build.version"
	ConfigKeyInfoBuildBuildNumber   = "info.build.buildNumber"
	ConfigKeyInfoBuildBuildDateTime = "info.build.buildDateTime"
)

var (
	logger      = log.NewLogger("msx.discovery.consulprovider")
	ErrDisabled = errors.New("Consul registration provider disabled")
)

type RegistrationProviderConfig struct {
	Enabled             bool          `config:"default=false"`
	Name                string        `config:"default=${info.app.name}"`
	Register            bool          `config:"default=true"`
	IpAddress           string        `config:"default="`
	Interface           string        `config:"default="`
	Port                int           `config:"default=0"`
	RegisterHealthCheck bool          `config:"default=true"`
	HealthCheckPath     string        `config:"default=/admin/health"`
	HealthCheckInterval time.Duration `config:"default=10s"`
	HealthCheckTimeout  time.Duration `config:"default=2s"`
	Tags                string        `config:"default="`
	InstanceId          string        `config:"default=local"` // uuid, hostname, or any static string
	InstanceName        string        `config:"default=${info.app.name}"`
}

type AppRegistrationDetails struct {
	ServiceAddress string
	ServicePort    string
	InstanceUuid   string
	InstanceId     string
	ContextPath    string
	SwaggerPath    string
	Name           string
	Application    string
	DisplayName    string
	Description    string
	Parent         string
	Type           string
	BuildVersion   string
	BuildDateTime  string
	BuildNumber    string
}

func (d AppRegistrationDetails) SocketAddress() string {
	return d.ServiceAddress + ":" + d.ServicePort
}

func (d AppRegistrationDetails) Tags() []string {
	return []string{
		"managedMicroservice",
		"contextPath=" + d.ContextPath,
		"swaggerPath=" + d.SwaggerPath,
		"instanceUuid=" + d.InstanceUuid,
		"name=" + d.DisplayName,
		"version=" + d.BuildVersion,
		"buildDateTime=" + d.BuildDateTime,
		"buildNumber=" + d.BuildNumber,
		"secure=false",
		"application=" + d.Application,
		"componentAttributes=" + marshalComponentAttributes(map[string]string{
			"serviceName": d.Name,
			"context":     d.contextPath(),
			"name":        d.DisplayName,
			"description": d.Description,
			"parent":      d.Parent,
			"type":        d.Type,
		}),
	}
}

func (d AppRegistrationDetails) Meta() map[string]string {
	return map[string]string{
		"buildDateTime": d.BuildDateTime,
		"buildNumber":   d.BuildNumber,
		"context":       d.contextPath(),
		"description":   d.Description,
		"instanceUuid":  d.InstanceUuid,
		"name":          d.DisplayName,
		"parent":        d.Parent,
		"serviceName":   d.Name,
		"type":          d.Type,
		"version":       d.BuildVersion,
		"application":   d.Application,
	}
}

func (d AppRegistrationDetails) contextPath() string {
	contextPath := path.Clean(d.ContextPath)
	if contextPath == "/" || contextPath == "." {
		contextPath = ""
	}
	return contextPath
}

type RegistrationProvider struct {
	config  *RegistrationProviderConfig
	conn    *consul.Connection
	details *AppRegistrationDetails
}

func (c *RegistrationProvider) tags() []string {
	var tags []string

	if c.config.Tags != "" {
		tags = append(tags, strings.Split(c.config.Tags, ",")...)
	}

	tags = append(tags, c.details.Tags()...)

	return tags
}

func (c *RegistrationProvider) meta() map[string]string {
	return c.details.Meta()
}

func (c *RegistrationProvider) healthCheck() *api.AgentServiceCheck {
	if c.config.RegisterHealthCheck == false {
		return nil
	}

	checkPath := path.Clean(path.Join(c.details.ContextPath, c.config.HealthCheckPath))

	return &api.AgentServiceCheck{
		Interval:      c.config.HealthCheckInterval.String(),
		Timeout:       c.config.HealthCheckTimeout.String(),
		TLSSkipVerify: true,
		HTTP:          fmt.Sprintf("%s://%s%s", "http", c.details.SocketAddress(), checkPath),
	}
}

func (c *RegistrationProvider) serviceRegistration() *api.AgentServiceRegistration {
	registration := &api.AgentServiceRegistration{
		ID:      c.details.InstanceId,
		Name:    c.details.Name,
		Address: c.details.ServiceAddress,
		Port:    c.config.Port,
		Tags:    c.tags(),
		Meta:    c.meta(),
	}

	if c.config.RegisterHealthCheck {
		registration.Check = c.healthCheck()
	}

	return registration
}

func (c *RegistrationProvider) Register(ctx context.Context) error {
	if c.details != nil {
		logger.Infof("Registering service in consul: %v", *c.details)
		return c.conn.RegisterService(ctx, c.serviceRegistration())
	} else {
		return nil
	}
}

func (c *RegistrationProvider) Deregister(ctx context.Context) error {
	if c.details != nil {
		logger.Infof("De-registering service in consul: %v", *c.details)
		return c.conn.DeregisterService(ctx, c.serviceRegistration())
	} else {
		return nil
	}
}

func marshalComponentAttributes(attributes map[string]string) string {
	var stringBuilder strings.Builder
	for k, v := range attributes {
		if stringBuilder.Len() > 0 {
			stringBuilder.WriteRune(RowSeparator)
		}
		stringBuilder.WriteString(k)
		stringBuilder.WriteRune(FieldSeparator)
		stringBuilder.WriteString(v)
	}

	return stringBuilder.String()
}

func detailsFromConfig(cfg *config.Config, rpConfig *RegistrationProviderConfig) (result *AppRegistrationDetails, err error) {
	result = &AppRegistrationDetails{}

	if rpConfig.IpAddress == "" {
		if result.ServiceAddress, err = getIp(rpConfig.Interface); err != nil {
			return nil, err
		}
	} else {
		result.ServiceAddress = rpConfig.IpAddress
	}

	if rpConfig.Port == 0 {
		if result.ServicePort, err = cfg.StringOr(ConfigKeyServerPort, strconv.Itoa(rpConfig.Port)); err != nil {
			return nil, err
		} else {
			rpConfig.Port, _ = strconv.Atoi(result.ServicePort)
		}
	} else {
		result.ServicePort = strconv.Itoa(rpConfig.Port)
	}

	if result.InstanceUuid, err = uuid.GenerateUUID(); err != nil {
		return nil, err
	}

	var instanceIdSuffix string
	switch rpConfig.InstanceId {
	case InstanceIdHostname:
		instanceIdSuffix, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	case InstanceIdUuid:
		instanceIdSuffix = result.InstanceUuid
	default:
		instanceIdSuffix = rpConfig.InstanceId
	}

	if result.ContextPath, err = cfg.String(ConfigKeyServerContextPath); err != nil {
		return nil, err
	}

	if result.SwaggerPath, err = cfg.StringOr(ConfigKeyServerSwaggerPath, "/swagger"); err != nil {
		return nil, err
	}

	result.Name = rpConfig.Name
	if result.Name == "" {
		return nil, errors.New("Registration name not configured")
	}

	var instanceIdPrefix string
	switch {
	case rpConfig.InstanceName != "":
		instanceIdPrefix = rpConfig.InstanceName
	case rpConfig.Name != "":
		instanceIdPrefix = rpConfig.Name
	}

	result.InstanceId = instanceIdPrefix + "-" + instanceIdSuffix

	if result.Application, err = cfg.String(ConfigKeyInfoAppName); err != nil {
		return nil, err
	}

	if result.Description, err = cfg.String(ConfigKeyInfoAppDescription); err != nil {
		return nil, err
	}

	if result.DisplayName, err = cfg.String(ConfigKeyInfoAppAttributesDisplayName); err != nil {
		return nil, err
	}

	if result.Parent, err = cfg.String(ConfigKeyInfoAppAttributesParent); err != nil {
		return nil, err
	}

	if result.Type, err = cfg.String(ConfigKeyInfoAppAttributesType); err != nil {
		return nil, err
	}

	if result.BuildVersion, err = cfg.String(ConfigKeyInfoBuildVersion); err != nil {
		return nil, err
	}

	if result.BuildNumber, err = cfg.String(ConfigKeyInfoBuildBuildNumber); err != nil {
		return nil, err
	}

	if result.BuildDateTime, err = cfg.String(ConfigKeyInfoBuildBuildDateTime); err != nil {
		return nil, err
	}

	return
}

func getIp(iface string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var ip net.IP
	for _, i := range ifaces {
		name := i.Name
		if iface != name && strings.Contains(name, "utun") {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					ip = v.IP
				}
			case *net.IPAddr:
				if v.IP.To4() != nil {
					ip = v.IP
				}
			}
		}
		if iface == name {
			break
		}
	}
	if ip == nil {
		if iface == "" {
			return "", errors.New("No valid interface or address found")
		} else {
			return "", errors.Errorf("Interface %s not found or no address", iface)
		}
	}

	return ip.String(), nil
}

func NewRegistrationProviderFromConfig(cfg *config.Config) (*RegistrationProvider, error) {
	var providerConfig = &RegistrationProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootRegistrationProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled || !providerConfig.Register {
		logger.Warn(ErrDisabled)
		return nil, nil
	}

	details, err := detailsFromConfig(cfg, providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create registration details")
	}

	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, nil
	}

	return &RegistrationProvider{
		config:  providerConfig,
		details: details,
		conn:    conn,
	}, nil
}
