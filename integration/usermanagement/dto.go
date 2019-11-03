package usermanagement

type UserCapabilityListDTO struct {
	Capabilities []UserCapabilityDTO `json:"capabilities"`
}

type UserCapabilityDTO struct {
	Name string `json:"name"`
}
