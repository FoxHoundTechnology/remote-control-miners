package ant_miner_cgi

// cgi-bin/get_network_info.cgi

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

type NetworkInfoResponse struct {
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
}

// TODO: default values for response objects
func AntMinerCGIGetNetworkInfo(ipAddress string) (*NetworkInfoResponse, error) {
	return nil, nil
}
