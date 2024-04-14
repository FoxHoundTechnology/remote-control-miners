package queries

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

func AntMinerCGIGetNetworkInfo(ipAddress string) (*GetNetworkInfoResponse, error) {
	return nil, nil
}
