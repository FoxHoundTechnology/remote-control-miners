package ant_miner_cgi

//  cgi-bin/get_miner_conf.cgi: Get miner configuration.

type rawGetMinerConfigResponse struct {
	// rawGetMinerConfResponse represents the configuration settings of a miner.
	APIAllow         string        `json:"api-allow"`
	APIGroups        string        `json:"api-groups"`
	APIListen        bool          `json:"api-listen"`
	APINetwork       bool          `json:"api-network"`
	BitmainCCDelay   string        `json:"bitmain-ccdelay"`
	BitmainFanCtrl   bool          `json:"bitmain-fan-ctrl"` // Whether Bitmain fan control is enabled
	BitmainFanPWM    string        `json:"bitmain-fan-pwm"`  // Fan PWM value, example: "100"
	BitmainFreq      string        `json:"bitmain-freq"`     // Frequency setting for Bitmain, example: "400"
	BitmainFreqLevel string        `json:"bitmain-freq-level"`
	BitmainPWTH      string        `json:"bitmain-pwth"`
	BitmainUseVIL    bool          `json:"bitmain-use-vil"`
	BitmainVoltage   string        `json:"bitmain-voltage"`
	BitmainWorkMode  string        `json:"bitmain-work-mode"` // Work mode for Bitmain, example: "0"
	Pools            []PoolSetting `json:"pools"`             // List of pool settings
}

type PoolSetting struct {
	Pass string `json:"pass"`
	URL  string `json:"url"`
	User string `json:"user"`
}

// NOTE: identical to set_miner_conf's payload struct
type GetMinerConfigResponse struct {
	BitmainFanCtrl bool   `json:"bitmain-fan-ctrl"` // Whether Bitmain fan control is enabled
	BitmainFanPWM  string `json:"bitmain-fan-pwm"`  // Fan PWM value, example: "100"
	FreqLevel      string `json:"freq-level"`       // Frequency level. NOTE: Different key name from rawGetMinerConfResponse
	MinerMode      string `json:"miner-mode"`       // Miner mode. NOTE: Different key name from rawGetMinerConfResponse
	Pools          []Pool `json:"pools"`            // List of pool settings
}

func AntMinerCGIGetMinerConfig(username, password, ipAddress string) (*GetMinerConfigResponse, error) {
	return nil, nil
}
