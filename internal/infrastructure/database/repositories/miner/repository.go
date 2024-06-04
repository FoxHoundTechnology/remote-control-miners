package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TODO: WithContext method with the logic of timeout cancellation
// TODO: add preloads
// TODO: clause operation
// TODO: delete with cascade

type MinerRepository struct {
	db *gorm.DB
}

func NewMinerRepository(db *gorm.DB) *MinerRepository {
	return &MinerRepository{
		db: db,
	}
}

func (r *MinerRepository) GetByMacAddress(macAddress string) (*Miner, error) {
	miner := &Miner{}
	err := r.db.Preload("Pools").First(&miner, "mac_address = ?", macAddress).Error
	if err != nil {
		return nil, err
	}

	return miner, nil
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

	err = r.db.Save(&miner).Error
	if err != nil {
		return 0, err
	}

	return miner.ID, nil
}

func (r *MinerRepository) ListByFleetID(fleetId uint) ([]*Miner, error) {
	var miners []*Miner
	// TODO: Select statement
	// TODO: test a different way of defining the query with struct
	err := r.db.Where("fleet_id = ?", fleetId).Preload("Pools").Find(&miners).Error
	if err != nil {
		return nil, err
	}
	return miners, err
}

func (r *MinerRepository) ListByMacAddresses(macAddresses []string) ([]*Miner, error) {
	var miners []*Miner
	err := r.db.Where("mac_address IN (?)", macAddresses).Preload("Pools").Find(&miners).Error
	if err != nil {
		return nil, err
	}

	return miners, err
}

func (r *MinerRepository) List() ([]*Miner, error) {
	var miners []*Miner
	err := r.db.Preload("Pools").Find(&miners).Error
	if err != nil {
		return nil, err
	}

	return miners, err
}

func (r *MinerRepository) CreateMinersInBatch(miners []*Miner) error {

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, miner := range miners {

		// TODO! ideally insert on conclict operation
		if err := tx.Save(&miner).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error saving miner: %w", err)
		}

		// Save related Pools
		for _, pool := range miner.Pools {
			pool.MinerID = miner.ID // Ensure foreign key is set
			if err := tx.Save(&pool).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error saving pool: %w", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

// TODO!: separate createInBatches and updateInBatches
// TODO: R&D for association bulk update
func (r *MinerRepository) UpdateMinersInBatch(miners []*Miner) error {

	fmt.Println("BulkUpdateMinersWithPools,", len(miners))
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// TODO: Add mac address condition too
	// TODO: remove where clause and Model

	for _, miner := range miners {
		if err := tx.Omit("Pools").
			Save(miner).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating miners: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
