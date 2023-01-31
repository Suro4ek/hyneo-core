package sources

import (
	"context"
	"encoding/json"
	"hyneo/internal/antivpn"
	"io"
	"net/http"
	"time"
)

const VPNURL = "https://vpnapi.io/api/"

type VPNApi struct {
}

func NewVPNApi() antivpn.Source {
	return &VPNApi{}
}

func (V VPNApi) GetResult(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := http.NewRequestWithContext(ctx, "GET", VPNURL+ip, nil)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	var VPNAPI antivpn.VPNAPI
	err = json.Unmarshal(b, &VPNAPI)
	if err != nil {
		return false
	}
	if VPNAPI.Security.VPN || VPNAPI.Security.Proxy || VPNAPI.Security.Tor {
		return true
	}
	return false
}
