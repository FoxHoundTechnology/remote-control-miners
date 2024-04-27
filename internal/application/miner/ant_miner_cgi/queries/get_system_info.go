package queries

import (
	"encoding/json"
	"fmt"
	"foxhound/pkg/http_auth"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// cgi-bin/get_system_info.cgi: Get system information.
// TODO: default values for response objects
// TODO: modify the raw response object
type RawGetSystemInfoResponse struct {
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

type GetSystemInfoResponse struct {
	FirmwareType string `json:"firmware_type"`
	IPAddress    string `json:"ip_address"`
	MacAddress   string `json:"mac_address"`
	MinerType    string `json:"miner_type"` // used to determine which vendor the miner is from
}

func AntMinerCGIGetSystemInfo(username, password, ipAddress string) (*GetSystemInfoResponse, error) {
	t := http_auth.NewTransport(username, password)

	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_system_info.cgi", ipAddress), nil)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":      err,
			"newRequest": newRequest,
		}).Info("Error creating new request")

		return nil, err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err,
			"resp":  resp,
		}).Info("Error creating new request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  body,
		}).Info("Error reading response body")
		return nil, err
	}

	var rawGetSystemInfoResponse RawGetSystemInfoResponse
	err = json.Unmarshal(body, &rawGetSystemInfoResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Info("Error unmarshalling response body")
		return nil, err
	}

	return &GetSystemInfoResponse{
		FirmwareType: rawGetSystemInfoResponse.FirmwareType,
		IPAddress:    rawGetSystemInfoResponse.IPAddress,
		MacAddress:   rawGetSystemInfoResponse.MACAddress,
		MinerType:    rawGetSystemInfoResponse.MinerType,
	}, nil
}
