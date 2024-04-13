package ant_miner_cgi

// cgi-bin/get_system_info.cgi: Get system information.

type rawGetSystemInfoResponse struct {
	DNSServers              string `json:"dnsservers"`
	FirmwareType            string `json:"firmware_type"`
	Gateway                 string `json:"gateway"`
	Hostname                string `json:"hostname"`
	IPAddress               string `json:"ipaddress"`
	MACAddress              string `json:"macaddr"`
	MinerType               string `json:"minertype"`
	NetDevice               string `json:"netdevice"`
	NetMask                 string `json:"netmask"`
	NetType                 string `json:"nettype"`
	SerialNumber            string `json:"serinum"`
	SystemFilesystemVersion string `json:"system_filesystem_version"`
	SystemKernelVersion     string `json:"system_kernel_version"`
	SystemMode              string `json:"system_mode"`
}

type SystemInfoResponse struct {
	FirmwareType string `json:"firmware_type"`
	IPAddress    string `json:"ip_address"`
	MacAddress   string `json:"mac_address"`
	MinerType    string `json:"miner_type"` // used to determine which vendor the miner is from
}

func AntMinerCGIGetSystemInfo(ipAddress string) (*SystemInfoResponse, error) {

	return nil, nil
}
