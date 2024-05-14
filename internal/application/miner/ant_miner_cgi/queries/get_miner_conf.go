package queries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"

	"github.com/sirupsen/logrus"
)

//  cgi-bin/get_miner_conf.cgi: Get miner configuration.

type rawGetMinerConfigResponse struct {
	// rawGetMinerConfResponse represents the configuration settings of a miner.
	APIAllow         string `json:"api-allow"`
	APIGroups        string `json:"api-groups"`
	APIListen        bool   `json:"api-listen"`
	APINetwork       bool   `json:"api-network"`
	BitmainCCDelay   string `json:"bitmain-ccdelay"`
	BitmainFanCtrl   bool   `json:"bitmain-fan-ctrl"` // Whether Bitmain fan control is enabled
	BitmainFanPWM    string `json:"bitmain-fan-pwm"`  // Fan PWM value, example: "100"
	BitmainFreq      string `json:"bitmain-freq"`     // Frequency setting for Bitmain, example: "400"
	BitmainFreqLevel string `json:"bitmain-freq-level"`
	BitmainPWTH      string `json:"bitmain-pwth"`
	BitmainUseVIL    bool   `json:"bitmain-use-vil"`
	BitmainVoltage   string `json:"bitmain-voltage"`
	BitmainWorkMode  string `json:"bitmain-work-mode"` // Work mode for Bitmain, example: "0"
	// Pools            []domain.Pool `json:"pools"`             // NOTE: List of pool settings. Only settings, no stats included
}

// NOTE: identical to set_miner_conf's payload struct
type GetMinerConfigResponse struct {
	BitmainFanCtrl bool   `json:"bitmain-fan-ctrl"` // Whether Bitmain fan control is enabled
	BitmainFanPWM  string `json:"bitmain-fan-pwm"`  // Fan PWM value, example: "100"
	FreqLevel      string `json:"freq-level"`       // Frequency level. NOTE: Different key name from rawGetMinerConfResponse
	MinerMode      int    `json:"miner-mode"`       // Miner mode. NOTE: Different key name from rawGetMinerConfResponse
	// Pools          []domain.Pool `json:"pools"`            // NOTE: List of pool settings. Only settings, no stats included
}

func AntMinerCGIGetMinerConfig(username, password, ipAddress string) (*GetMinerConfigResponse, error) {

	t := http_auth.NewTransport(username, password)

	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/get_miner_conf.cgi", ipAddress), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"newRequest": newRequest,
		}).Debug("Error creating new request")
		return nil, err
	}

	resp, err := t.RoundTrip(newRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"resp":  resp,
		}).Debug("Error creating new request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  body,
		}).Debug("Error reading response body")
		return nil, err
	}

	var rawGetMinerConfigResponse rawGetMinerConfigResponse
	err = json.Unmarshal(body, &rawGetMinerConfigResponse)
	if err != nil {

		return nil, err
	}

	minerMode, err := strconv.Atoi(rawGetMinerConfigResponse.BitmainWorkMode)
	if err != nil {
		return nil, err
	}

	return &GetMinerConfigResponse{
		BitmainFanCtrl: rawGetMinerConfigResponse.BitmainFanCtrl,
		BitmainFanPWM:  rawGetMinerConfigResponse.BitmainFanPWM,
		FreqLevel:      rawGetMinerConfigResponse.BitmainFreqLevel,
		MinerMode:      minerMode,
	}, nil
}
