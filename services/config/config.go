package config

import (
	"errors"
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/utils/logger"
	"github.com/common-nighthawk/go-figure"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

const version = "1.0.4"

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
		return "", "", "", errors.New("no RFC 1918 private IPv4 address found")
	}

	publicIPv4, err = fetchPublicIP("https://api.ipify.org")
	if err != nil {
		return privateIPv4, "", "", err
	}

	publicIPv6, err = fetchPublicIP("https://api6.ipify.org")
	if err != nil {
		return privateIPv4, publicIPv4, "", err
	}

	return privateIPv4, publicIPv4, publicIPv6, nil
}

func fetchPublicIP(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
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
	} else {
		logger.Error("Unknown node detected, defaulting to Nulink")
		config.Node = "Nulink"
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
