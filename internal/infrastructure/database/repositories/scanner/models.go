package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/scanner/domain"

	fleet "foxhound/internal/infrastructure/database/repositories/fleet"
)

type Scanner struct {
	gorm.Model
	Scanner domain.Scanner `gorm:"embedded;"`
	Status  domain.Status  `gorm:"embedded;"`
	Config  domain.Config  `gorm:"embedded;"`

	MinerType string `gorm:"type:miner_type"`
	Owner     string
	Location  string

	Fleet fleet.Fleet // foreign key to Fleet, one-to-one
}
