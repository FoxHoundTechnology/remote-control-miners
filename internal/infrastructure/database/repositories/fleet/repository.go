package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// TODO: WithContext method with the logic of timeout cancellation
// TODO: clause operation
// TODO: delete with cascade
type FleetRepository struct {
	db *gorm.DB
}

func NewFleetRepository(db *gorm.DB) *FleetRepository {
	return &FleetRepository{
		db: db,
	}
}

func (r *FleetRepository) Upsert(ctx context.Context, fleet *Fleet) (uint, error) {
	err := r.db.First(&fleet, "name = ?", fleet.Name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := r.db.Create(&fleet).Error
		if err != nil {
			return 0, err
		}
		return fleet.ID, nil
	}

	// Save is a combined function.
	// If save value does not contain its primary key,
	// it executes Create. Otherwise it executes Update (with all fields).
	err = r.db.Save(&fleet).Error
	if err != nil {
		return 0, err
	}

	return fleet.ID, nil
}

func (r *FleetRepository) List() ([]*Fleet, error) {
	var fleets []*Fleet
	err := r.db.Find(&fleets).Error
	if err != nil {
		return nil, err
	}
	return fleets, err
}
