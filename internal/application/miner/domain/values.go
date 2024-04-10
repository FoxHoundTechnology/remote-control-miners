package miner

import (
	"time"
)

type Mode string

const (
	NormalMode   Mode = "normal"
	SleepMode    Mode = "sleep"
	LowPowerMode Mode = "low_power"
)

type Status string

const (
	Online   Status = "online"
	Offline  Status = "offline"
	Disabled Status = "disabled"
	Error    Status = "error"
	Warning  Status = "warning"
)

type Config struct {
	Username string
	Password string
	Firmware string
}

type Stats struct {
	HashRate    float64 // in GH/s
	RateIdeal   float64 // in GH/s
	Uptime      int     // in seconds
	Location    string
	FanSpeed    int
	LastUpdated time.Time // UTC timestamp
}

type Pool struct {
	Url  string
	User string
	Pass string
}

type Fleet struct {
	Miners []Miner
}

type Temperature struct {
	Temperature []TemperatureSensor
}

type TemperatureSensor struct {
	Name string
	Temp float64
}

type Fan struct {
	Fan []FanSensor
}

type FanSensor struct {
	Name  string
	Speed int
}
