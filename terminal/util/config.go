package terminal

import (
	"encoding/json"
	"os"
)

type FleetConfig struct {
	Fleets []Fleet `json:"fleets"`
}

type Fleet struct {
	Name    string  `json:"name"`
	Scanner Scanner `json:"scanner"`
	Alert   Alert   `json:"alert"`
}

type Scanner struct {
	Name      string `json:"name"`
	StartIP   string `json:"start_ip"`
	EndIP     string `json:"end_ip"`
	Active    bool   `json:"active"`
	Location  string `json:"location"`
	Interval  int    `json:"interval"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	MinerType int    `json:"miner_type"`
}

type Alert struct {
	Name      string      `json:"name"`
	Action    int         `json:"action"`
	Active    bool        `json:"active"`
	Condition []Condition `json:"condition"`
}

type Condition struct {
	TriggerValue  int `json:"trigger_value"`
	MachineCount  int `json:"machine_count"`
	ThresholdType int `json:"threshold_type"`
	ConditionType int `json:"condition_type"`
	LayerType     int `json:"layer_type"`
}

func LoadFleetConfig(filename string) (*FleetConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config FleetConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
