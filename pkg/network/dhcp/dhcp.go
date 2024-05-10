package network

import (
	"fmt"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	"github.com/ranjankuldeep/lab/pkg/constants"
)

// sets up the lease time for the allocated ip address.
var leaseDuration, _ = time.ParseDuration(constants.DHCP_INFINITE_LEASE)

type DHCPInterface struct {
	VMIPNet    *net.IPNet
	GatewayIP  *net.IP
	VMTAP      string
	Bridge     string
	Hostname   string
	MACFilter  string
	dnsServers []byte
}

// Justifying dhcp interface
func (i *DHCPInterface) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var respMsg dhcp.MessageType

	switch msgType {
	case dhcp.Discover:
		respMsg = dhcp.Offer
	case dhcp.Request:
		respMsg = dhcp.ACK
	}

	fmt.Printf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, msgType.String(), options, respMsg.String())
	if respMsg != 0 {
		requestingMAC := p.CHAddr().String()
		if requestingMAC == i.MACFilter {
			opts := dhcp.Options{
				dhcp.OptionSubnetMask:       []byte(i.VMIPNet.Mask),
				dhcp.OptionRouter:           []byte(*i.GatewayIP),
				dhcp.OptionDomainNameServer: i.dnsServers,
				dhcp.OptionHostName:         []byte(i.Hostname),
			}

			optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
			fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), i.GatewayIP.String(), i.VMIPNet.IP.String(), optSlice, requestingMAC)
			return dhcp.ReplyPacket(p, respMsg, *i.GatewayIP, i.VMIPNet.IP, leaseDuration, optSlice)
		}
	}

	return nil
}

func (i *DHCPInterface) StartBlockingServer() error {
	packetConn, err := conn.NewUDP4BoundListener(i.Bridge, ":67")
	if err != nil {
		return err
	}

	return dhcp.Serve(packetConn, i)
}

// Parse the DNS servers for the DHCP server
func (i *DHCPInterface) SetDNSServers(dns []string) {
	for _, server := range dns {
		i.dnsServers = append(i.dnsServers, []byte(net.ParseIP(server).To4())...)
	}
}
