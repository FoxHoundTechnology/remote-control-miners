package main

import (
	"fmt"
	postgres "foxhound/internal/infrastructure/database/postgres"
	"log"

	// TODO: db migration/seed
	alert "foxhound/internal/infrastructure/database/repositories/alert"
	miner "foxhound/internal/infrastructure/database/repositories/miner"
	scanner "foxhound/internal/infrastructure/database/repositories/scanner"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	postgresDB := postgres.Init()
	devMigrate(postgresDB)

	err := postgresDB.AutoMigrate(
		&alert.Alert{},
		&scanner.Scanner{},
		&miner.Miner{},
		&miner.Pool{},
		&miner.Temperature{},
		&miner.TemperatureSensor{},
		&miner.Fan{},
		&miner.FanSensor{},
		&miner.Fleet{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	// antMinerCGI := service.AntminerCGI{
	// 	Miner:       &domain.Miner{},
	// 	Mode:        domain.Mode(domain.NormalMode),
	// 	Status:      domain.Status(domain.Online),
	// 	Config:      &domain.Config{},
	// 	Stats:       &domain.Stats{},
	// 	Pools:       &domain.Pool{},
	// 	Temperature: &domain.Temperature{},
	// 	Fan:         &domain.Fan{},
	// }
	// list := []domain.MinerController{
	// 	&antMinerCGI,
	// }
	// for _, miner := range list {
	// 	fmt.Println(miner)
	// }

	router := gin.Default()

	router.Run(":8080")

}

// TODO: migration/seed
func devMigrate(db *gorm.DB) {

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

	err = db.Exec("CREATE TYPE miner_type AS ENUM ('antimner_cgi')")
	if err != nil {
		fmt.Println("type already exists (expected)")
	}

	//..... add more types here until all the custom types are finalized
}
