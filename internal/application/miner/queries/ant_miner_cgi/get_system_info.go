package queries

// cgi-bin/get_system_info.cgi: Get system information.

type getSystemInfoResponse struct {
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

type SystemInfoResponse struct{}
