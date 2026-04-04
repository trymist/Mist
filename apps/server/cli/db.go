package cli

import (
	"fmt"
	"os"

	"github.com/trymist/mist/db"
	"github.com/trymist/mist/models"
)

const dbPath = "/var/lib/mist/mist.db"

func initDB() error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found at %s. Please ensure Mist is installed and running", dbPath)
	}

	dbInstance, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	models.SetDB(dbInstance)
	return nil
}
