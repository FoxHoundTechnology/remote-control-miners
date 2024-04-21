package main

import (
	"fmt"
	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	miner_domain "foxhound/internal/application/miner/domain"
	scanner_domain "foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

func DevMigrateFleet(db *gorm.DB) error {
	fleet := fleet_repo.Fleet{
		Name: "test_fleet",
	}

	db.Where(fleet_repo.Fleet{
		Name: "test_fleet",
	}).FirstOrCreate(&fleet)

	fmt.Println("fleet ID", fleet.ID)

	miners := []miner_repo.Miner{
		{
			Miner: miner_domain.Miner{
				MacAddress: "00:1A:2B:3C:4D:5E",
				IPAddress:  "10.0.1.100",
				Owner:      "Owner1",
			},
			Stats: miner_domain.Stats{
				HashRate:  5000.0,
				RateIdeal: 5200.0,
				Uptime:    100000,
			},
			Config: miner_domain.Config{
				Username: "user",
				Password: "pass",
				Firmware: "v1.2.3",
			},
			Mode:    miner_domain.NormalMode,
			Status:  miner_domain.Online,
			FleetID: fleet.ID,
		},
		{
			Miner: miner_domain.Miner{
				MacAddress: "00:1A:2B:3C:4D:5T",
				IPAddress:  "192.168.1.101",
				Owner:      "Owner2",
			},
			Stats: miner_domain.Stats{
				HashRate:  5000.0,
				RateIdeal: 5200.0,
				Uptime:    100000,
			},
			Config: miner_domain.Config{
				Username: "user",
				Password: "pass",
				Firmware: "v1.2.3",
			},
			Mode:    miner_domain.NormalMode,
			Status:  miner_domain.Online,
			FleetID: fleet.ID,
		},
	}

	// TODO: batch or tx. operation
	for _, miner := range miners {
		result := db.Where("mac_address = ?", miner.Miner.MacAddress).First(&miner.Miner)

		if result.RowsAffected == 0 {
			err := db.Create(&miner).Error
			fmt.Println("ERROR IN ROWS", err)

		}
	}

	return nil
}

func DevMigrateScanerAndAlert(db *gorm.DB) error {

	scanner := scanner_repo.Scanner{
		Scanner: scanner_domain.Scanner{
			Name:     "scanner test",
			StartIP:  "10.0.0.137",
			EndIP:    "10.0.0.180",
			Active:   true,
			Location: "TEST LOCATION",
		},
		Config: scanner_domain.Config{
			Interval: 5,
			Username: "user",
			Password: "pass",
		},
		MinerType: scanner_domain.AntminerCgi,
		Owner:     "test owner",
		// FleetID
	}

	result := db.Where("name = ?", scanner.Scanner.Name).First(&scanner)
	if result.RowsAffected == 0 {
		err := db.Create(&scanner).Error
		fmt.Println("ERROR IN ROWS", err)
	}

	alertA := scanner_repo.Alert{
		Name:      "alert A",
		Value:     50, // 50%
		Threshold: scanner_domain.ThresholdRate,
		Condition: scanner_domain.Hashrate,
		Action:    scanner_domain.Reboot,
		Layer:     scanner_domain.InfoAlert,
		Log: []scanner_repo.AlertLog{
			{
				Log: "test log from a",
			},
		},
		ScannerID: scanner.ID,
	}
	alertB := scanner_repo.Alert{
		Name:      "alert B",
		Value:     10, // 50 machines
		Threshold: scanner_domain.ThresholdCount,
		Condition: scanner_domain.Temperature,
		Action:    scanner_domain.Sleep,
		Layer:     scanner_domain.InfoAlert,
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
