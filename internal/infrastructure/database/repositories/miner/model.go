package repositories

import (
	"foxhound/internal/application/miner/domain"
	scanner_domain "foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Miner     domain.Miner  `gorm:"embedded"`
	Stats     domain.Stats  `gorm:"embedded"`
	Config    domain.Config `gorm:"embedded"`
	MinerType scanner_domain.MinerType

	Mode   domain.Mode   `gorm:"comment: Mode: 0=Normal, 1=Sleep, 2=LowPower"`
	Status domain.Status `gorm:"comment: Status: 0=Online, 1=Offline, 2=Disabled, 3=Error, 4=Warning"`

	Pools       []Pool              `gorm:"foreignKey:MinerID;references:ID"`
	Temperature []TemperatureSensor `gorm:"foreignKey:MinerID;references:ID; comment: A collection of highest temperatures"`
	Fan         []FanSensor         `gorm:"foreignKey:MinerID;references:ID"`
	Log         []MinerLog          `gorm:"foreignKey:MinerID;references:ID"`

	FleetID uint `gorm:"foreignKey:FleetID;"`
}

type Pool struct {
	gorm.Model
	Pool    domain.Pool `gorm:"embedded;"`
	MinerID uint        `gorm:"foreignKey:MinerID;"`
}

type TemperatureSensor struct {
	gorm.Model
	Name        string
	Temperature int
	MinerID     uint `gorm:"foreignKey:MinerID;"`
}

type FanSensor struct {
	gorm.Model
	Sensor  domain.FanSensor `gorm:"embedded;"`
	MinerID uint             `gorm:"foreignKey:MinerID;"`
}

type MinerLog struct {
	gorm.Model
	Log       domain.Log       `gorm:"embedded;"`
	EventType domain.EventType `gorm:"comment: EventType: 0=Operational, 1=SystemIssue, 2=UserActivity"`
	MinerID   uint             `gorm:"foreignKey:MinerID;"`
}
