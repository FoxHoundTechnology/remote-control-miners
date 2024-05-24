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

// [ ]
// JOIN with miner config
/*
	Struct db.Find(&users, User{Age: 20})
	SELECT * FROM users WHERE age = 20;
*/
func (r *MinerRepository) ListByFleetID(fleetId uint) ([]*Miner, error) {
	var miners []*Miner
	// TODO: test preload
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

// TODO: R&D for association bulk update
func (r *MinerRepository) BulkUpdateMinersWithPools(miners []*Miner) error {

	fmt.Println("BulkUpdateMinersWithPools,", len(miners))

	temp := &gorm.Session{
		FullSaveAssociations: true,
		// SkipDefaultTransaction: true,
	}

	// Start a transaction
	tx := r.db.Session(temp).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// TODO: ideally with association,
	// 	     only one transaction should go here
	for _, miner := range miners {

		// Update the miner
		if err := tx.Save(&miner).Error; err != nil {
			tx.Rollback()
			return err
		}

		// db.Model(&person.Car).Update("Value", 9000)
		// db.Model(&person).Updates(Person{Name: "Jinzhu 2"})

		// db.Model(&person).Updates(map[string]interface{}{"Name": "Jinzhu 2"})

		// for _, pool := range miner.Pools {
		// 	// Check if the pool already exists
		// 	existingPool := Pool{}
		// 	err := tx.Where("miner_id = ? AND id = ?", miner.ID, pool.ID).First(&existingPool).Error

		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 		// If the pool doesn't exist, create a new one
		// 		pool.MinerID = miner.ID
		// 		if err := tx.Create(&pool).Error; err != nil {
		// 			tx.Rollback()
		// 			return err
		// 		}

		// 	} else if err == nil {
		// 		// If the pool exists, update it
		// 		existingPool.Pool = pool.Pool
		// 		if err := tx.Save(&existingPool).Error; err != nil {
		// 			tx.Rollback()
		// 			return err
		// 		}

		// 	} else {
		// 		tx.Rollback()
		// 		return err

		// 	}
		// }
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
