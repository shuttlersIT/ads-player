// backend/routes/playlist_routes.go

package routes

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shuttlersit/ads-player/backend/controllers"
)

// RegisterPlaylistRoutes registers routes related to playlists
func RegisterPlaylistRoutes(r *gin.Engine, db *gorm.DB) {
	playlistController := controllers.NewPlaylistController(db)

	playlists := r.Group("/playlists")
	{
		playlists.GET("", playlistController.GetPlaylists)
		playlists.GET("/:id", playlistController.GetPlaylistByID)
		playlists.POST("", playlistController.CreatePlaylist)
		playlists.PUT("/:id", playlistController.UpdatePlaylist)
		playlists.DELETE("/:id", playlistController.DeletePlaylist)
	}
}
