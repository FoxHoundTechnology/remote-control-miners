package repositories

import (
	miner_domain "foxhound/internal/application/miner/domain"
	scanner_domain "foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Miner     miner_domain.Miner  `gorm:"embedded"`
	Stats     miner_domain.Stats  `gorm:"embedded"`
	Config    miner_domain.Config `gorm:"embedded"`
	MinerType scanner_domain.MinerType

	Mode   miner_domain.Mode   `gorm:"comment: Mode: 0=Normal, 1=Sleep, 2=LowPower"`
	Status miner_domain.Status `gorm:"comment: Status: 0=Online, 1=Offline, 2=Disabled, 3=Error, 4=Warning"`

	Pools       []Pool
	Temperature []TemperatureSensor `gorm:"comment: A collection of highest temperatures"`
	Fan         []FanSensor
	Log         []MinerLog

	FleetID uint
}

type Pool struct {
	gorm.Model
	Pool    miner_domain.Pool `gorm:"embedded;"`
	MinerID uint
}

type TemperatureSensor struct {
	gorm.Model
	Name        string
	Temperature int
	MinerID     uint
}

type FanSensor struct {
	gorm.Model
	Sensor  miner_domain.FanSensor `gorm:"embedded;"`
	MinerID uint
}

type MinerLog struct {
	gorm.Model
	Log       miner_domain.Log       `gorm:"embedded;"`
	EventType miner_domain.EventType `gorm:"comment: EventType: 0=Operational, 1=SystemIssue, 2=UserActivity"`
	MinerID   uint
}
