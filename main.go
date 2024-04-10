package main

import (
	"fmt"
	postgres "foxhound/internal/infrastructure/database/postgres"
	"log"

	// TODO: db migration/seed
	alert "foxhound/internal/infrastructure/database/repositories/alert"
	fleet "foxhound/internal/infrastructure/database/repositories/fleet"
	miner "foxhound/internal/infrastructure/database/repositories/miner"
	scanner "foxhound/internal/infrastructure/database/repositories/scanner"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	postgresDB := postgres.Init()
	DevMigrate(postgresDB)

	err := postgresDB.AutoMigrate(
		&alert.Alert{},
		&scanner.Scanner{},
		&miner.Miner{},
		&fleet.Fleet{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	router := gin.Default()

	router.Run(":8080")

}

// TODO: migration/seed
func DevMigrate(db *gorm.DB) {
	err := db.Exec(("CREATE TYPE alert_threshold AS ENUM ('count', 'rate')"))
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	err = db.Exec("CREATE TYPE alert_condition AS ENUM ('hashrate', 'temperature', 'fan_speed', 'pool_shares', 'offline_miners', 'missing_hashboards')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	err = db.Exec("CREATE TYPE alert_action AS ENUM ('reboot', 'sleep', 'normal', 'change_pool')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	err = db.Exec("CREATE TYPE alert_layer AS ENUM ('info', 'warning', 'error', 'fatal')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	err = db.Exec("CREATE TYPE miner_status AS ENUM ('online', 'offline', 'disabled', 'warning', 'error')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	err = db.Exec("CREATE TYPE miner_type AS ENUM ('antimner_cgi')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	//..... add more types here until all the custom types are finalized
}
