package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	miner_domain "github.com/FoxHoundTechnology/remote-control-miners/internal/application/miner/domain"
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

// TODO: FIXME
func (r *MinerRepository) CreateMinersInBatch(miners []*Miner) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, miner := range miners {

		conditions := map[string]interface{}{
			"mac_address": miner.Miner.MacAddress,
			"fleet_id":    miner.FleetID,
		}

		if err := tx.Where(conditions).Save(&miner).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error saving miner: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (r *MinerRepository) UpdateMinersInBatch(miners []*Miner) error {

	if len(miners) == 0 {
		return nil
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Prepare the values for bulk upsert
	valueStrings := make([]string, 0, len(miners))
	valueArgs := make([]interface{}, 0, len(miners)*16) // Adjusted for 16 fields

	for _, miner := range miners {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			miner.Miner.MacAddress,
			miner.FleetID,
			time.Now(),
			miner.Miner.IPAddress,
			miner.Stats.HashRate,
			miner.Stats.RateIdeal,
			miner.Stats.Uptime,
			miner.Config.Username,
			miner.Config.Password,
			miner.Config.Firmware,
			miner.MinerType,
			miner.ModelName,
			miner.Mode,
			miner.Status,
			miner.Fan,
			miner.Temperature,
		)
	}

	// Construct the SQL query
	query := `
			 INSERT INTO public.miners (
				 mac_address, fleet_id, updated_at, ip_address, 
				 hash_rate, rate_ideal, uptime, username, password, firmware, 
				 miner_type, model_name, mode, status, fan, temperature
			 ) 
			 VALUES %s
			 ON CONFLICT (mac_address, fleet_id) DO UPDATE SET
				 updated_at = EXCLUDED.updated_at,
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
				 temperature = EXCLUDED.temperature
			 `
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

	// TODO: POOL UPDATE
	if err := tx.Exec(query, valueArgs...).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating miners: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *MinerRepository) UpdatePoolsInBatch(miners []*Miner) error {
	if len(miners) == 0 {
		return nil
	}
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	valueStrings := make([]string, 0)
	valueArgs := make([]interface{}, 0)

	for _, miner := range miners {
		// Skip miners with no pools
		if miner == nil || len(miner.Pools) == 0 {
			continue
		}

		for i := 0; i < 3; i++ { // Assuming a maximum of 3 pools per miner
			var pool Pool
			if i < len(miner.Pools) {
				pool = miner.Pools[i]
			} else {
				// If this pool index doesn't exist, use empty values
				pool = Pool{Pool: miner_domain.Pool{}}
			}

			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				miner.Miner.MacAddress,
				i,
				pool.Pool.Url,
				pool.Pool.User,
				pool.Pool.Pass,
				pool.Pool.Status,
				pool.Pool.Accepted,
				pool.Pool.Rejected,
				pool.Pool.Stale,
				time.Now(),
			)
		}
	}

	if len(valueStrings) == 0 {
		// No pools to update
		tx.Rollback()
		return nil
	}

	query := `
    INSERT INTO public.pools (
        miner_mac_address, index, url, user, pass, status, accepted, rejected, stale, updated_at
    ) 
    VALUES %s
    ON CONFLICT (miner_mac_address, index) DO UPDATE SET
        url = EXCLUDED.url,
        user = EXCLUDED.user,
        pass = EXCLUDED.pass,
        status = EXCLUDED.status,
        accepted = EXCLUDED.accepted,
        rejected = EXCLUDED.rejected,
        stale = EXCLUDED.stale,
        updated_at = EXCLUDED.updated_at
    `

	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

	if err := tx.Exec(query, valueArgs...).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating pools: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
