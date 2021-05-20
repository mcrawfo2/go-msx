package prepared

import (
	"database/sql/driver"
	"fmt"
	"net"
	"strings"
)

type Ip net.IP

func (i Ip) String() string {
	value, _ := i.MarshalText()
	return value
}

func (i *Ip) UnmarshalText(data string) error {
	ip := net.ParseIP(data)
	if strings.Contains(data, ".") {
		ip = ip.To4()
	}
	*i = Ip(ip)
	return nil
}

func (i Ip) MarshalText() (string, error) {
	return net.IP(i).String(), nil
}

func (i Ip) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	return i.MarshalText()
}

func (i *Ip) Scan(src interface{}) error {
	switch tsrc := src.(type) {
	case []byte:
		return i.UnmarshalText(string(tsrc))
	case string:
		return i.UnmarshalText(tsrc)
	case nil:
		*i = Ip(net.IPv4zero)
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Ip", src)
}

func (i Ip) Add(count int) Ip {
	resultBytes := i[:]
	for n := len(resultBytes) - 1; n >= 0; n-- {
		count += int(resultBytes[n])
		resultBytes[n] = byte(count & 0xFF)
		count >>= 8
	}
	return resultBytes
}

func (i Ip) Equals(o Ip) bool {
	for n := 0; n < len(i) && n < len(o); n++ {
		if i[n] != o[n] {
			return false
		}
	}
	return len(i) == len(o)
}

func NewIp(ip string) (Ip, error) {
	var i = new(Ip)
	err := i.UnmarshalText(ip)
	return *i, err
}

// MustNewIp ignores errors returned by NewIp.
// Should be used only for known-good inputs.
func MustNewIp(ip string) Ip {
	result, _ := NewIp(ip)
	return result
}
