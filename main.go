package main

import (
	"net"

	network "github.com/ranjankuldeep/lab/pkg/network/dhcp"
)

func main() {
	// Create a new instance of DHCPInterface
	dhcpInterface := &network.DHCPInterface{
		VMIPNet:   &net.IPNet{IP: net.IP{172, 16, 1, 1}, Mask: net.IPv4Mask(255, 255, 255, 0)},
		GatewayIP: &net.IP{172, 16, 1, 254},
		VMTAP:     "tap0",
		Bridge:    "br0",
		Hostname:  "example-host",
		MACFilter: "00:11:22:33:44:55", // Example MAC address to filter
	}

	// Set DNS servers for the DHCP server
	dhcpInterface.SetDNSServers([]string{"8.8.8.8", "8.8.4.4"}) // Example DNS servers

	// Start the DHCP server in a blocking manner
	if err := dhcpInterface.StartBlockingServer(); err != nil {
		panic(err)
	}
}
