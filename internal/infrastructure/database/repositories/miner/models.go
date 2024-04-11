package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/miner/domain"
)

// TODO: cascade/association
// TODO: pool config
type Fleet struct {
	gorm.Model
	Miners    []Miner `gorm:"foreignKey:FleetID;references:ID"` // one-to-many
	ScannerID uint    // foreign key to Scanner, one-to-one
}

type Miner struct {
	gorm.Model
	Miner  domain.Miner  `gorm:"embedded;"`
	Stats  domain.Stats  `gorm:"embedded;"`
	Config domain.Config `gorm:"embedded;"`

	Mode   domain.Mode   `gorm:"comment: Mode: 0=Normal, 1=Sleep, 2=LowPower;"`
	Status domain.Status `gorm:"comment: Status: 0=Online, 1=Offline, 2=Disabled, 3=Error, 4=Warning;"`

	Pools       []Pool      `gorm:"foreignKey:MinerID;"`
	Temperature Temperature `gorm:"foreignKey:MinerID;"`
	Fan         Fan         `gorm:"foreignKey:MinerID;"`

	FleetID uint // foreign key to Fleet, one-to-many
}

type Pool struct {
	gorm.Model
	Pool    domain.Pool `gorm:"embedded;"`
	MinerID uint        // foreign key to Miner, one-to-many
}

// TODO: composite pattern for Miner structure
type Temperature struct {
	gorm.Model
	Sensors []TemperatureSensor `gorm:"foreignKey:TemperatureID;"`
}
type TemperatureSensor struct {
	gorm.Model
	Sensor domain.TemperatureSensor `gorm:"embedded;"`
}

type Fan struct {
	gorm.Model
	Fans []FanSensor `gorm:"foreignKey:FanID;"`
}
type FanSensor struct {
	gorm.Model
	Sensor domain.FanSensor `gorm:"embedded;"`
}
