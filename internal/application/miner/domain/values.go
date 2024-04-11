package miner

import (
	"time"
)

type Mode int

const (
	NormalMode   Mode = iota // 0
	SleepMode                // 1
	LowPowerMode             // 2
)

type Status int

const (
	Online   Status = iota // 0
	Offline                // 1
	Disabled               // 2
	Error                  // 3
	Warning                // 4
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

type Pools struct {
	Pools []Pool
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
	TemperatureSensors []TemperatureSensor
}

type TemperatureSensor struct {
	Name string
	Temp float64
}

type Fan struct {
	FanSensors []FanSensor
}

type FanSensor struct {
	Name  string
	Speed int
}