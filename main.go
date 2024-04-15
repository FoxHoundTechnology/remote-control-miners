package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	postgres "foxhound/internal/infrastructure/database/postgres"
	"log"

	// TODO: db migration/seed
	miner "foxhound/internal/infrastructure/database/repositories/miner"
	scanner "foxhound/internal/infrastructure/database/repositories/scanner"
)

func main() {

	postgresDB := postgres.Init()
	devMigrate(postgresDB)

	err := postgresDB.AutoMigrate(
		&scanner.Scanner{},
		&miner.Fleet{},
		&miner.Miner{},
		&miner.Pool{},
		&miner.TemperatureSensor{},
		&miner.FanSensor{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	router := gin.Default()

	router.Run(":8080")

}

// TODO: migration/seed
func devMigrate(db *gorm.DB) {

	// err := db.Exec("CREATE TYPE miner_type AS ENUM ('antimner_cgi')")
	// if err != nil {
	// 	fmt.Println("type already exists (expected)")
	// }

	//..... add more types here until all the custom types are finalized
}
