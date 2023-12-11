package database

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shuttlersit/ads-player/backend/models"
)

var db *gorm.DB
var err error
var status string

func ConnectMySqlite() (string, *gorm.DB) {

	// Replace with your database credentials
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Playlist{}, &models.Video{})

	fmt.Println("Database is on.")
	return status, db
}

// Check if a table exists in the database
func TableExists(db *sql.DB, tableName string) bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
