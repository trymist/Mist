package cmd

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/trymist/mist/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbPath = "/var/lib/mist/mist.db"

func initDB() error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found at %s. Please ensure Mist is installed and running", dbPath)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	models.SetDB(db)
	return nil
}
