// backend/services/video_service.go

package services

import (
	"errors"
	"fmt"

	"github.com/shuttlersit/ads-player/backend/models"
)

// VideoService provides methods related to video management
type VideoServiceInterface interface {
	// Include any other video-related methods
	GetVideoURL(videoID uint) (string, error)
	SetVideoURL(videoID uint, s3Key string) error
	generateS3URL(s3Key string) string
	SetCurrentVideoURL(videoURL string) string
	GetCurrentVideoURL() string
	GetVideoByID(id uint) (*models.Video, error)
	GetRelatedVideos(videoID uint) ([]*models.Video, error)
	GetMostLikedVideos(limit int) ([]*models.Video, error)
	GetVideosByUploader(userID uint) ([]*models.Video, error)
	GetVideosByCategory(categoryID uint) ([]*models.Video, error)
	IncrementViews(videoID uint) error
	IncrementLikes(videoID uint) error
	IncrementDislikes(videoID uint) error
	GetTopTags(limit int) ([]string, error)
	CreateVideo(video *models.Video) error
	UpdateVideo(video *models.Video) error
	GetVideoViews(videoID uint) (uint, error)
	GetVideoLikes(videoID uint) (uint, error)
	GetVideoDislikes(videoID uint) (uint, error)
	// Add any other methods as needed
}

// VideoService handles video-related operations.
type VideoService struct {
	videoDBModel models.VideoDBModel
	videoPlayer  *RTSPVideoPlayer
	// Add any dependencies or data needed for the service
}

// NewVideoService creates a new instance of VideoService.
func NewVideoService(videoDBModel models.VideoDBModel, videoPlayer *RTSPVideoPlayer) *VideoService {
	return &VideoService{
		videoDBModel: videoDBModel,
		videoPlayer:  videoPlayer,
	}
}

// ... (existing methods)

// New methods

// GetVideoByID retrieves a video by its ID.
func (vs *VideoService) GetVideoByIDInternal(id uint) (*models.Video, error) {
	video, err := vs.videoDBModel.GetVideoByID(id)
	if err != nil {
		return nil, errors.New("failed to fetch video by ID")
	}
	return video, nil
}

func (vs *VideoService) GetRelatedVideosInternal(videoID uint) ([]*models.Video, error) {
	return vs.videoDBModel.GetRelatedVideos(videoID)
}

func (vs *VideoService) GetMostLikedVideosInternal(limit int) ([]*models.Video, error) {
	return vs.videoDBModel.GetMostLikedVideos(limit)
}

func (vs *VideoService) GetVideosByUploaderInternal(userID uint) ([]*models.Video, error) {
	return vs.videoDBModel.GetVideosByUploader(userID)
}

func (vs *VideoService) GetVideosByCategoryInternal(categoryID uint) ([]*models.Video, error) {
	return vs.videoDBModel.GetVideosByCategory(categoryID)
}

func (vs *VideoService) IncrementViewsInternal(videoID uint) error {
	return vs.videoDBModel.IncrementViews(videoID)
}

func (vs *VideoService) IncrementLikesInternal(videoID uint) error {
	return vs.videoDBModel.IncrementLikes(videoID)
}

func (vs *VideoService) IncrementDislikesInternal(videoID uint) error {
	return vs.videoDBModel.IncrementDislikes(videoID)
}

func (vs *VideoService) GetTopTagsInternal(limit int) ([]string, error) {
	return vs.videoDBModel.GetTopTags(limit)
}

// New methods
func (vs *VideoService) CreateVideoInternal(video *models.Video) error {
	return vs.videoDBModel.CreateVideo(video)
}

func (vs *VideoService) UpdateVideoInternal(video *models.Video) error {
	return vs.videoDBModel.UpdateVideo(video)
}

func (vs *VideoService) GetVideoViewsInternal(videoID uint) (uint, error) {
	// Implement logic to fetch views for the specified videoID
	return 0, nil
}

func (vs *VideoService) GetVideoLikesInternal(videoID uint) (uint, error) {
	// Implement logic to fetch likes for the specified videoID
	return 0, nil
}

func (vs *VideoService) GetVideoDislikesInternal(videoID uint) (uint, error) {
	// Implement logic to fetch dislikes for the specified videoID
	return 0, nil
}

// GetVideoURL retrieves the URL of a video by its ID.
func (vs *VideoService) GetVideoURLInternal(videoID uint) (string, error) {
	video, err := vs.videoDBModel.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	return video.URL, nil
}

// SetVideoURL sets the S3 key for a video and updates the S3 URL.
func (vs *VideoService) SetVideoURLInternal(videoID uint, s3Key string) error {
	video, err := vs.videoDBModel.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Assuming a function to generate S3 URL based on the key
	s3URL := vs.videoDBModel.generateS3URL(s3Key)

	// Update the video's S3Key and URL
	video.S3Key = s3Key
	video.URL = s3URL

	// Save the updated video
	if err := vs.videoDBModel.UpdateVideo(video); err != nil {
		return err
	}

	return nil
}

// generateS3URL generates the S3 URL based on the S3 key.
func generateS3URLInternal(s3Key string) string {
	// Implement your logic to generate S3 URL based on the key
	// Example: Concatenate the key with the S3 base URL
	s3BaseURL := "https://s3.example.com/"
	return fmt.Sprintf("%s%s", s3BaseURL, s3Key)
}

// SetCurrentVideoURL sets the current video URL for streaming.
func (vs *VideoService) SetCurrentVideoURLInternal(url string) error {
	// Implement logic to set the current video URL using the video player.
	return vs.videoPlayer.Play(url)
}

// GetCurrentVideoURL gets the current video URL from the video service
func (vs *VideoService) GetCurrentVideoURLInternal() string {
	// Implement the logic to get the current video URL from the video service
	// Return the current video URL

	// Simulate getting the current video URL
	fmt.Printf("Getting current video URL from video service: %s\n", vs.videoDBModel.GetCurrentVideoURL())

	// Return the current video URL
	return vs.videoDBModel.GetCurrentVideoURL()
}
