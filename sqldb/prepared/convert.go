package prepared

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/google/uuid"
	"time"
)

func ToModelUuid(v types.UUID) uuid.UUID {
	return uuid.MustParse(v.String())
}

func ToOptionalModelUuid(v *types.UUID) NullUUID {
	if v == nil {
		return NullUUID{}
	}

	return NullUUID{
		valid: true,
		value: v.ToByteArray(),
	}
}

func ToApiUuid(v uuid.UUID) types.UUID {
	return types.MustParseUUID(v.String())
}

func ToOptionalApiUuid(v NullUUID) *types.UUID {
	if !v.valid {
		return nil
	}

	uuidValue := v.UUID()
	value := types.UUID(uuidValue[:])
	return &value
}

func ToApiTime(v time.Time) types.Time {
	return types.Time(v)
}

func ToOptionalApiTime(v *time.Time) *types.Time {
	if v == nil {
		return nil
	}

	r := ToApiTime(*v)
	return &r
}

func ToModelTime(v types.Time) time.Time {
	return v.ToTimeTime()
}

func ToOptionalModelTime(v *types.Time) *time.Time {
	if v == nil {
		return nil
	}

	r := ToModelTime(*v)
	return &r
}

func ToApiIp(i Ip) string {
	return i.String()
}

func ToModelIp(ip string) Ip {
	i, _ := NewIp(ip)
	return i
}

func ToApiCidr(c Cidr) string {
	return c.String()
}

func ToModelCidr(cidr string) Cidr {
	c, _ := NewCidr(cidr)
	return c
}

func ToApiIpPort(i IpPort) string {
	return i.String()
}

func ToModelIpPort(ipPort string) IpPort {
	i, _ := NewIpPort(ipPort)
	return i
}
