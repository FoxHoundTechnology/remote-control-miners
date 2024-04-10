package repositories

import (
	"context"

	"gorm.io/gorm"
)

type FleetRepository struct {
	db *gorm.DB
}

type RepositoryInterface interface {
	Create(ctx context.Context, fleet *Fleet) error
	FindByID(ctx context.Context, id uint) (*Fleet, error)
	Update(ctx context.Context, fleet *Fleet) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*Fleet, error)
}

func (r *FleetRepository) Create(ctx context.Context, fleet *Fleet) error {
	return r.db.WithContext(ctx).Create(fleet).Error
}

func (r *FleetRepository) FindByID(ctx context.Context, id uint) (*Fleet, error) {
	var fleet Fleet
	result := r.db.WithContext(ctx).First(&fleet, id)
	return &fleet, result.Error
}

func (r *FleetRepository) Update(ctx context.Context, fleet *Fleet) error {
	return r.db.WithContext(ctx).Save(fleet).Error
}

func (r *FleetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Fleet{}, id).Error
}

func (r *FleetRepository) List(ctx context.Context) ([]*Fleet, error) {
	var fleets []*Fleet
	result := r.db.WithContext(ctx).Find(&fleets)
	return fleets, result.Error
}
