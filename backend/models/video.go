// backend/models/video.go

package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Video model
type Video struct {
	gorm.Model
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	ThumbnailURL    string    `json:"thumbnailUrl"`
	Views           uint      `json:"views" gorm:"default:0"`
	Likes           uint      `json:"likes" gorm:"default:0"`
	Dislikes        uint      `json:"dislikes" gorm:"default:0"`
	Comments        []Comment `gorm:"foreignKey:VideoID"`
	PlaylistID      uint      `json:"-"`
	Playlist        Playlist  `json:"playlist"`
	IsAdvertisement bool      `json:"isAdvertisement" gorm:"default:false"`
	//Duration        time.Duration  `json:"duration"` // Duration in seconds
	Order           int            `json:"order" gorm:"default:0"`
	Tags            []string       `json:"tags" gorm:"type:varchar(255)[]"`
	UploadDate      int            `json:"uploadDate" gorm:"autoCreateTime"`
	Uploader        User           `json:"uploader"`
	Category        Category       `json:"category"`
	PrivacySetting  PrivacySetting `json:"privacySetting" gorm:"embedded"`
	CommentsEnabled bool           `json:"commentsEnabled" gorm:"default:true"`
	RelatedVideos   []RelatedVideo `json:"relatedVideos" gorm:"foreignKey:VideoID"`
	Duration        uint           `json:"duration"`
	CreatorUserID   uint           `json:"creator_user_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	S3Key           string         `json:"s3key"`
	// Add more video-related fields as needed
}

// RelatedVideo model for representing related videos
type RelatedVideo struct {
	gorm.Model
	VideoID        uint
	RelatedVideoID uint
}

type VideoDBModel struct {
	DB *gorm.DB
}

type VideoModel interface {
	GetTotalDuration(videoID uint) (time.Duration, error)
	GetCurrentVideo(playlistID uint) (*Video, error)
	GetNextVideo(playlistID uint) (*Video, error)
	GetVideoByID(videoID uint) (*Video, error)
	GetPlaybackHistory(userID uint) ([]*Video, error)
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewVideoDBModel(db *gorm.DB) *VideoDBModel {
	return &VideoDBModel{
		DB: db,
	}
}

func (vm *VideoDBModel) GetTotalDuration(videoID uint) (time.Duration, error) {
	// Implement logic to fetch total duration from the database for the specified videoID
	// For example:
	// var v Video
	// vm.DB.First(&v, videoID)
	// return v.Duration, nil
	return 5 * time.Minute, nil // Placeholder value, replace with actual logic
}

func (vm *VideoDBModel) GetCurrentVideo(playlistID uint) (*Video, error) {
	// Implement logic to fetch the currently playing video from the database for the specified playlistID
	// For example:
	// var v Video
	// vm.DB.Where("playlist_id = ? AND is_playing = ?", playlistID, true).First(&v)
	// return &v, nil
	return &Video{ID: 1, Title: "Current Video", PlaylistID: playlistID}, nil // Placeholder value, replace with actual logic
}

func (vm *VideoDBModel) GetNextVideo(playlistID uint) (*Video, error) {
	// Implement logic to fetch the next video in the playlist from the database for the specified playlistID
	// For example:
	// var v Video
	// vm.DB.Where("playlist_id = ? AND is_playing = ?", playlistID, false).First(&v)
	// return &v, nil
	return &Video{ID: 2, Title: "Next Video", PlaylistID: playlistID}, nil // Placeholder value, replace with actual logic
}

// GetVideoByID retrieves a video by its ID
func (m *VideoDBModel) GetVideoByID(videoID uint) (*Video, error) {
	// Replace with the actual implementation to get the video by its ID
	// For example, query the database, fetch from cache, etc.
	// This is just a simulated example

	var video Video
	if err := m.DB.Where("id = ?", videoID).First(&video).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

func (vm *VideoDBModel) GetPlaybackHistory(userID uint) ([]*Video, error) {
	// Implement logic to fetch the playback history based on the provided user ID
	// Use the database connection (vm.DB) to query the data
	// Replace the following line with your actual database query logic
	return nil, errors.New("not implemented")
}
