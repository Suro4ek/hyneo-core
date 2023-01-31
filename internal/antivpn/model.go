package antivpn

type ProxyCheck struct {
	Proxy        string `json:"proxy"`
	Type         string `json:"type"`
	Provider     string `json:"provider"`
	Organisation string `json:"organisation"`
}

type VPNBlocker struct {
	Status            string `json:"status"`
	Package           string `json:"package"`
	RemainingRequests string `json:"remaining_requests"`
	Ipaddress         string `json:"ipaddress"`
	HostIp            string `json:"host-ip"`
	Organisation      string `json:"org"`
}

type Shodan struct {
	Tags map[ShodanTags]struct{} `json:"tags"`
}

type ShodanTags struct {
	VPN   *string `json:"vpn"`
	Proxy *string `json:"proxy"`
}

//type IpInfo struct {
//
//}

type VPNAPI struct {
	Security VPNApiSecurity `json:"security"`
}

type VPNApiSecurity struct {
	VPN   bool `json:"vpn"`
	Proxy bool `json:"proxy"`
	Tor   bool `json:"tor"`
}
