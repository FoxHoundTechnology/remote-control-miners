package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/miner/domain"
)

// TODO: decorator pattern for miner fleet
type Miner struct {
	gorm.Model
	Miner  domain.Miner  `gorm:"embedded;"`
	Stats  domain.Stats  `gorm:"embedded;"`
	Config domain.Config `gorm:"embedded;"`
	Status domain.Status `gorm:"type:miner_status;"`

	Pools []Pool

	FleetID uint // foreign key to Fleet, one-to-many
}

type Pool struct {
	gorm.Model
	Pool    domain.Pool `gorm:"embedded;"`
	MinerID uint        // foreign key to Miner, one-to-many
}
