// backend/routes/remote_control_routes.go

package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/controllers"
	"github.com/shuttlersit/ads-player/backend/middleware"
	"github.com/shuttlersit/ads-player/backend/services"
)

// RegisterRemoteControlRoutes registers remote control routes
func RegisterRemoteControlRoutes(router *gin.Engine, advertisementController *controllers.AdvertisementController, playlistService *services.DefaultPlaylistService, playbackService *services.PlaybackService) {

	// WebSocket route
	webSocketController := controllers.NewRemoteControlController(playlistService, playbackService)
	router.GET("/ws", webSocketController.WebSocketHandler)

	// API route for remote control
	apiController := controllers.NewRemoteControlController(playlistService, playbackService)
	router.GET("/api/control", apiController.APIControlHandler)

	// Register WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		advertisementController.HandleWebSocketConnection(c)
	})

	// Register API endpoints with authentication and authorization middleware
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.AuthorizeRequest()) // Add your authentication middleware

	apiGroup.POST("/pause", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to pause the advertisement playback
		fmt.Println("Advertisement playback paused.")
		c.JSON(200, gin.H{"message": "Playback paused"})
	})

	apiGroup.POST("/resume", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to resume the advertisement playback
		fmt.Println("Advertisement playback resumed.")
		c.JSON(200, gin.H{"message": "Playback resumed"})
	})

	apiGroup.POST("/skip", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to skip to the next advertisement
		fmt.Println("Skipping to the next advertisement.")
		c.JSON(200, gin.H{"message": "Skipping to the next advertisement"})
	})

	apiGroup.POST("/volume_up", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to increase the volume
		fmt.Println("Volume increased.")
		c.JSON(200, gin.H{"message": "Volume increased"})
	})

	apiGroup.POST("/volume_down", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to decrease the volume
		fmt.Println("Volume decreased.")
		c.JSON(200, gin.H{"message": "Volume decreased"})
	})

	apiGroup.POST("/mute", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to mute the audio
		fmt.Println("Audio muted.")
		c.JSON(200, gin.H{"message": "Audio muted"})
	})

	apiGroup.POST("/unmute", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to unmute the audio
		fmt.Println("Audio unmuted.")
		c.JSON(200, gin.H{"message": "Audio unmuted"})
	})

	apiGroup.POST("/info", func(c *gin.Context) {
		// Authorization check (if needed)
		// Add logic to request information about the current playback state
		// Return the information as JSON in the response
		playbackInfo := map[string]interface{}{
			"status":           "playing",   // Replace with the actual playback status
			"current_ad":       "Ad123",     // Replace with the ID or details of the currently playing advertisement
			"current_playlist": "Playlist1", // Replace with the ID or details of the currently playing playlist
			// Add more relevant information as needed
		}
		c.JSON(200, playbackInfo)
	})
}
