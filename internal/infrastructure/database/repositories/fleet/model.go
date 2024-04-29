package repositories

import (
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	"gorm.io/gorm"
)

type Fleet struct {
	gorm.Model
	Name    string `gorm:"unique"`
	Miners  []miner_repo.Miner
	Scanner scanner_repo.Scanner
}
