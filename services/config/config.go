package config

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/utils/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/pion/stun"
	"net"
	"os"
	"strings"
)

const version = "1.1.6"

type Config struct {
	Node                    string `json:"node"`
	PrivateIPv4             string `json:"private_ip"`
	IPv4                    string `json:"ip"`
	IPv6                    string `json:"ip6"`
	Owner                   string `json:"owner"`
	Port                    int    `json:"port"`
	NodeboxDashboardVersion string `json:"nodebox_dashboard_version"`
}

var loadedConfig *Config

func isRFC1918Private(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}
	return false
}

func getIPAddresses() (privateIPv4, publicIPv4, publicIPv6 string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || !isRFC1918Private(ip) {
				continue
			}
			ip = ip.To4()
			if ip != nil {
				privateIPv4 = ip.String()
				break
			}
		}
		if privateIPv4 != "" {
			break
		}
	}
	if privateIPv4 == "" {
		privateIPv4 = "127.0.0.1"
	}

	publicIPv4, publicIPv6, err = getPublicIPsViaSTUN()
	if err != nil {
		return privateIPv4, "", "", err
	}

	return privateIPv4, publicIPv4, publicIPv6, nil
}

func getPublicIPsViaSTUN() (publicIPv4, publicIPv6 string, err error) {
	stunServers := []string{
		"[2001:4860:4864:5:8000::1]:19302",
		"74.125.250.129:19302",
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

func LoadConfig() (*Config, error) {
	var config Config
	hostname, err := os.Hostname()

	logger.Info("Hostname: ", hostname, " detected")

	if err != nil {
		return nil, err
	}

	var badHostname bool

	if strings.Contains(hostname, "-linea") {
		logger.Info("Linea node detected")
		config.Node = "Linea"
	} else if strings.Contains(hostname, "-dusk") {
		logger.Info("Dusk node detected")
		config.Node = "Dusk"
	} else if strings.Contains(hostname, "-nulink") {
		logger.Info("Nulink node detected")
		config.Node = "Nulink"
	} else if strings.Contains(hostname, "-babylon") {
		logger.Info("Babylon node detected")
		config.Node = "Babylon"
	} else if strings.Contains(hostname, "-xcally") || strings.Contains(hostname, "-xally") {
		logger.Info("Xally node detected")
		config.Node = "Xally"
	} else if strings.Contains(hostname, "-juneo") || strings.Contains(hostname, "-june") || strings.Contains(hostname, "-juneogo") {
		logger.Info("Juneo node detected")
		config.Node = "Juneo"
	} else {
		logger.Error("Unknown node detected, defaulting to Juneo")
		config.Node = "Juneo"
		badHostname = true
	}
	
	privateIpv4, ipv4, ipv6, err := getIPAddresses()

	if err != nil {
		return nil, err
	}

	config.PrivateIPv4 = privateIpv4
	config.IPv4 = ipv4
	config.IPv6 = ipv6

	if !badHostname {
		config.Owner = hostname[:len(hostname)-len(config.Node)-1]
	} else {
		config.Owner = "Unknown"
	}

	config.Port = 3000
	config.NodeboxDashboardVersion = version

	loadedConfig = &config

	return &config, nil
}

func GetNodeType() string {
	return loadedConfig.Node
}

func ShowAsciiArt() {
	myFigure := figure.NewFigure("Nodebox", "doom", true)
	myFigure.Print()
	fmt.Println("\n\t\t\t\t\tVersion: ", version)
	fmt.Println("\n")
}
