// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"github.com/jackpal/gateway"
	"github.com/pkg/errors"
	"net"
	"os"
	"strings"
)

type NetworkProvider struct {
	Describer
	SilentNotifier
}

func (p *NetworkProvider) Load(ctx context.Context) (results ProviderEntries, err error) {
	var outboundAddress net.IP
	var outboundInterface string
	var outboundSubnet net.IPNet
	if outboundAddress, outboundInterface, outboundSubnet, err = p.loadOutboundAddress(); err != nil {
		return
	}

	// Load the default gateway
	var outboundGateway net.IP
	if outboundGateway, err = gateway.DiscoverGateway(); err != nil {
		switch err.Error() {
		case "no gateway found":
			logger.WithContext(ctx).Warn(err)
		default:
			return nil, err
		}
	}

	// Load the hostname
	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		return nil, err
	}

	return ProviderEntries{
		NewEntry(p, "network.outbound.address", outboundAddress.String()),
		NewEntry(p, "network.outbound.interface", outboundInterface),
		NewEntry(p, "network.outbound.subnet", outboundSubnet.String()),
		NewEntry(p, "network.outbound.gateway", outboundGateway.String()),
		NewEntry(p, "network.hostname", hostname),
	}, nil
}

func (p *NetworkProvider) loadOutboundAddress() (addr net.IP, ifName string, subnet net.IPNet, err error) {
	// Create the underlying transport to an address that is
	// guaranteed to be _more_ public.   This does not create
	// a connection (as it is UDP), but instead merely
	// calculates the local interface from which such a connection
	// would be made.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	_ = conn.Close()
	addr = conn.LocalAddr().(*net.UDPAddr).IP

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		var addrs []net.Addr
		addrs, err = iface.Addrs()
		if err != nil {
			return
		}

		for _, ifAddr := range addrs {
			ifAddrString := ifAddr.String()
			slashOffset := strings.Index(ifAddrString, "/")
			if slashOffset > 0 {
				ifAddrString = ifAddrString[:slashOffset]
			}

			if addr.String() != ifAddrString {
				continue
			}

			ifName = iface.Name

			switch v := ifAddr.(type) {
			case *net.IPNet:
				subnet = net.IPNet{
					IP:   v.IP.Mask(v.Mask),
					Mask: v.Mask,
				}

			case *net.IPAddr:
				subnet = net.IPNet{
					IP:   v.IP.Mask(v.IP.DefaultMask()),
					Mask: v.IP.DefaultMask(),
				}
			}

			return
		}
	}

	err = errors.Errorf("Could not determine interface for address %q", addr.String())
	return
}

func NewNetworkProvider() *NetworkProvider {
	return &NetworkProvider{
		Describer: Named{
			name: "Network",
		},
	}
}
