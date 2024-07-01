package query

import (
	"time"
)

type Miner struct {
	MacAddress string `json:"MacAddress"`
	IPAddress  string `json:"IPAddress"`
}
type MinerStats struct {
	HashRate  float64 `json:"HashRate"`
	RateIdeal int     `json:"RateIdeal"`
	Uptime    int     `json:"Uptime"`
}
type MinerConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Firmware string `json:"Firmware"`
}

type PoolInfo struct {
	Url      string `json:"Url"`
	User     string `json:"User"`
	Pass     string `json:"Pass"`
	Status   string `json:"Status"`
	Accepted int    `json:"Accepted"`
	Rejected int    `json:"Rejected"`
	Stale    int    `json:"Stale"`
}

type PoolData struct {
	ID        int       `json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Pool      PoolInfo  `json:"Pool"`
	MinerID   int       `json:"MinerID"`
}

type MinerData struct {
	ID          int         `json:"ID"`
	CreatedAt   time.Time   `json:"CreatedAt"`
	UpdatedAt   time.Time   `json:"UpdatedAt"`
	Miner       Miner       `json:"Miner"`
	Stats       MinerStats  `json:"Stats"`
	Config      MinerConfig `json:"Config"`
	MinerType   int         `json:"MinerType"`
	ModelName   string      `json:"ModelName"`
	Mode        int         `json:"Mode"`
	Status      int         `json:"Status"`
	Pools       []PoolData  `json:"Pools"`
	Fan         []int       `json:"Fan"`
	Temperature []int       `json:"Temperature"`
	Log         interface{} `json:"Log"`
	FleetID     int         `json:"FleetID"`
}

type QueryMinerResponse struct {
	Data []MinerData `json:"data"`
}
