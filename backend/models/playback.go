// backend/models/playback.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// PlayingVideo represents the currently playing video
type PlayingVideo struct {
	gorm.Model
	Video       *Video    `json:"video"`
	StartTime   time.Time `json:"start_time"`
	CurrentRate float64   `json:"current_rate" gorm:"default:0"`
	PlaylistID  *uint     `json:"playlist_id" gorm:"foreignKey:PlaylistID"`
	// Add more fields as needed
}

// TableName sets the table name for the Video model.
func (PlayingVideo) TableName() string {
	return "playingVideo"
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewPlayingVideoModel(video *Video, currentRate float64, playlistID uint) *PlayingVideo {
	return &PlayingVideo{
		Video:       video,
		StartTime:   time.Now(),
		CurrentRate: 1.0,
	}
}

type PlayingVideoDBModel struct {
	DB *gorm.DB
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewPlayingVideoDBModel(db *gorm.DB) *PlayingVideoDBModel {
	return &PlayingVideoDBModel{
		DB: db,
	}
}
