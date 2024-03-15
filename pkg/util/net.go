package util

import (
	"net"
	"strings"
)

const (
	ProtocolIPv4 = "IPv4"
	ProtocolIPv6 = "IPv6"
	ProtocolDual = "Dual"
)

func CheckProtocol(address string) string {
	ips := strings.Split(address, ",")
	if len(ips) == 2 {
		IP1 := net.ParseIP(strings.Split(ips[0], "/")[0])
		IP2 := net.ParseIP(strings.Split(ips[1], "/")[0])
		if IP1.To4() != nil && IP2.To4() == nil && IP2.To16() != nil {
			return ProtocolDual
		}
		if IP2.To4() != nil && IP1.To4() == nil && IP1.To16() != nil {
			return ProtocolDual
		}
		return ""
	}

	address = strings.Split(address, "/")[0]
	ip := net.ParseIP(address)
	if ip.To4() != nil {
		return ProtocolIPv4
	} else if ip.To16() != nil {
		return ProtocolIPv6
	}

	// cidr formal error
	return ""
}