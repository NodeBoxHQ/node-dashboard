package models

type Host struct {
	ID          int    `json:"id"`
	Hostname    string `json:"hostname"`
	Owner       string `json:"owner"`
	PrivateIPv4 string `json:"privateIpv4"`
	PrivateIPv6 string `json:"privateIpv6"`
	IPv4        string `json:"ipv4"`
	IPv6        string `json:"ipv6"`
	Node        string `json:"node"`
}
