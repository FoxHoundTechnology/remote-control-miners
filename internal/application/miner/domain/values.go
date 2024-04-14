package domain

import (
	"database/sql/driver"
	"errors"
	"time"
)

// OPTION: segregate the domain values with value/scan methods into their respective files as it grows
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
	HashRate    float64   // in GH/s
	RateIdeal   float64   // in GH/s
	Uptime      int       // in seconds
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

type TemperatureSensor struct {
	Name    string
	TempPcb []int
}

type FanSensor struct {
	Name  string
	Speed int
}

type Log struct {
	Name      string
	EventType EventType
	Timestamp time.Time
}

type EventType int

const (
	Operational  = iota // 0
	SystemIssue         // 1
	UserActivity        // 2
)

// Value/Scan methods for Mode
func (m Mode) Value() (driver.Value, error) {
	// Convert the Mode to its string representation
	switch m {
	case NormalMode:
		return "normal", nil
	case SleepMode:
		return "sleep", nil
	case LowPowerMode:
		return "lowpower", nil
	default:
		return nil, errors.New("invalid Mode value")
	}
}

func (m *Mode) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("conversion error: failed to convert the interface to string")
	}
	switch s {
	case "normal":
		*m = NormalMode
	case "sleep":
		*m = SleepMode
	case "lowpower":
		*m = LowPowerMode
	default:
		return errors.New("invalid Mode string")
	}
	return nil
}

// Value converts the Status to a string that can be stored in the database.
func (s Status) Value() (driver.Value, error) {
	switch s {
	case Online:
		return "online", nil
	case Offline:
		return "offline", nil
	case Disabled:
		return "disabled", nil
	case Error:
		return "error", nil
	case Warning:
		return "warning", nil
	default:
		return nil, errors.New("invalid Status value")
	}
}

// Scan reads a string from the database and converts it to a Status.
func (s *Status) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("conversion error: failed to convert the interface to string")
	}
	switch str {
	case "online":
		*s = Online
	case "offline":
		*s = Offline
	case "disabled":
		*s = Disabled
	case "error":
		*s = Error
	case "warning":
		*s = Warning
	default:
		return errors.New("invalid Status string")
	}
	return nil
}
