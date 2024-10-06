package utils

import (
	"fmt"
	"net"

	"gortc.io/stun"
)

func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("fd00::"), Mask: net.CIDRMask(8, 128)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func GetPublicIPs() (publicIPv4, publicIPv6 string, err error) {
	stunServers := []string{
		"[2001:4860:4864:5:8000::1]:19302",
		"74.125.250.129:19302",
		"stun.difuse.io:3478",
	}

	for _, server := range stunServers {
		client, err := stun.Dial("udp", server)
		if err != nil {
			continue
		}
		message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
		var response *stun.Message
		err = client.Do(message, func(event stun.Event) {
			if event.Error != nil {
				err = event.Error
				return
			}
			response = event.Message
		})
		if err != nil {
			continue
		}

		var addr stun.XORMappedAddress
		if err := addr.GetFrom(response); err != nil {
			continue
		}

		if addr.IP.To4() != nil {
			publicIPv4 = addr.IP.String()
		} else {
			publicIPv6 = addr.IP.String()
		}

		if publicIPv4 != "" {
			break
		}
	}

	if publicIPv4 == "" {
		err = fmt.Errorf("No valid IPv4 address retrieved from any STUN server")
	}

	return publicIPv4, publicIPv6, err
}

func GetPrivateIPs() (privateIPv4, privateIPv6 string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if ip := v.IP.To4(); ip != nil && isPrivateIP(ip) {
					privateIPv4 = ip.String()
				} else if ip := v.IP.To16(); ip != nil && isPrivateIP(ip) {
					privateIPv6 = ip.String()
				}
			}
		}
	}

	return privateIPv4, privateIPv6
}
