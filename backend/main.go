// backend/main.go

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/database"
	"github.com/shuttlersit/ads-player/backend/models"
	"github.com/shuttlersit/ads-player/backend/routes"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func main() {
	//initiate mysql database
	status, db := database.ConnectMySqlite()
	fmt.Println(status)
	//database.TableExists(db, "tickets")

	// Migrate the schema
	db.AutoMigrate(&models.Playlist{}, &models.Video{})

	// Initialize Gin router
	r := gin.Default()

	// Register playlist routes
	routes.RegisterPlaylistRoutes(r, db)

	// Run the server
	r.Run(":8080")
}
