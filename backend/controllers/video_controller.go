// backend/controllers/video_controller.go

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/models"
	"github.com/shuttlersit/ads-player/backend/services"
)

// VideoController handles video-related HTTP requests.
type VideoController struct {
	VideoService *services.VideoService
	//VideoDBModel *models.VideoDBModel
}

// NewVideoController creates a new instance of VideoController.
func NewVideoController(videoService *services.VideoService) *VideoController {
	return &VideoController{
		VideoService: videoService,
		//VideoDBModel: videoDBModel,
	}
}

// New methods

// HandleSetVideoURL handles the logic for setting the S3 key and updating the S3 URL.
func (vh *VideoController) HandleSetVideoURL(videoID uint, s3Key string) error {
	return vh.VideoService.SetVideoURL(videoID, s3Key)
}

// HandleGetVideoURL handles the logic for retrieving the video URL.
func (vh *VideoController) HandleGetVideoURL(videoID uint) (string, error) {
	return vh.VideoService.GetVideoURL(videoID)
}

func (vc *VideoController) GetRelatedVideos(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	videos, err := vc.VideoService.GetRelatedVideos(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch related videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func (vc *VideoController) GetMostLikedVideos(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10 // Default limit
	}

	videos, err := vc.VideoService.GetMostLikedVideos(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch most liked videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func (vc *VideoController) GetVideosByUploader(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	videos, err := vc.VideoService.GetVideosByUploader(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos by uploader"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func (vc *VideoController) GetVideosByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	videos, err := vc.VideoService.GetVideosByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos by category"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func (vc *VideoController) IncrementViews(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	err = vc.VideoService.IncrementViews(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment views"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Views incremented"})
}

func (vc *VideoController) IncrementLikes(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	err = vc.VideoService.IncrementLikes(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Likes incremented"})
}

func (vc *VideoController) IncrementDislikes(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	err = vc.VideoService.IncrementDislikes(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment dislikes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dislikes incremented"})
}

func (vc *VideoController) GetTopTags(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10 // Default limit
	}

	tags, err := vc.VideoService.GetTopTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// New methods
func (vc *VideoController) CreateVideo(c *gin.Context) {
	var input models.Video
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data"})
		return
	}

	err := vc.VideoService.CreateVideo(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Video created successfully"})
}

func (vc *VideoController) UpdateVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	var input models.Video
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video data"})
		return
	}

	existingVideo, err := vc.VideoService.GetVideoByID(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	// Update relevant fields
	existingVideo.Title = input.Title
	existingVideo.Description = input.Description
	// Update other fields as needed

	err = vc.VideoService.UpdateVideo(existingVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video updated successfully"})
}

func (vc *VideoController) GetVideoByID(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	video, err := vc.VideoService.GetVideoByID(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (vc *VideoController) GetVideoAnalytics(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	views, err := vc.VideoService.GetVideoViews(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video views"})
		return
	}

	likes, err := vc.VideoService.GetVideoLikes(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video likes"})
		return
	}

	dislikes, err := vc.VideoService.GetVideoDislikes(uint(videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video dislikes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"views":    views,
		"likes":    likes,
		"dislikes": dislikes,
	})
}
