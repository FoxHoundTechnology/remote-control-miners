package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Logger
// TODO: Handle fatal errors
// TODO: env

type TimescaleDBConnectionSettings struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var settings = TimescaleDBConnectionSettings{
	Host:     "timescale_db", // Placeholder values
	Port:     "5434",         // Placeholder values
	User:     "user",         // Placeholder values
	Password: "1234",         // Placeholder values
	Database: "timescale",    // Placeholder values
}

// TODO: automate instantiation with init
func Init() *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		settings.Host, settings.User, settings.Password, settings.Database, settings.Port)

	TimescaleDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Successfully established the connection with the TimescaleDB database.")
	return TimescaleDB
}

func Close(db *gorm.DB) {
	timescaleDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}
	timescaleDB.Close()
	fmt.Println("Successfully closed the connection with the TimescaleDB database.")
}
