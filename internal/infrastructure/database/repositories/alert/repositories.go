package repositories

import (
	"context"

	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

type AlertRepositoryInterface interface {
	Create(ctx context.Context, alert *Alert) error
	FindByID(ctx context.Context, id string) (*Alert, error)
	Update(ctx context.Context, alert *Alert) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Alert, error)
}

func (r *AlertRepository) Create(ctx context.Context, alert *Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *AlertRepository) FindByID(ctx context.Context, id string) (*Alert, error) {
	var alert Alert
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&alert)
	return &alert, result.Error
}

func (r *AlertRepository) Update(ctx context.Context, alert *Alert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

func (r *AlertRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Alert{}, id).Error
}

func (r *AlertRepository) List(ctx context.Context) ([]*Alert, error) {
	var alerts []*Alert
	result := r.db.WithContext(ctx).Find(&alerts)
	return alerts, result.Error
}
