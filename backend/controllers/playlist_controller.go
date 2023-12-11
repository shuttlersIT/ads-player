// backend/controllers/playlist_controller.go

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/models"
	"gorm.io/gorm"
)

// PlaylistController handles CRUD operations for playlists
type PlaylistController struct {
	DB *gorm.DB
}

// NewPlaylistController creates a new PlaylistController
func NewPlaylistController(db *gorm.DB) *PlaylistController {
	return &PlaylistController{
		DB: db,
	}
}

// GetPlaylists retrieves all playlists
func (pc *PlaylistController) GetPlaylists(c *gin.Context) {
	var playlists []models.Playlist
	if err := pc.DB.Find(&playlists).Error; err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, playlists)
}

// GetPlaylistByID retrieves a playlist by ID
func (pc *PlaylistController) GetPlaylistByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var playlist models.Playlist
	if err := pc.DB.Preload("Videos").First(&playlist, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, playlist)
}

// CreatePlaylist creates a new playlist
func (pc *PlaylistController) CreatePlaylist(c *gin.Context) {
	var playlist models.Playlist
	if err := c.ShouldBindJSON(&playlist); err != nil {
		c.AbortWithStatus(400)
		return
	}
	if err := pc.DB.Create(&playlist).Error; err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, playlist)
}

// UpdatePlaylist updates a playlist by ID
func (pc *PlaylistController) UpdatePlaylist(c *gin.Context) {
	id := c.Params.ByName("id")
	var playlist models.Playlist
	if err := pc.DB.First(&playlist, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	if err := c.ShouldBindJSON(&playlist); err != nil {
		c.AbortWithStatus(400)
		return
	}
	if err := pc.DB.Save(&playlist).Error; err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, playlist)
}

// DeletePlaylist deletes a playlist by ID
func (pc *PlaylistController) DeletePlaylist(c *gin.Context) {
	id := c.Params.ByName("id")
	var playlist models.Playlist
	if err := pc.DB.First(&playlist, id).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	if err := pc.DB.Delete(&playlist).Error; err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
