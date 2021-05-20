package prepared

import (
	"database/sql/driver"
	"fmt"
	"github.com/lib/pq"
	"net"
)

type Cidr net.IPNet

func (c Cidr) String() string {
	v, _ := c.MarshalText()
	return v
}

func (c *Cidr) UnmarshalText(data string) error {
	_, ipNet, err := net.ParseCIDR(data)
	if err != nil {
		return err
	}
	c.IP = ipNet.IP
	c.Mask = ipNet.Mask
	return nil
}

func (c Cidr) MarshalText() (string, error) {
	ipAddr := c.IP.String()
	prefixLen, _ := c.Mask.Size()
	return fmt.Sprintf("%s/%d", ipAddr, prefixLen), nil
}

func (c Cidr) Value() (driver.Value, error) {
	return c.MarshalText()
}

func (c *Cidr) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return c.UnmarshalText(string(v))
	case string:
		return c.UnmarshalText(v)
	}

	return fmt.Errorf("cannot convert %T to Cidr", src)
}

func (c Cidr) Network() Ip {
	return Ip(c.IP.Mask(c.Mask))
}

func (c Cidr) FirstHost() Ip {
	ones, _ := c.Mask.Size()
	if ones < 30 {
		return c.Network().Add(1)
	} else {
		return c.Network()
	}
}

func (c Cidr) Contains(ip Ip) bool {
	return net.IP(ip).Mask(c.Mask).Equal(c.IP.Mask(c.Mask))
}

func NewCidr(cidr string) (Cidr, error) {
	var c Cidr
	err := c.UnmarshalText(cidr)
	return c, err
}

// MustNewCidr ignores errors returned by NewCidr.
// Should be used only for known-good inputs.
func MustNewCidr(cidr string) Cidr {
	result, _ := NewCidr(cidr)
	return result
}

type CidrArray []Cidr

func (a CidrArray) Value() (driver.Value, error) {
	var result = make([]string, 0, len(a))
	for _, elem := range a {
		result = append(result, elem.String())
	}
	return pq.StringArray(result).Value()
}

func (a *CidrArray) Scan(value interface{}) error {
	v := &pq.StringArray{}
	err := v.Scan(value)
	if err != nil {
		return err
	}
	for _, elem := range *v {
		c, err := NewCidr(elem)
		if err != nil {
			return err
		}
		*a = append(*a, c)
	}
	return nil
}
