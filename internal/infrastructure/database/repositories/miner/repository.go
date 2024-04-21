package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// TODO: WithContext method with the logic of timeout cancellation
// TODO: add preloads
// TODO: clause operation
// TODO: delete with cascade

// receiver functions go here
type MinerRepository struct {
	db *gorm.DB
}

func NewMinerRepository(db *gorm.DB) *MinerRepository {
	return &MinerRepository{
		db: db,
	}
}

func (r *MinerRepository) Upsert(ctx context.Context, miner *Miner) (uint, error) {
	err := r.db.First(&miner, "mac_address = ?", miner.Miner.MacAddress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := r.db.Create(&miner).Error
		if err != nil {
			return 0, err
		}
		return miner.ID, nil
	}

	// Save is a combined function.
	// If save value does not contain primary key,
	// it executes Create. Otherwise it executes Update (with all fields).
	err = r.db.Save(&miner).Error
	if err != nil {
		return 0, err
	}

	return miner.ID, nil
}

// [ ]
// NOT TESTED
// JOIN with miner config
/*
	Struct db.Find(&users, User{Age: 20})
	SELECT * FROM users WHERE age = 20;
*/
func (r *MinerRepository) ListByFleetID(ctx context.Context, miner *Miner) ([]*Miner, error) {
	var miners []*Miner
	// TODO: test preload
	// TODO: test different way of defining the query with struct
	err := r.db.Find(&miners, "fleet_id = ?", miners).Error
	if err != nil {
		return nil, err
	}
	return miners, err
}

func (r *MinerRepository) List() ([]*Miner, error) {
	var miners []*Miner
	err := r.db.Find(&miners).Error
	if err != nil {
		return nil, err
	}
	return miners, err
}
