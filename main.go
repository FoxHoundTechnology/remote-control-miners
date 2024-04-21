package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"

	postgres "foxhound/internal/infrastructure/database/postgres"

	// TODO: db migration/seed
	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"
)

// TODO: R&D for pool library's memory leak
func main() {

	postgresDB := postgres.Init()

	err := postgresDB.AutoMigrate(
		// NOTE: The order matters
		&fleet_repo.Fleet{},

		&scanner_repo.Scanner{},
		&scanner_repo.Alert{},
		&scanner_repo.AlertLog{},

		&miner_repo.Miner{},
		&miner_repo.Pool{},
		&miner_repo.TemperatureSensor{},
		&miner_repo.FanSensor{},
		&miner_repo.MinerLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	DevMigrate(postgresDB)

	router := gin.Default()

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure that resources are freed on return

	fleetRepo := fleet_repo.NewFleetRepository(postgresDB)
	// minerRepo := miner_repo.NewMinerRepository(postgresDB)
	// scannerRepo := scanner_repo.NewScannerRepository(postgresDB)

	pool := pond.New(100, 1000)
	// go func() { }()
	// ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {

		case <-ctx.Done():
			fmt.Println("Shutting down the scheduled tasks...")
			return

		case <-ticker.C:
			fmt.Println("Running scheduled tasks...")
			fleets, err := fleetRepo.List()
			if err != nil {
				// alert layer =  4
				fmt.Println("Error getting fleet list:", err)
				continue
			}
			fmt.Println("fleet list", fleets)

			// spawn a goroutine for each fleet->scanner (which is associated with a scanner and a list of miners)
			for _, fleet := range fleets {
				pool.Submit(func() {
					fmt.Printf("Processing scanner ID: %d\n", fleet.ID)
					fmt.Println("fleet", fleet)

				})
			}

			// (in each goroutine)
			// ARP scan within the ip range
			// using controller to get the miner list
			// while retrieving the raw response and injest it to the miner payload struct
			// methods that will be used is the followings:

			// 1, CheckSystemInfo

			// 2, CheckStats

			// 3, CheckPools

			// 4, CheckConfig

			// Using alert service to check the alert conditions
			// (5), CheckAlerts

			// (5-A), Condition met, process the alert action and log the alert activity

			// (5-B), Condition not met, update the alert state with the lastUpdatedAt timestamp

			// After updating the miner payload (in the application layer),
			// update the miner payload to the database with upsert operation

			// kill go routines after the fleet->scanner->miner list is processed

		}
	}

}
