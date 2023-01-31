package sources

import "hyneo/internal/antivpn"

type proxyCheck struct {
}

const ProxyCheckURL = "https://vpnapi.i1o/api/"

func NewProxyCheck() antivpn.Source {
	return &proxyCheck{}
}

func (p proxyCheck) GetResult(ip string) bool {

}
