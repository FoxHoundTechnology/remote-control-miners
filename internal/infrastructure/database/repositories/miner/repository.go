package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// TODO!: separate createInBatches and updateInBatches
// TODO: R&D for association bulk update
func (r *MinerRepository) BulkUpdateMinersWithPools(miners []*Miner) error {

	fmt.Println("BulkUpdateMinersWithPools,", len(miners))

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Construct the bulk upsert query
	var valueStrings []string
	var valueArgs []interface{}

	for _, miner := range miners {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, miner.Miner.MacAddress, miner.Miner.IPAddress, miner.Stats.HashRate, miner.Stats.RateIdeal, miner.Stats.Uptime, miner.Config.Username, miner.Config.Password, miner.Config.Firmware, miner.MinerType, miner.ModelName, miner.Mode, miner.Status, miner.Fan, miner.Temperature, miner.FleetID)
	}

	query := fmt.Sprintf(`
		INSERT INTO miners (mac_address, ip_address, hash_rate, rate_ideal, uptime, username, password, firmware, miner_type,
		model_name, mode, status, fan, temperature, fleet_id)
		VALUES %s
		ON CONFLICT (mac_address) DO UPDATE 
		SET 
		ip_address = EXCLUDED.ip_address,
		hash_rate = EXCLUDED.hash_rate,
		rate_ideal = EXCLUDED.rate_ideal,
		uptime = EXCLUDED.uptime,
		username = EXCLUDED.username,
		password = EXCLUDED.password,
		firmware = EXCLUDED.firmware,
		miner_type = EXCLUDED.miner_type,
		model_name = EXCLUDED.model_name,
		mode = EXCLUDED.mode,
		status = EXCLUDED.status,
		fan = EXCLUDED.fan,
		temperature = EXCLUDED.temperature;
	`, strings.Join(valueStrings, ","))

	if err := tx.Exec(query, valueArgs...).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error upserting miners: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// TODO: ideally raw sql with onconclict,
	// 	     only one transaction should go here
	// if len(miners) == 0 {
	// 	fmt.Println("No miners to update")
	// 	return nil
	// }

	// shardTable := fmt.Sprintf("miners_%02d", miners[0].FleetID)
	// var valueStrings []string
	// var valueArgs []interface{}

	// for _, miner := range miners {
	// 	valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
	// 	valueArgs = append(valueArgs, miner.ID, miner.Stats.HashRate, miner.MinerType, miner.ModelName, miner.Mode, miner.Status, miner.FleetID)
	// }

	// query := fmt.Sprintf(`
	// 	INSERT INTO %s (id, hash_rate, miner_type, model_name, mode, status, fleet_id)
	// 	VALUES %s
	// 	ON CONFLICT (id) DO UPDATE
	// 	SET hash_rate = EXCLUDED.hash_rate,
	// 		miner_type = EXCLUDED.miner_type,
	// 		model_name = EXCLUDED.model_name,
	// 		mode = EXCLUDED.mode,
	// 		status = EXCLUDED.status;
	// `, shardTable, strings.Join(valueStrings, ","))

	// if err := tx.Exec(query, valueArgs...).Error; err != nil {
	// 	return fmt.Errorf("error upserting miners: %w", err)
	// }

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
	// }

	return nil
}
