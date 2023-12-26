// backend/controllers/playlist_controller.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/models"
	"github.com/shuttlersit/ads-player/backend/services"
	"gorm.io/gorm"
)

// PlaylistController handles HTTP requests related to playlists
type PlaylistController struct {
	playlistService services.DefaultPlaylistService
}

// NewPlaylistController creates a new PlaylistController
func NewPlaylistController(playlistService services.DefaultPlaylistService) *PlaylistController {
	return &PlaylistController{
		playlistService: playlistService,
	}
}

// PlaylistController handles CRUD operations for playlists
type PlaylistDBController struct {
	PlaylistService *services.DefaultPlaylistService
	DB              *gorm.DB
}

// NewPlaylistController creates a new PlaylistController
func NewPlaylistDBController(db *gorm.DB, playlistService *services.DefaultPlaylistService) *PlaylistDBController {
	return &PlaylistDBController{
		DB:              db,
		PlaylistService: playlistService,
	}
}

// GetPlaylists retrieves all playlists
func (pc *PlaylistDBController) GetPlaylists(c *gin.Context) {
	var playlists []models.Playlist
	if err := pc.DB.Find(&playlists).Error; err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, playlists)
}

// UpdatePlaylist updates a playlist by ID
func (pc *PlaylistDBController) UpdatePlaylist(c *gin.Context) {
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
func (pc *PlaylistDBController) DeletePlaylist(c *gin.Context) {
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

// RegisterRoutesA registers playlist routes using the Gin framework
func (pc *PlaylistDBController) RegisterRoutes(r *gin.Engine) {
	playlists := r.Group("/playlists")
	{
		playlists.GET("", pc.GetPlaylists)
		playlists.GET("/:id", pc.GetPlaylistByID)
		playlists.POST("", pc.CreatePlaylist)
		playlists.PUT("/:id", pc.UpdatePlaylist)
		playlists.DELETE("/:id", pc.DeletePlaylist)
	}
}

// GetPlaylistByID handles the HTTP request to retrieve a playlist by ID.
func (pc *PlaylistDBController) GetPlaylistByID(ctx *gin.Context) {
	playlistID, _ := strconv.Atoi(ctx.Param("id"))
	playlist, err := pc.PlaylistService.GetPlaylistByID(uint(playlistID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	ctx.JSON(http.StatusOK, playlist)
}

// GetAllPlaylists handles the HTTP request to retrieve all playlists.
func (pc *PlaylistDBController) GetAllPlaylists2(ctx *gin.Context) {
	playlists, err := pc.PlaylistService.GetAllPlaylists()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playlists"})
		return
	}
	ctx.JSON(http.StatusOK, playlists)
}

// CreatePlaylist handles the HTTP request to create a new playlist.
func (pc *PlaylistDBController) CreatePlaylist2(ctx *gin.Context) {
	var playlist models.Playlist
	if err := ctx.ShouldBindJSON(&playlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	createdPlaylist, err := pc.PlaylistService.CreatePlaylist(&playlist)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	ctx.JSON(http.StatusCreated, createdPlaylist)
}

// GetAllPlaylists retrieves all playlists
func (pc *PlaylistDBController) GetAllPlaylists3(c *gin.Context) {
	// Fetch all playlists from the database
	playlists, err := models.NewPlaylistModel(pc.DB).GetAllPlaylists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch playlists"})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// CreatePlaylist creates a new playlist
func (pc *PlaylistDBController) CreatePlaylist3(c *gin.Context) {
	// Bind the request body to a Playlist struct
	var newPlaylist models.Playlist
	if err := c.ShouldBindJSON(&newPlaylist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create the playlist in the database
	err := models.NewPlaylistModel(pc.DB).CreatePlaylist(&newPlaylist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	c.JSON(http.StatusCreated, &newPlaylist)
}

// GetAllPlaylists retrieves all playlists
func (pc *PlaylistDBController) GetAllPlaylists(c *gin.Context) {
	// Fetch all playlists using the PlaylistService
	playlists, err := pc.PlaylistService.GetAllPlaylists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch playlists"})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// CreatePlaylist creates a new playlist
func (pc *PlaylistDBController) CreatePlaylist(c *gin.Context) {
	// Bind the request body to a Playlist struct
	var newPlaylist models.Playlist
	if err := c.ShouldBindJSON(&newPlaylist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create the playlist using the PlaylistService
	createdPlaylist, err := pc.PlaylistService.CreatePlaylist(&newPlaylist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	c.JSON(http.StatusCreated, createdPlaylist)
}
