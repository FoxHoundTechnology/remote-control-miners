package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Logger
// TODO: Handle fatal errors

type PostgresConnectionSettings struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Init() *gorm.DB {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, db_name, port)

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
