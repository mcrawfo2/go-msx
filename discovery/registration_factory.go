package discovery

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
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
	ConfigKeyNetworkAddress               = "network.outbound.address"

	ConfigKeyInfoBuildVersion       = "info.build.version"
	ConfigKeyInfoBuildBuildNumber   = "info.build.buildNumber"
	ConfigKeyInfoBuildBuildDateTime = "info.build.buildDateTime"
)

type RegistrationFactory struct{}

func (f RegistrationFactory) NewRegistrationDetails(ctx context.Context, rpConfig *RegistrationConfig) (result *RegistrationDetails, err error) {
	result = &RegistrationDetails{}
	cfg := config.FromContext(ctx)

	if rpConfig.IpAddress == "" {
		if rpConfig.Interface == "" {
			if result.ServiceAddress, err = cfg.String(ConfigKeyNetworkAddress); err != nil {
				return nil, err
			}
		} else {
			if result.ServiceAddress, err = getIp(rpConfig.Interface); err != nil {
				return nil, err
			}
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
		switch {
		case iface == name:
			// Keep it
		case strings.Contains(name, "utun"),
			strings.Contains(name, "vEthernet"),
			strings.Contains(name, "Loopback"):
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
