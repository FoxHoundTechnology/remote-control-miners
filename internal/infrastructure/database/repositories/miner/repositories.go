package repositories

import (
	"context"

	"gorm.io/gorm"
)

// TODO: association
// TODO: preload
type MinerRepositoryInterface interface {
	Create(ctx context.Context, miner *Miner) error
	FindByID(ctx context.Context, id uint) (*Miner, error)
	FindByMacAddress(ctx context.Context, macAddress string) (*Miner, error)
	Update(ctx context.Context, miner *Miner) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*Miner, error)
}

type MinerRepository struct {
	db *gorm.DB
}

func (r *MinerRepository) Create(ctx context.Context, miner *Miner) error {
	return r.db.WithContext(ctx).Create(miner).Error
}

func (r *MinerRepository) FindByID(ctx context.Context, id uint) (*Miner, error) {
	var miner Miner
	result := r.db.WithContext(ctx).First(&miner, id)
	return &miner, result.Error
}

func (r *MinerRepository) FindByMacAddress(ctx context.Context, macAddress string) (*Miner, error) {
	var miner Miner
	result := r.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&miner)
	return &miner, result.Error
}

func (r *MinerRepository) Update(ctx context.Context, miner *Miner) error {
	return r.db.WithContext(ctx).Save(miner).Error
}

func (r *MinerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Miner{}, id).Error
}

func (r *MinerRepository) List(ctx context.Context) ([]*Miner, error) {
	var miners []*Miner
	result := r.db.WithContext(ctx).Find(&miners)
	return miners, result.Error
}

func (r *MinerRepository) ListByFleetID(ctx context.Context, fleetID uint) ([]*Miner, error) {
	var miners []*Miner
	result := r.db.WithContext(ctx).Where("fleet_id = ?", fleetID).Find(&miners)
	return miners, result.Error
}
