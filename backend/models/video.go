// backend/models/video.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Video model
type Video struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	URL             string         `json:"url"`
	ThumbnailURL    string         `json:"thumbnailUrl"`
	Views           uint           `json:"views" gorm:"default:0"`
	Likes           uint           `json:"likes" gorm:"default:0"`
	Dislikes        uint           `json:"dislikes" gorm:"default:0"`
	Comments        []Comment      `gorm:"foreignKey:VideoID"`
	PlaylistID      uint           `json:"-"`
	Playlist        Playlist       `json:"playlist"`
	IsAdvertisement bool           `json:"isAdvertisement" gorm:"default:false"`
	Order           int            `json:"order" gorm:"default:0"`
	Tags            []string       `json:"tags" gorm:"type:varchar(255)[]"`
	UploadDate      int            `json:"uploadDate" gorm:"autoCreateTime"`
	Uploader        User           `json:"uploader"`
	Category        Category       `json:"category"`
	PrivacySetting  PrivacySetting `json:"privacySetting" gorm:"embedded"`
	CommentsEnabled bool           `json:"commentsEnabled" gorm:"default:true"`
	RelatedVideos   []RelatedVideo `json:"relatedVideos" gorm:"foreignKey:VideoID"`
	DurationInt     uint           `json:"duration_int"`
	Duration        time.Duration  `json:"duration"`
	CreatorUserID   uint           `json:"creator_user_id"`
	S3Key           string         `json:"s3key"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       time.Time      `gorm:"index" json:"-"`

	// Add more video-related fields as needed
}

// TableName sets the table name for the Video model.
func (Video) TableName() string {
	return "videos"
}

// RelatedVideo model for representing related videos
type RelatedVideo struct {
	gorm.Model
	VideoID        uint
	RelatedVideoID uint
}

// TableName sets the table name for the Video model.
func (RelatedVideo) TableName() string {
	return "related_videos"
}

// Add this to your configuration setup
var S3BaseURL = "https://your-s3-bucket.s3.amazonaws.com/"

type VideoDBModel struct {
	DB        *gorm.DB
	S3BaseURL string // Optionally, you can store S3BaseURL here as well
}

type VideoModel interface {
	GetTotalDuration(videoID uint) (time.Duration, error)
	GetCurrentVideo(playlistID uint) (*Video, error)
	GetNextVideo(playlistID uint) (*Video, error)
	GetVideoByID2(videoID uint) (*Video, error)
	GetPlaybackHistory(userID uint) ([]*Video, error)
	GetPlaylistVideos(playlistID uint) ([]Video, error)
	CreateVideo(video *Video) error
	GetVideoByID(id uint) (*Video, error)
	UpdateVideo(video *Video) error
	DeleteVideo(id uint) error
	GetAllVideos() ([]Video, error)
	GetRelatedVideos(videoID uint) ([]*Video, error)
	GetMostLikedVideos(limit int) ([]*Video, error)
	GetVideosByUploader(userID uint) ([]*Video, error)
	GetVideosByCategory(categoryID uint) ([]*Video, error)
	IncrementViews(videoID uint) error
	IncrementLikes(videoID uint) error
	IncrementDislikes(videoID uint) error
	GetTopTags(limit int) ([]string, error)
	//
	GetVideoURL(videoID uint) (string, error)
	SetVideoURL(videoID uint, s3Key string) error
	generateS3URL(s3Key string) string
	SetCurrentVideoURL(videoURL string) string
	GetCurrentVideoURL() string
	GetVideoViews(videoID uint) (uint, error)
	GetVideoLikes(videoID uint) (uint, error)
	GetVideoDislikes(videoID uint) (uint, error)
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewVideoDBModel(db *gorm.DB, s3BaseURL string) *VideoDBModel {
	return &VideoDBModel{
		DB:        db,
		S3BaseURL: s3BaseURL,
	}
}

// CreateVideo creates a new video in the database.
func (m *VideoDBModel) CreateVideo(video *Video) error {
	return m.DB.Create(video).Error
}

// GetVideoByID retrieves a video by its ID.
func (m *VideoDBModel) GetVideoByID(id uint) (*Video, error) {
	var video Video
	err := m.DB.Where("id = ?", id).First(&video).Error
	return &video, err
}

// UpdateVideo updates the details of an existing video.
func (m *VideoDBModel) UpdateVideo(video *Video) error {
	return m.DB.Save(video).Error
}

// DeleteVideo deletes a video from the database.
func (m *VideoDBModel) DeleteVideo(id uint) error {
	return m.DB.Delete(&Video{}, id).Error
}

// GetAllVideos retrieves all videos from the database.
func (m *VideoDBModel) GetAllVideos() ([]Video, error) {
	var videos []Video
	err := m.DB.Find(&videos).Error
	return videos, err
}

// GetPlaylistVideos retrieves all videos associated with a playlist.
func (m *VideoDBModel) GetPlaylistVideos(playlistID uint) ([]Video, error) {
	var videos []Video
	if err := m.DB.Where("playlist_id = ?", playlistID).Order("order").Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (v *Video) GetS3URL() string {
	return S3BaseURL + v.S3Key
}

// GetVideoURL retrieves the complete S3 URL for a video.
func (vm *VideoDBModel) GetVideoURL(videoID uint) (string, error) {
	video, err := vm.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	return video.GetS3URL(), nil
}

// SetVideoURL sets the S3 key for a video and updates the S3 URL.
func (vm *VideoDBModel) SetVideoURL(videoID uint, s3Key string) error {
	video, err := vm.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	video.SetS3Key(s3Key)
	video.UpdateS3URL()

	return vm.UpdateVideo(video)
}

func (v *Video) SetS3Key(s3Key string) {
	v.S3Key = s3Key
}

// UpdateS3URL updates the S3 URL based on the S3 key.
func (v *Video) UpdateS3URL() {
	v.URL = S3BaseURL + v.S3Key
}

/*-----------------------------------------------------------------------------------------------------------------------*/

// GetTotalDuration retrieves the total duration of a video.
func (vm *VideoDBModel) GetTotalDuration(videoID uint) (time.Duration, error) {
	var v Video
	if err := vm.DB.First(&v, videoID).Error; err != nil {
		return 0, err
	}
	return time.Duration(v.Duration) * time.Second, nil
}

// GetCurrentVideo retrieves the currently playing video for the specified playlistID.
func (vm *VideoDBModel) GetCurrentVideo(playlistID uint) (*Video, error) {
	var v Video
	if err := vm.DB.Where("playlist_id = ? AND is_playing = ?", playlistID, true).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// GetNextVideo retrieves the next video in the playlist for the specified playlistID.
func (vm *VideoDBModel) GetNextVideo(playlistID uint) (*Video, error) {
	var v Video
	if err := vm.DB.Where("playlist_id = ? AND order = ?", playlistID, 1).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (vm *VideoDBModel) GetTotalDuration2(videoID uint) (time.Duration, error) {
	var v Video
	if err := vm.DB.First(&v, videoID).Error; err != nil {
		return 0, err
	}
	return time.Duration(v.Duration) * time.Second, nil
}

func (vm *VideoDBModel) GetCurrentVideo2(playlistID uint) (*Video, error) {
	var v Video
	if err := vm.DB.Where("playlist_id = ? AND order = (SELECT MAX(order) FROM videos WHERE playlist_id = ?)", playlistID, playlistID).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// GetNextVideo retrieves the next video in the playlist for the specified playlistID
func (vm *VideoDBModel) GetNextVideo2(playlistID uint) (*Video, error) {
	// Implement logic to fetch the next video in the playlist from the database for the specified playlistID
	var v Video
	if err := vm.DB.Where("playlist_id = ? AND is_playing = ?", playlistID, false).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// GetVideoByID2 retrieves a video by its ID.
func (vm *VideoDBModel) GetVideoByID2(videoID uint) (*Video, error) {
	var v Video
	if err := vm.DB.Where("id = ?", videoID).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// GetNextVideo2 retrieves the next video in the playlist for the specified playlistID.
func (vm *VideoDBModel) GetNextVideo3(playlistID uint) (*Video, error) {
	var v Video
	if err := vm.DB.Where("playlist_id = ? AND order = (SELECT MAX(order) FROM videos WHERE playlist_id = ?) + 1", playlistID, playlistID).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// GetVideoByID3 retrieves a video by its ID.
func (m *VideoDBModel) GetVideoByID3(videoID uint) (*Video, error) {
	var video Video
	if err := m.DB.Where("id = ?", videoID).First(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// GetPlaybackHistory retrieves the playback history for a user.
func (vm *VideoDBModel) GetPlaybackHistory(userID uint) ([]*Video, error) {
	var videos []*Video
	if err := vm.DB.Joins("JOIN playback_histories ON videos.id = playback_histories.video_id").
		Where("playback_histories.user_id = ?", userID).
		Order("playback_histories.created_at desc").
		Limit(10).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetRelatedVideos retrieves videos related to the specified videoID.
func (vm *VideoDBModel) GetRelatedVideos(videoID uint) ([]*Video, error) {
	// Implement logic to fetch related videos based on your application's requirements.
	// This is a placeholder and should be replaced with the actual implementation.
	var videos []*Video
	return videos, nil
}

// GetMostLikedVideos retrieves the most liked videos up to the specified limit.
func (vm *VideoDBModel) GetMostLikedVideos(limit int) ([]*Video, error) {
	// Implement logic to fetch most liked videos based on your application's requirements.
	// This is a placeholder and should be replaced with the actual implementation.
	var videos []*Video
	return videos, nil
}

// GetVideosByUploader retrieves videos uploaded by the specified user.
func (vm *VideoDBModel) GetVideosByUploader(userID uint) ([]*Video, error) {
	// Implement logic to fetch videos uploaded by the specified user.
	// This is a placeholder and should be replaced with the actual implementation.
	var videos []*Video
	return videos, nil
}

// GetVideosByCategory retrieves videos in the specified category.
func (vm *VideoDBModel) GetVideosByCategory(categoryID uint) ([]*Video, error) {
	// Implement logic to fetch videos in the specified category.
	// This is a placeholder and should be replaced with the actual implementation.
	var videos []*Video
	return videos, nil
}

// IncrementViews increments the view count for the specified video.
func (vm *VideoDBModel) IncrementViews(videoID uint) error {
	// Implement logic to increment views for the specified video.
	// This is a placeholder and should be replaced with the actual implementation.
	return nil
}

// IncrementLikes increments the like count for the specified video.
func (vm *VideoDBModel) IncrementLikes(videoID uint) error {
	// Implement logic to increment likes for the specified video.
	// This is a placeholder and should be replaced with the actual implementation.
	return nil
}

// IncrementDislikes increments the dislike count for the specified video.
func (vm *VideoDBModel) IncrementDislikes(videoID uint) error {
	// Implement logic to increment dislikes for the specified video.
	// This is a placeholder and should be replaced with the actual implementation.
	return nil
}

// GetTopTags retrieves the top tags up to the specified limit.
func (vm *VideoDBModel) GetTopTags(limit int) ([]string, error) {
	// Implement logic to fetch top tags based on your application's requirements.
	// This is a placeholder and should be replaced with the actual implementation.
	var tags []string
	return tags, nil
}
