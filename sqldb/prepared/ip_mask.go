// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package prepared

import (
	"database/sql/driver"
	"fmt"
	"github.com/lib/pq"
	"net"
)

type IpMask net.IPNet

func (c IpMask) String() string {
	v, _ := c.MarshalText()
	return v
}

func (c *IpMask) UnmarshalText(data string) error {
	ip, ipNet, err := net.ParseCIDR(data)
	if err != nil {
		return err
	}
	c.IP = ip
	c.Mask = ipNet.Mask
	return nil
}

func (c IpMask) MarshalText() (string, error) {
	ipAddr := c.IP.String()
	prefixLen, _ := c.Mask.Size()
	return fmt.Sprintf("%s/%d", ipAddr, prefixLen), nil
}

func (c IpMask) Value() (driver.Value, error) {
	return c.MarshalText()
}

func (c *IpMask) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return c.UnmarshalText(string(v))
	case string:
		return c.UnmarshalText(v)
	}

	return fmt.Errorf("cannot convert %T to IpMask", src)
}

func (c IpMask) Network() Ip {
	return Ip(c.IP.Mask(c.Mask))
}

func (c IpMask) FirstHost() Ip {
	ones, _ := c.Mask.Size()
	if ones < 30 {
		return c.Network().Add(1)
	} else {
		return c.Network()
	}
}

func (c IpMask) Contains(ip Ip) bool {
	return net.IP(ip).Mask(c.Mask).Equal(c.IP.Mask(c.Mask))
}

func NewIpMask(ipMask string) (IpMask, error) {
	var c IpMask
	err := c.UnmarshalText(ipMask)
	return c, err
}

// MustNewIpMask ignores errors returned by NewIpMask.
// Should be used only for known-good inputs.
func MustNewIpMask(ipMask string) IpMask {
	result, _ := NewIpMask(ipMask)
	return result
}

type IpMaskArray []IpMask

func (a IpMaskArray) Value() (driver.Value, error) {
	var result = make([]string, 0, len(a))
	for _, elem := range a {
		result = append(result, elem.String())
	}
	return pq.StringArray(result).Value()
}

func (a *IpMaskArray) Scan(value interface{}) error {
	v := &pq.StringArray{}
	err := v.Scan(value)
	if err != nil {
		return err
	}
	for _, elem := range *v {
		c, err := NewIpMask(elem)
		if err != nil {
			return err
		}
		*a = append(*a, c)
	}
	return nil
}
