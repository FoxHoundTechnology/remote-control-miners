package domain

import (
	"time"
)

// TODO: fix pool settings and pool stats
type Mode int

const (
	NormalMode   Mode = iota // 0
	SleepMode                // 1
	LowPowerMode             // 2
)

type Command int

const (
	Normal Command = iota // 0
	Sleep
	LowPower
	Reboot
	// ... other miner commands go here
)

type Status int

const (
	Online                Status = iota // 0
	Offline                             // 1
	Disabled                            // 2
	HashrateError                       // 3
	TemperatureError                    // 4
	FanSpeedError                       // 5
	MissingHashboardError               // 6
	PoolShareError                      // 7
)

type Config struct {
	Username string
	Password string
	Firmware string
}

type Stats struct {
	HashRate  float64 // in GH/s
	RateIdeal float64 // in GH/s
	Uptime    int     // in seconds
}

type Pool struct {
	Url      string
	User     string
	Pass     string
	Status   string // NOTE: string for now
	Accepted int
	Rejected int
	Stale    int
}

type TemperatureSensor struct {
	Name       string
	PcbSensors []int
}

type FanSensor struct {
	Name  string
	Speed int
}

type Log struct {
	Description string
	Timestamp   time.Time
}

type EventType int

const (
	Operational  = iota // 0
	SystemIssue         // 1
	UserActivity        // 2
)
