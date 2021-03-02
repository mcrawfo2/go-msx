package api

type DeviceDTO struct {
	// Deprecated
	Id        string           `json:"id"`
	DeviceId  string           `json:"deviceId"`
	ServiceId string           `json:"serviceId"`
	Ip        string           `json:"ip"`
	Metrics   DeviceMetricsDTO `json:"metrics,omitempty"`
	Tags      []string         `json:"tags,omitempty"`
}

type DeviceMetricsDTO struct {
	ProtocolUpperCamel BeatMetricsDTO `json:"protocolLowerCamel"`
}

type BeatMetricsDTO struct {
	// TODO: add beat device config
}
