package commands

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"net/http"

	"github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	"github.com/FoxHoundTechnology/remote-control-miners/pkg/http_auth"
)

// cgi-bin/set_miner_conf.cgi: Set miner configuration.

type SetMinerConfigPayload struct {
	BitmainFanCtrl bool          `json:"bitmain-fan-ctrl"` // Whether Bitmain fan control is enabled
	BitmainFanPWM  string        `json:"bitmain-fan-pwm"`  // Fan PWM value, example: "100"
	FreqLevel      string        `json:"freq-level"`       // Frequency level. NOTE: Different key name from rawGetMinerConfResponse
	MinerMode      string        `json:"miner-mode"`       // Miner mode. NOTE: Different key name from rawGetMinerConfResponse
	Pools          []domain.Pool `json:"pools"`            // List of pool settings
}

type rawSetMinerConfigResponse struct {
	Stats string `json:"stats"` // Status of the response, example: "success" TODO: enum
	Code  string `json:"code"`  // Code associated with the response, example: "M000"
	Msg   string `json:"msg"`   // Message detailing the response, example: "OK!"
}

type SetMinerConfigResponse struct {
	Stats string `json:"stats"`
}

func AntminerCGISetMinerConfig(username, password, ipAddress string, payload SetMinerConfigPayload) (*SetMinerConfigResponse, error) {
	marshalledPayload, err := json.Marshal(payload) // original payload is retrieved from get_miner_conf endpoint
	if err != nil {
		return nil, err
	}

	t := http_auth.NewTransport(username, password)
	newRequest, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/set_miner_conf.cgi", ipAddress), bytes.NewBuffer(marshalledPayload))
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

	var rawSetMinerConfigResponse rawSetMinerConfigResponse
	err = json.Unmarshal(body, &rawSetMinerConfigResponse)
	if err != nil {
		return nil, err
	}

	return &SetMinerConfigResponse{
		Stats: rawSetMinerConfigResponse.Stats,
	}, nil
}
