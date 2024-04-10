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

type PostgresConnectionSettings struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var settings = PostgresConnectionSettings{
	Host:     "postgres_db", // Placeholder values
	Port:     "5432",        // Placeholder values
	User:     "user",        // Placeholder values
	Password: "1234",        // Placeholder values
	Database: "postgres",    // Placeholder values
}

// TODO: automate instantiation with init
func Init() *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		settings.Host, settings.User, settings.Password, settings.Database, settings.Port)

	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Successfully established the connection with the PostgreSQL database.")
	return postgresDB
}

func Close(db *gorm.DB) {
	postgresDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}
	postgresDB.Close()
	fmt.Println("Successfully closed the connection with the PostgreSQL database.")
}
