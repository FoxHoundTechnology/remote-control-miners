package repositories

import (
	"gorm.io/gorm"

	domain "foxhound/internal/application/scanner/domain"

	miner "foxhound/internal/infrastructure/database/repositories/miner"
)

type Scanner struct {
	gorm.Model
	Scanner domain.Scanner `gorm:"embedded;"`
	Status  domain.Status  `gorm:"embedded;"`
	Config  domain.Config  `gorm:"embedded;"`

	MinerType string `gorm:"type:miner_type"`
	Owner     string
	Location  string

	Fleet miner.Fleet // foreign key to Fleet, one-to-one
}
