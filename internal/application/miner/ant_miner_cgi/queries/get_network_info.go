package queries

import (
	"encoding/json"
	"fmt"

	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	"io"
	"net/http"
)

// cgi-bin/get_network_info.cgi
// TODO: default values for response objects
type rawGetNetworkInfoResponse struct {
	NetType        string `json:"nettype"`
	NetDevice      string `json:"netdevice"`
	MacAddress     string `json:"macaddr"`
	IPAddress      string `json:"ipaddress"`
	NetMask        string `json:"netmask"`
	ConfNetType    string `json:"conf_nettype"`
	ConfHostname   string `json:"conf_hostname"`
	ConfIPAddress  string `json:"conf_ipaddress"`
	ConfNetMask    string `json:"conf_netmask"`
	ConfGateway    string `json:"conf_gateway"`
	ConfDNSServers string `json:"conf_dnsservers"`
}

type GetNetworkInfoResponse struct {
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
}

func AntMinerCGIGetNetworkInfo(username, password, ipAddress string) (*GetNetworkInfoResponse, error) {

	t := http_auth.NewTransport(username, password)

	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_network_info.cgi", ipAddress), nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawGetNetworkInfoResponse rawGetNetworkInfoResponse
	err = json.Unmarshal(body, &rawGetNetworkInfoResponse)
	if err != nil {
		return nil, err
	}

	return &GetNetworkInfoResponse{
		MacAddress: rawGetNetworkInfoResponse.MacAddress,
		IPAddress:  rawGetNetworkInfoResponse.IPAddress,
	}, nil
}
