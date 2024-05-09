package repositories

import "time"

type MinerTimeSeries struct {
	MacAddress string `json:"mac_address"`
	HashRate   int    `json:"hashrate"`
	TempSensor []int  `json:"temp_sensor"` // Assuming a maximum of 10 temperature sensors
	FanSensor  []int  `json:"fan_sensor"`  // Assuming a maximum of 10 fan sensors
}

type MinerTimeSeriesResponse struct {
	Record     []MinerTimeSeries `json:"miner_time_series_record"`
	TimeStamps []time.Time       `json:"timestamp"`
}

type PoolTimeSeries struct {
	MacAddress string `json:"mac_address"`
	Accepted   int    `json:"accepted"`
	Rejected   int    `json:"rejected"`
	Stale      int    `json:"stale"`
}

type PoolTimeSeriesResponse struct {
	Record     []PoolTimeSeries `json:"pool_time_series_record"`
	TimeStamps []time.Time      `json:"timestamps"`
}
