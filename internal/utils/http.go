package utils

import (
	"net"
	"strings"
)

func IsAllowedOrigin(origin string) bool {
	origin = strings.TrimPrefix(origin, "http://")
	origin = strings.TrimPrefix(origin, "https://")
	host := strings.Split(origin, ":")[0]

	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	ip := net.ParseIP(host)
	if ip != nil && isPrivateIP(ip) {
		return true
	}

	return false
}
