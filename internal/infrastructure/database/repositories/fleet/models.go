package repositories

import (
	"gorm.io/gorm"

	miner "foxhound/internal/infrastructure/database/repositories/miner"
)

type Fleet struct {
	gorm.Model
	Miners    []miner.Miner `gorm:"foreignKey:FleetID;references:ID"` // one-to-many
	ScannerID uint          // foreign key to Scanner, one-to-one
	// TODO: cascade
	// TODO: pool config
}
