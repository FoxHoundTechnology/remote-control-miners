package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"

	scanner_domain "foxhound/internal/application/scanner/domain"

	"gorm.io/gorm"
)

type Config struct {
	Fleets []struct {
		Name    string `json:"name"`
		Scanner struct {
			Name      string `json:"name"`
			StartIP   string `json:"start_ip"`
			EndIP     string `json:"end_ip"`
			Active    bool   `json:"active"`
			Location  string `json:"location"`
			Username  string `json:"username"`
			Password  string `json:"password"`
			MinerType int    `json:"miner_type"`
			Owner     string `json:"owner"`
		} `json:"scanner"`
		Alert struct {
			Name       string `json:"name"`
			Action     int    `json:"action"`
			Conditions []struct {
				TriggerValue  float64 `json:"trigger_value"`
				MachineCount  int     `json:"machine_count"`
				ThresholdType int     `json:"threshold_type"`
				ConditionType int     `json:"condition_type"`
				LayerType     int     `json:"layer_type"`
			} `json:"condition"`
		} `json:"alert"`
	} `json:"fleets"`
}

func DevMigrate(db *gorm.DB, configFile *os.File) error {

	var loadedConfig Config
	if err := json.NewDecoder(configFile).Decode(&loadedConfig); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
		return err
	}

	// Process each fleet
	for index, fleetConfig := range loadedConfig.Fleets {
		fleet := fleet_repo.Fleet{
			Name: fleetConfig.Name,
		}

		// Check if fleet exists and create if not
		var existingFleet fleet_repo.Fleet
		if err := db.Where("name = ?", fleet.Name).First(&existingFleet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&fleet).Error; err != nil {
					fmt.Println("ERROR IN DevMigrate", err)
					return err
				}
			} else {
				fmt.Println("ERROR IN DevMigrate", err)
				return err
			}
		}

		scanner := scanner_repo.Scanner{
			Name: fleet.Name,
			Scanner: scanner_domain.Scanner{
				StartIP:  fleetConfig.Scanner.StartIP,
				EndIP:    fleetConfig.Scanner.EndIP,
				Active:   true, // fleetConfig.Scanner.Active,
				Location: fleetConfig.Scanner.Location,
			},
			Config: scanner_domain.Config{
				Username: fleetConfig.Scanner.Username,
				Password: fleetConfig.Scanner.Password,
			},
			MinerType: scanner_domain.MinerType(fleetConfig.Scanner.MinerType),
			Owner:     fleetConfig.Scanner.Owner,
			FleetID:   fleet.ID,
		}

		if result := db.Where("name = ?", scanner.Name).First(&scanner_repo.Scanner{}); result.RowsAffected == 0 {
			if err := db.Create(&scanner).Error; err != nil {
				fmt.Println("ERROR IN SCANNER", err)
				return err
			}
		}

		alert := scanner_repo.Alert{
			Name:      fleet.Name,
			Action:    scanner_domain.AlertActionType(fleetConfig.Alert.Action),
			State:     scanner_domain.Monitoring,
			ScannerID: scanner.ID,
		}

		// Insert or update the scanner
		if result := db.Where("name = ?", fleetConfig.Alert.Name).First(&scanner_repo.Alert{}); result.RowsAffected == 0 {
			if err := db.Create(&scanner).Error; err != nil {
				fmt.Println("ERROR IN SCANNER", err)
				return err
			}
		}

		for _, alertCondition := range loadedConfig.Fleets[index].Alert.Conditions {
			condition := scanner_repo.AlertCondition{
				TriggerValue:  scanner_domain.AlertTriggerValue(alertCondition.TriggerValue),
				MachineCount:  scanner_domain.AlertMachineCount(alertCondition.MachineCount),
				ThresholdType: scanner_domain.AlertThresholdType(alertCondition.ThresholdType),
				ConditionType: scanner_domain.AlertConditionType(alertCondition.ConditionType),
				LayerType:     scanner_domain.AlertLayerType(alertCondition.LayerType),
				AlertID:       alert.ID,
			}

			if result := db.Where("condition_type = ?", alertCondition.ConditionType).First(&scanner_repo.Alert{}); result.RowsAffected == 0 {
				if err := db.Create(&alert).Error; err != nil {
					fmt.Println("ERROR IN ALERT", err)
					return err
				}
			}

		}
	}
	return nil
}
