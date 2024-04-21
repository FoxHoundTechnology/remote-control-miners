package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// TODO: preload with alerts and alert logs
type ScannerRepository struct {
	db *gorm.DB
}

func NewScannerRepository(db *gorm.DB) *ScannerRepository {
	return &ScannerRepository{
		db: db,
	}
}

func (r *ScannerRepository) Upsert(ctx context.Context, scanner *Scanner) (uint, error) {
	err := r.db.Create(scanner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := r.db.Create(scanner).Error
		if err != nil {
			return 0, err
		}
		return scanner.ID, nil
	}

	// Save is a combined function.
	// If save value does not contain primary key,
	// it executes Create. Otherwise it executes Update (with all fields).
	err = r.db.Save(&scanner).Error
	if err != nil {
		return 0, err
	}

	return scanner.ID, nil
}

func (r *ScannerRepository) List() ([]*Scanner, error) {
	var scanners []*Scanner
	err := r.db.Find(&scanners).Error
	if err != nil {
		return nil, err
	}

	return scanners, err
}

func (r *ScannerRepository) UpsertAlert(ctx context.Context, alert *Alert) (uint, error) {
	err := r.db.First(&alert, "name = ?", alert.Name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// with the scanner id
		err := r.db.Create(&alert).Error
		if err != nil {
			return 0, err
		}

		return alert.ID, nil
	}

	return alert.ID, nil
}

// where alert_id = ?
func (r *ScannerRepository) LogAlert(ctx context.Context, alertLog *AlertLog) (uint, error) {
	err := r.db.Create(alertLog).Error
	if err != nil {
		return 0, err
	}

	return alertLog.ID, nil
}
