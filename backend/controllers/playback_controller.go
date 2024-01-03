// backend/controllers/playback_controller.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/services"
)

// PlaybackController handles HTTP requests related to playback
type PlaybackController struct {
	PlaybackService services.PlaybackServiceImpl
}

// NewPlaybackController creates a new PlaybackController
func NewPlaybackController(playbackService services.PlaybackServiceImpl) *PlaybackController {
	return &PlaybackController{
		PlaybackService: playbackService,
	}
}

// RegisterRoutes registers playback routes using the Gin framework
func (plc *PlaybackController) RegisterRoutes(r *gin.Engine) {
	playback := r.Group("/playback")
	{
		playback.GET("/current", plc.GetCurrentVideo)
		playback.GET("/next", plc.GetNextVideo)
		playback.POST("/play/:id", plc.PlayVideo)
	}
}

// GetCurrentVideo handles the HTTP request to retrieve the currently playing video.
func (plc *PlaybackController) GetCurrentVideo(c *gin.Context) {
	// Implement logic to fetch and return the currently playing video
	currentVideo, err := plc.PlaybackService.GetCurrentVideo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current video"})
		return
	}
	c.JSON(http.StatusOK, currentVideo)
}

// GetNextVideo handles the HTTP request to retrieve the next video in the playlist.
func (plc *PlaybackController) GetNextVideo(c *gin.Context) {
	// Implement logic to fetch and return the next video in the playlist
	nextVideo, err := plc.PlaybackService.GetNextVideo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve next video"})
		return
	}
	c.JSON(http.StatusOK, nextVideo)
}

// PlayVideo handles the HTTP request to play a specific video.
func (plc *PlaybackController) PlayVideo(c *gin.Context) {
	videoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Implement logic to play the specified video
	if err := plc.PlaybackService.PlayVideo(uint(videoID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to play video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video playback started"})
}
