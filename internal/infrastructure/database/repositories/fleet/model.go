package repositories

import (
	repositories "foxhound/internal/infrastructure/database/repositories/miner"

	"gorm.io/gorm"
)

type Fleet struct {
	gorm.Model
	Name   string
	Miners []repositories.Miner `gorm:"foreignKey:FleetID;references:ID"` // one-to-many
}
