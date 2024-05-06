package repositories

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
	scanner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/scanner/domain"

	"gorm.io/gorm"
)

// TODO: cascade for log

type Fan []int
type Temperature []int

type Miner struct {
	gorm.Model
	Miner     miner_domain.Miner  `gorm:"embedded"`
	Stats     miner_domain.Stats  `gorm:"embedded"`
	Config    miner_domain.Config `gorm:"embedded"`
	MinerType scanner_domain.MinerType

	Mode   miner_domain.Mode   `gorm:"comment: Mode: 0=Normal, 1=Sleep, 2=LowPower"`
	Status miner_domain.Status `gorm:"comment: Status: 0=Online, 1=Offline, 2=Disabled, 3=Error, 4=Warning"`

	Pools       []Pool      `gorm:"onDelete:CASCADE; onUpdate:CASCADE"`
	Fan         Fan         `gorm:"type:VARCHAR(255)"`
	Temperature Temperature `gorm:"type:VARCHAR(255)"`
	Log         []MinerLog

	FleetID uint
}

type Pool struct {
	gorm.Model
	Pool    miner_domain.Pool `gorm:"embedded;"`
	MinerID uint
}

type MinerLog struct {
	gorm.Model
	Log       miner_domain.Log       `gorm:"embedded;"`
	EventType miner_domain.EventType `gorm:"comment: EventType: 0=Operational, 1=SystemIssue, 2=UserActivity"`
	MinerID   uint
}

// ============== Scan/Values for Fan and Temp ==============
func (f *Fan) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, f)
	case string:
		return json.Unmarshal([]byte(v), f)
	default:
		return errors.New("failed to unmarshal value into bytes")
	}
}

func (f Fan) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (t *Temperature) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return errors.New("failed to unmarshal value into bytes")
	}
}
func (t Temperature) Value() (driver.Value, error) {
	return json.Marshal(t)
}
