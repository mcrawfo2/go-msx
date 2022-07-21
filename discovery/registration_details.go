package discovery

import (
	"path"
	"strings"
)

const (
	RowSeparator   = '~'
	FieldSeparator = ':'
)

type RegistrationDetails struct {
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

func (d RegistrationDetails) SocketAddress() string {
	return d.ServiceAddress + ":" + d.ServicePort
}

func (d RegistrationDetails) Tags() []string {
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

func (d RegistrationDetails) Meta() map[string]string {
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

func (d RegistrationDetails) contextPath() string {
	contextPath := path.Clean(d.ContextPath)
	if contextPath == "/" || contextPath == "." {
		contextPath = ""
	}
	return contextPath
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
