package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection. If DATABASE_URL is provided it
// uses Postgres; otherwise it falls back to an in-memory sqlite database which
// is useful for tests and CI.
func Init(models ...interface{}) {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	if dsn != "" {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
