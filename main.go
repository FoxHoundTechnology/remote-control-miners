package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	//"github.com/alitto/pond"
	postgres "foxhound/internal/infrastructure/database/postgres"
	"log"

	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"

	// TODO: db migration/seed
	fleet_repo "foxhound/internal/infrastructure/database/repositories/fleet"
	miner_repo "foxhound/internal/infrastructure/database/repositories/miner"
	scanner_repo "foxhound/internal/infrastructure/database/repositories/scanner"
)

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

	DevMigrateFleet(postgresDB)
	DevMigrateScanerAndAlert(postgresDB)

	router := gin.Default()

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure that resources are freed on return

	scannerRepo := scanner_repo.NewScannerRepository(postgresDB)

	pool := pond.New(100, 1000)
	go func() {
		// ticker
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Shutting down the scheduled tasks...")
				return

			case <-ticker.C: // This case is selected every 5 minutes
				fmt.Println("Running scheduled tasks...")
				scanners, err := scannerRepo.List()
				if err != nil {
					fmt.Println("Error getting scanner list:", err)
					continue
				}
				fmt.Println("scanner list from db", scanners)
				fmt.Println("Scanner list retrieved, number of scanners:", len(scanners))
				for _, scanner := range scanners {
					sc := scanner // capture range variable
					pool.Submit(func() {
						fmt.Printf("Processing scanner ID: %d\n", sc.ID)
					})
				}
			}
		}
	}()

}
