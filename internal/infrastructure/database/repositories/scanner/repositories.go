package repositories

import (
	"context"

	"gorm.io/gorm"
)

// TODO: add ListActiveScanners

type ScannerRepositoryInterface interface {
	Create(ctx context.Context, scanner *Scanner) error
	FindByID(ctx context.Context, id uint) (*Scanner, error)
	Update(ctx context.Context, scanner *Scanner) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*Scanner, error)
}

type ScannerRepository struct {
	db *gorm.DB
}

func NewScannerRepository(db *gorm.DB) *ScannerRepository {
	return &ScannerRepository{
		db: db,
	}
}

func (r *ScannerRepository) Create(ctx context.Context, scanner *Scanner) error {
	return r.db.WithContext(ctx).Create(scanner).Error
}

func (r *ScannerRepository) FindByID(ctx context.Context, id uint) (*Scanner, error) {
	var scanner Scanner
	result := r.db.WithContext(ctx).First(&scanner, id)
	return &scanner, result.Error
}

func (r *ScannerRepository) Update(ctx context.Context, scanner *Scanner) error {
	return r.db.WithContext(ctx).Save(scanner).Error
}

func (r *ScannerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Scanner{}, id).Error
}

func (r *ScannerRepository) List(ctx context.Context) ([]*Scanner, error) {

	var scanners []*Scanner

	result := r.db.WithContext(ctx).Preload("Alert").Find(&scanners)

	return scanners, result.Error
}
