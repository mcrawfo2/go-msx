package prepared

import (
	"database/sql/driver"
	"fmt"
	"net"
	"strconv"
)

type IpPort struct {
	Ip   Ip
	Port int
}

func (i IpPort) String() string {
	v, _ := i.MarshalText()
	return v
}

func (i *IpPort) UnmarshalText(s string) error {
	ip, port, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}

	err = i.Ip.UnmarshalText(ip)
	if err != nil {
		return err
	}

	i.Port, err = strconv.Atoi(port)
	if err != nil {
		return err
	}

	return nil
}

func (i IpPort) MarshalText() (string, error) {
	ip, err := i.Ip.MarshalText()
	if err != nil {
		return "", err
	}

	port := strconv.Itoa(i.Port)

	return net.JoinHostPort(ip, port), nil
}

func (i IpPort) Value() (driver.Value, error) {
	return i.MarshalText()
}

func (i *IpPort) Scan(src interface{}) error {
	switch tsrc := src.(type) {
	case []byte:
		return i.UnmarshalText(string(tsrc))
	case string:
		return i.UnmarshalText(tsrc)
	case nil:
		*i = IpPort{
			Ip:   Ip(net.IPv4zero),
			Port: 0,
		}
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to IpPort", src)
}

func NewIpPort(ipPort string) (IpPort, error) {
	var i IpPort
	err := i.UnmarshalText(ipPort)
	return i, err
}
