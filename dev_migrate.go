package main

import (
	"errors"
	"fmt"
	"os"

	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	scanner_domain "foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

func DevMigrate(db *gorm.DB) error {
	fleet := fleet_repo.Fleet{
		Name: "test_fleet",
	}

	err := db.First(&fleet, "name = ?", fleet.Name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := db.Create(&fleet).Error
		if err != nil {
			fmt.Println("ERROR IN FLEET", err)
			return err
		}
	}

	start_ip := os.Getenv("START_IP")
	end_ip := os.Getenv("END_IP")
	temp_user := os.Getenv("TEMP_USER")
	temp_pass := os.Getenv("TEMP_PASS")

	scanner := scanner_repo.Scanner{
		Name: "scanner test",
		Scanner: scanner_domain.Scanner{
			StartIP:  start_ip,
			EndIP:    end_ip,
			Active:   true,
			Location: "TEST LOCATION",
		},
		Config: scanner_domain.Config{
			Interval: 5,
			Username: temp_user,
			Password: temp_pass,
		},
		MinerType: scanner_domain.AntminerCgi,
		Owner:     "test owner",
		FleetID:   fleet.ID,
	}

	result := db.Where("name = ?", scanner.Name).First(&scanner)
	if result.RowsAffected == 0 {
		err := db.Create(&scanner).Error
		fmt.Println("ERROR IN ROWS", err)
	}

	alertA := scanner_repo.Alert{
		Name:   "alert A",
		Action: scanner_domain.Reboot,
		Condition: []scanner_repo.AlertCondition{
			{
				TriggerValue:  50,                  // 50 TH/s
				MachineCount:  100,                 // 100 machines
				ThresholdType: scanner_domain.Rate, // %
				ConditionType: scanner_domain.Hashrate,
				LayerType:     scanner_domain.InfoAlert,
			},
		},

		Log: []scanner_repo.AlertLog{
			{
				Log: "test log from a",
			},
		},
		ScannerID: scanner.ID,
	}
	alertB := scanner_repo.Alert{
		Name:   "alert B",
		Action: scanner_domain.Sleep,
		Condition: []scanner_repo.AlertCondition{
			{
				TriggerValue:  80,                   // 80C
				MachineCount:  100,                  // 100 machines
				ThresholdType: scanner_domain.Count, // machines
				ConditionType: scanner_domain.Temperature,
				LayerType:     scanner_domain.InfoAlert,
			},
		},
		Log: []scanner_repo.AlertLog{
			{
				Log: "test log from b",
			},
		},
		ScannerID: scanner.ID,
	}

	result = db.Where("name = ?", alertA.Name).First(&alertA)
	if result.RowsAffected == 0 {
		db.Create(&alertA)
	}

	result = db.Where("name = ?", alertB.Name).First(&alertB)
	if result.RowsAffected == 0 {
		db.Create(&alertB)
	}

	return nil
}
